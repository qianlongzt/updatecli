source:
  kind: helmChart
  spec:
    name: prometheus
    url: https://prometheus-community.github.io/helm-charts
conditions:
  isNamePrometheus:
    kind: yaml
    spec:
      file: "helmfile.d/prometheus.yaml"
      key: "releases[0].name"
      value: "prometheus"
    scm:
      github:
        user: {{ .github.user }}
        email: {{ .github.email }}
        owner: {{ .github.owner }}
        repository: {{ .github.repository }}
        token: {{ requiredEnv "GITHUB_TOKEN" }}
        username: {{ .github.username }}
        branch: {{ .github.branch }}
targets:
  chartVersion:
    name: "prometheus Helm Chart"
    kind: yaml
    spec:
      file: "helmfile.d/prometheus.yaml"
      key: "releases[0].version"
    scm:
      github:
        user: {{ .github.user }}
        email: {{ .github.email }}
        owner: {{ .github.owner }}
        repository: {{ .github.repository }}
        token: {{ requiredEnv "GITHUB_TOKEN" }}
        username: {{ .github.username }}
        branch: {{ .github.branch }}
