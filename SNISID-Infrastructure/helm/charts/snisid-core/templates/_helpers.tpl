{{/*
Expand the name of the chart.
*/}}
{{- define "snisid-core.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "snisid-core.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "snisid-core.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "snisid-core.labels" -}}
helm.sh/chart: {{ include "snisid-core.chart" . }}
{{ include "snisid-core.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
snisid.gov/tier: {{ .Values.podLabels.tier | default "tier-1" | quote }}
snisid.gov/region: {{ .Values.podLabels.region | default "core" | quote }}
snisid.gov/service: {{ .Values.podLabels.service | default "core-api" | quote }}
snisid.gov/owner: {{ .Values.podLabels.owner | default "equipe-nationale" | quote }}
snisid.gov/data-classification: {{ .Values.podLabels.dataClassification | default "restricted" | quote }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "snisid-core.selectorLabels" -}}
app.kubernetes.io/name: {{ include "snisid-core.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: api
app.kubernetes.io/part-of: snisid-core
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "snisid-core.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "snisid-core.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
