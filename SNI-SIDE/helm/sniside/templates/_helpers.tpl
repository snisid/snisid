{{- define "sniside.labels" -}}
app.kubernetes.io/name: {{ .name }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: {{ .component | default "microservice" }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version }}
{{- end -}}

{{- define "sniside.selectorLabels" -}}
app.kubernetes.io/name: {{ .name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "sniside.image" -}}
{{ .image | default (printf "%s/%s" .Values.global.imageRegistry .name) }}:{{ .Values.global.imageTag }}
{{- end -}}

{{- define "sniside.env.common" -}}
- name: KAFKA_BOOTSTRAP_SERVERS
  value: "{{ .Values.global.kafkaBootstrapServers }}"
- name: NEO4J_URI
  value: "{{ .Values.global.neo4jUri }}"
- name: REDIS_HOST
  value: "{{ .Values.global.redisHost }}"
- name: MILVUS_HOST
  value: "{{ .Values.global.milvusHost }}"
- name: MINIO_ENDPOINT
  value: "{{ .Values.global.minioEndpoint }}"
- name: CLICKHOUSE_HOST
  value: "{{ .Values.global.clickhouseHost }}"
- name: OTEL_EXPORTER_OTLP_ENDPOINT
  value: "{{ .Values.global.otelEndpoint }}"
- name: ENVIRONMENT
  value: "{{ .Values.global.environment }}"
{{- end -}}

{{- define "sniside.pod" -}}
- name: {{ .name }}
  image: "{{ .image }}:{{ $.Values.global.imageTag }}"
  imagePullPolicy: "{{ $.Values.global.imagePullPolicy }}"
  ports:
    - containerPort: 8080
      name: http
    - containerPort: 9100
      name: metrics
  env:
    {{- include "sniside.env.common" $ | nindent 4 }}
    {{- with .extraEnv }}{{ toYaml . | nindent 4 }}{{ end }}
  resources:
    {{- toYaml .resources | nindent 4 }}
  livenessProbe:
    httpGet:
      path: /health
      port: 8080
    initialDelaySeconds: 30
    periodSeconds: 15
  readinessProbe:
    httpGet:
      path: /health
      port: 8080
    initialDelaySeconds: 15
    periodSeconds: 10
  securityContext:
    runAsNonRoot: true
    runAsUser: 1001
    capabilities:
      drop: ["ALL"]
  {{- if .gpu }}
  resources:
    limits:
      nvidia.com/gpu: {{ .gpuCount | default 1 }}
  {{- end }}
{{- end -}}
