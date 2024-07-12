package components

func WhoAmI() string {
	return `
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
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: le-example-http
  namespace: whoami
spec:
  acme:
    email: milan@kappen.name
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: storyteller-plus-le-example-http
    solvers:
      - http01:
          ingress:
            class: traefik
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: whoami
  namespace: whoami
  annotations:
    cert-manager.io/issuer: "le-example-http"
    kubernetes.io/ingress.class: "traefik"
    cert-manager.io/acme-challenge-type: http01
spec:
  tls:
    - hosts:
        - whoami.storyteller.plus
      secretName: tls-whoami-ingress-http
  rules:
    - host: whoami.storyteller.plus
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: whoami
                port:
                  name: web
`
}
