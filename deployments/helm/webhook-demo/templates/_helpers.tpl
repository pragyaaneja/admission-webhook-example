{{/*
Expand the name of the chart.
*/}}
{{- define "webhook-demo.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "webhook-demo.fullname" -}}
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
{{- define "webhook-demo.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "webhook-demo.labels" -}}
helm.sh/chart: {{ include "webhook-demo.chart" . }}
{{ include "webhook-demo.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "webhook-demo.selectorLabels" -}}
app.kubernetes.io/name: {{ include "webhook-demo.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "webhook-demo.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "webhook-demo.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Generate certificates for custom-metrics api server 
*/}}
{{- define "webhook-demo-certs" -}}
{{- $altNames := list ( printf "%s.%s" (include "webhook-demo.name" .) .Release.Namespace ) ( printf "%s.%s.svc" (include "webhook-demo.name" .) .Release.Namespace ) -}}
{{- $ca := genCA "webhook-demo-ca" 365 -}}
{{- $cert := genSignedCert ( include "webhook-demo.name" . ) nil $altNames 365 $ca -}}
tls.crt: {{ $cert.Cert | b64enc }}
tls.key: {{ $cert.Key | b64enc }}
{{- end -}}

{{/*
Generate CaBundle for custom-metrics api server 
*/}}
{{- define "webhook-demo-ca" -}}
{{- $altNames := list ( printf "%s.%s" (include "webhook-demo.name" .) .Release.Namespace ) ( printf "%s.%s.svc" (include "webhook-demo.name" .) .Release.Namespace ) -}}
{{- $ca := genCA "webhook-demo-ca" 365 -}}
caBundle: {{ b64enc $ca.Cert }}
{{- end -}}