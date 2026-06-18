package webhook

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractRegistry(t *testing.T) {
	tests := []struct {
		name  string
		image string
		want  string
	}{
		{"docker hub official", "nginx:1.21", ""},
		{"docker hub user", "myuser/app:v1", ""},
		{"allowed registry", "registry.snisid.gouv.ht/app:1.0", "registry.snisid.gouv.ht"},
		{"ghcr.io", "ghcr.io/snisid/app:latest", "ghcr.io"},
		{"external registry", "docker.io/library/nginx:latest", "docker.io"},
		{"with tag and path", "docker.snisid.gouv.ht/team/app:v2", "docker.snisid.gouv.ht"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractRegistry(tt.image)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateLabels_MissingRequired(t *testing.T) {
	v := validateLabels(map[string]string{"foo": "bar"})
	assert.Len(t, v, 2)
	for _, pv := range v {
		assert.Equal(t, "ERROR", pv.Severity)
	}
}

func TestValidateLabels_AllPresent(t *testing.T) {
	v := validateLabels(map[string]string{
		"snisid.gouv.ht/component":   "api",
		"snisid.gouv.ht/environment": "prod",
	})
	assert.Empty(t, v)
}

func TestValidateAnnotations_NamespaceMismatch(t *testing.T) {
	v := validateAnnotations(map[string]string{"snisid.gouv.ht/namespace": "other"}, "expected-ns")
	require.Len(t, v, 1)
	assert.Equal(t, "ERROR", v[0].Severity)
}

func TestValidateAnnotations_InvalidSecurityTier(t *testing.T) {
	v := validateAnnotations(map[string]string{"snisid.gouv.ht/security-tier": "platinum"}, "ns")
	require.Len(t, v, 1)
	assert.Equal(t, "ERROR", v[0].Severity)
}

func TestValidateAnnotations_ValidSecurityTier(t *testing.T) {
	for _, tier := range []string{"gold", "silver", "bronze"} {
		v := validateAnnotations(map[string]string{"snisid.gouv.ht/security-tier": tier}, "ns")
		assert.Empty(t, v, "tier %s should be valid", tier)
	}
}

func TestValidateAnnotations_Nil(t *testing.T) {
	v := validateAnnotations(nil, "ns")
	assert.Empty(t, v)
}

func TestValidateContainerImage_Empty(t *testing.T) {
	c := Container{Name: "test"}
	v := validateContainerImage(c)
	require.Len(t, v, 1)
	assert.Equal(t, "ERROR", v[0].Severity)
}

func TestValidateContainerImage_LatestTag(t *testing.T) {
	c := Container{Name: "test", Image: "nginx:latest"}
	v := validateContainerImage(c)
	require.Len(t, v, 1)
	assert.Equal(t, "ERROR", v[0].Severity)
}

func TestValidateContainerImage_DisallowedRegistry(t *testing.T) {
	c := Container{Name: "test", Image: "evil.com/app:v1"}
	v := validateContainerImage(c)
	require.Len(t, v, 1)
	assert.Equal(t, "WARNING", v[0].Severity)
}

func TestValidateContainerImage_Allowed(t *testing.T) {
	c := Container{Name: "test", Image: "registry.snisid.gouv.ht/app:1.0"}
	v := validateContainerImage(c)
	assert.Empty(t, v)
}

func TestValidateContainerResources_NoLimits(t *testing.T) {
	c := Container{Name: "test", Image: "nginx:1.21"}
	v := validateContainerResources(c)
	require.Len(t, v, 1)
	assert.Equal(t, "WARNING", v[0].Severity)
}

func TestValidateContainerResources_MissingMemory(t *testing.T) {
	c := Container{
		Name:  "test",
		Image: "nginx:1.21",
		Resources: ResourceRequirements{
			Limits: map[string]string{"cpu": "500m"},
		},
	}
	v := validateContainerResources(c)
	require.Len(t, v, 1)
	assert.Contains(t, v[0].Field, "memory")
}

func TestValidateContainerResources_Full(t *testing.T) {
	c := Container{
		Name:  "test",
		Image: "nginx:1.21",
		Resources: ResourceRequirements{
			Limits: map[string]string{"cpu": "500m", "memory": "512Mi"},
		},
	}
	v := validateContainerResources(c)
	assert.Empty(t, v)
}

func TestValidateContainerSecurity_Nil(t *testing.T) {
	c := Container{Name: "test", Image: "nginx:1.21"}
	v := validateContainerSecurity(c)
	require.Len(t, v, 1)
	assert.Equal(t, "WARNING", v[0].Severity)
}

func TestValidateContainerSecurity_Privileged(t *testing.T) {
	priv := true
	c := Container{
		Name:  "test",
		Image: "nginx:1.21",
		SecurityContext: &ContainerSecurityContext{
			Privileged: &priv,
		},
	}
	v := validateContainerSecurity(c)
	found := false
	for _, p := range v {
		if p.Severity == "ERROR" && strings.Contains(p.Field, "privileged") {
			found = true
		}
	}
	assert.True(t, found)
}

func TestValidateContainerSecurity_PrivilegeEscalation(t *testing.T) {
	esc := true
	ro := true
	c := Container{
		Name:  "test",
		Image: "nginx:1.21",
		SecurityContext: &ContainerSecurityContext{
			AllowPrivilegeEscalation: &esc,
			ReadOnlyRootFilesystem:   &ro,
		},
	}
	v := validateContainerSecurity(c)
	found := false
	for _, p := range v {
		if p.Severity == "ERROR" && strings.Contains(p.Field, "allowPrivilegeEscalation") {
			found = true
		}
	}
	assert.True(t, found)
}

func TestValidateContainerSecurity_DropAllCapabilities(t *testing.T) {
	ro := true
	c := Container{
		Name:  "test",
		Image: "nginx:1.21",
		SecurityContext: &ContainerSecurityContext{
			ReadOnlyRootFilesystem: &ro,
			Capabilities: &Capabilities{
				Drop: []string{"ALL"},
			},
		},
	}
	v := validateContainerSecurity(c)
	for _, p := range v {
		assert.NotContains(t, p.Field, "capabilities")
	}
}

func TestValidateVolumeMounts_ForbiddenPath(t *testing.T) {
	c := Container{
		Name:  "test",
		Image: "nginx:1.21",
		VolumeMounts: []VolumeMount{
			{Name: "hostfs", MountPath: "/host/etc"},
		},
	}
	v := validateContainerVolumeMounts(c)
	require.Len(t, v, 1)
	assert.Equal(t, "ERROR", v[0].Severity)
}

func TestValidateVolumeMounts_Safe(t *testing.T) {
	c := Container{
		Name:  "test",
		Image: "nginx:1.21",
		VolumeMounts: []VolumeMount{
			{Name: "data", MountPath: "/var/lib/data"},
		},
	}
	v := validateContainerVolumeMounts(c)
	assert.Empty(t, v)
}

func TestValidateImagePullSecrets_Empty(t *testing.T) {
	v := validateImagePullSecrets(nil)
	require.Len(t, v, 1)
	assert.Equal(t, "WARNING", v[0].Severity)
}

func TestValidateImagePullSecrets_Configured(t *testing.T) {
	v := validateImagePullSecrets([]PullSecret{{Name: "regcred"}})
	assert.Empty(t, v)
}

func TestValidatePodSecurityContext_Nil(t *testing.T) {
	v := validatePodSecurityContext(nil)
	require.Len(t, v, 1)
}

func TestValidatePodSecurityContext_NonRoot(t *testing.T) {
	nr := false
	v := validatePodSecurityContext(&PodSecurityContext{RunAsNonRoot: &nr})
	require.Len(t, v, 1)
}

func TestValidateUserInfo_EmptyUsername(t *testing.T) {
	v := validateUserInfo(UserInfo{})
	require.Len(t, v, 2)
	assert.Equal(t, "ERROR", v[0].Severity)
}

func TestValidateUserInfo_Valid(t *testing.T) {
	v := validateUserInfo(UserInfo{
		Username: "admin",
		Groups:   []string{"snisid:admins"},
	})
	assert.Empty(t, v)
}

func TestValidateUserInfo_NoSNISIDGroup(t *testing.T) {
	v := validateUserInfo(UserInfo{
		Username: "admin",
		Groups:   []string{"developers"},
	})
	require.Len(t, v, 1)
	assert.Equal(t, "WARNING", v[0].Severity)
}

func TestBuildDenyMessage(t *testing.T) {
	v := []PolicyViolation{
		{Field: "a", Message: "err1", Severity: "ERROR"},
		{Field: "b", Message: "warn", Severity: "WARNING"},
	}
	msg := buildDenyMessage(v)
	assert.Contains(t, msg, "[a] err1")
	assert.NotContains(t, msg, "warn")
}

func TestValidatePod_SkipDelete(t *testing.T) {
	req := &AdmissionRequest{
		Kind:      GroupVersionKind{Kind: "Pod"},
		Operation: "DELETE",
	}
	v := validatePod(req)
	assert.Nil(t, v)
}

func TestValidatePod_InvalidJSON(t *testing.T) {
	req := &AdmissionRequest{
		Kind:      GroupVersionKind{Kind: "Pod"},
		Operation: "CREATE",
		Object:    json.RawMessage(`{bad json}`),
	}
	v := validatePod(req)
	require.Len(t, v, 1)
	assert.Equal(t, "ERROR", v[0].Severity)
}

func TestValidatePod_ValidPod(t *testing.T) {
	pod := Pod{
		APIVersion: "v1",
		Kind:       "Pod",
		Metadata: ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Labels: map[string]string{
				"snisid.gouv.ht/component":   "api",
				"snisid.gouv.ht/environment": "prod",
			},
		},
		Spec: PodSpec{
			Containers: []Container{
				{
					Name:  "app",
					Image: "registry.snisid.gouv.ht/app:1.0",
					Resources: ResourceRequirements{
						Limits: map[string]string{"cpu": "500m", "memory": "512Mi"},
					},
					SecurityContext: &ContainerSecurityContext{
						ReadOnlyRootFilesystem: boolPtr(true),
						Capabilities:           &Capabilities{Drop: []string{"ALL"}},
					},
				},
			},
			ImagePullSecrets: []PullSecret{{Name: "regcred"}},
			SecurityContext:  &PodSecurityContext{RunAsNonRoot: boolPtr(true)},
		},
	}
	raw, _ := json.Marshal(pod)
	req := &AdmissionRequest{
		UID:       "uid-1",
		Kind:      GroupVersionKind{Kind: "Pod"},
		Operation: "CREATE",
		Namespace: "default",
		Name:      "test-pod",
		Object:    raw,
		UserInfo:  UserInfo{Username: "admin", Groups: []string{"snisid:admins"}},
	}
	v := validatePod(req)
	assert.Empty(t, v)
}

func TestHandleAdmission_InvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{bad}`))
	HandleAdmission(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleAdmission_EmptyRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))
	HandleAdmission(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleAdmission_Rejected(t *testing.T) {
	review := AdmissionReview{
		Request: &AdmissionRequest{
			UID:       "uid-1",
			Kind:      GroupVersionKind{Kind: "Pod"},
			Operation: "CREATE",
			Namespace: "default",
			Object:    json.RawMessage(`{"metadata":{"labels":{}},"spec":{"containers":[{"name":"app","image":"nginx:latest"}]}}`),
			UserInfo:  UserInfo{Username: "admin", Groups: []string{"snisid:admins"}},
		},
	}
	body, _ := json.Marshal(review)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
	HandleAdmission(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp AdmissionReview
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.NotNil(t, resp.Response)
	assert.False(t, resp.Response.Allowed)
}

func boolPtr(b bool) *bool {
	return &b
}
