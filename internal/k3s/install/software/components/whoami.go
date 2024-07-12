package components

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/template"
	"strings"
)

func WhoAmI(ctx clustercontext.ClusterContext) string {
	return template.CompileTemplate(`
apiVersion: v1
kind: Namespace
metadata:
  name: whoami
---
apiVersion: v1
kind: Service
metadata:
  name: whoami
  namespace: whoami
spec:
  ports:
    - name: web
      port: 80
      targetPort: web
  selector:
    app: whoami
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: whoami
  namespace: whoami
spec:
  selector:
    matchLabels:
      app: whoami
  template:
    metadata:
      labels:
        app: whoami
    spec:
      containers:
        - name: whoami
          image: traefik/whoami
          ports:
            - name: web
              containerPort: 80
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: whoami
  namespace: whoami
  annotations:
    kubernetes.io/ingress.class: "traefik"
spec:
  tls:
    - hosts:
        - whoami.{{ .Domain }}
      secretName: {{ .DomainKebap }}-wildcard-tls
  rules:
    - host: whoami.{{ .Domain }}
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: whoami
                port:
                  name: web
`,
		map[string]interface{}{
			"Domain":      ctx.Config.Domain,
			"DomainKebap": strings.ReplaceAll(ctx.Config.Domain, ".", "-"),
		})
}
