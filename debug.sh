#!/bin/bash

# Function to log command output
log_command() {
    echo "===================================================" >> debug.log
    echo "$(date): Executing: $1" >> debug.log
    echo "===================================================" >> debug.log
    eval "$1" >> debug.log 2>&1
    echo -e "\n" >> debug.log
}

# Clear existing debug.log
> debug.log

log_command "./hcloud-k3s kubectl get svc,endpoints -n kube-system"
log_command "./hcloud-k3s kubectl describe svc kubernetes-api-proxy -n kube-system"
log_command "./hcloud-k3s kubectl get ingressroute kubernetes-api-tls -n kube-system -o yaml"
log_command "./hcloud-k3s kubectl get ingress kubernetes-api -n kube-system -o yaml"
log_command "./hcloud-k3s kubectl get pods -n kube-system -l component=kube-apiserver"
log_command "./hcloud-k3s kubectl get helmcharts.helm.cattle.io traefik -n kube-system -o yaml"
log_command "./hcloud-k3s kubectl logs -n traefik -l app.kubernetes.io/name=traefik"
log_command "./hcloud-k3s kubectl get nodes"
log_command "./hcloud-k3s kubectl get pods --all-namespaces"
log_command "./hcloud-k3s kubectl config current-context"
log_command "./hcloud-k3s kubectl config view -o jsonpath='{.contexts[*].name}'"
log_command "./hcloud-k3s kubectl cluster-info"
log_command "./hcloud-k3s kubectl get pods -n kube-system -l component=kube-apiserver -o wide"
log_command "./hcloud-k3s kubectl get services -n kube-system"
log_command "./hcloud-k3s kubectl describe service kubernetes-api-proxy -n kube-system"
log_command "./hcloud-k3s kubectl get deployment -n traefik"
log_command "./hcloud-k3s kubectl describe deployment traefik -n traefik"
log_command "./hcloud-k3s kubectl get service traefik -n traefik"
log_command "./hcloud-k3s kubectl describe service traefik -n traefik"
log_command "./hcloud-k3s kubectl get configmaps -n traefik"
log_command "./hcloud-k3s kubectl logs -n traefik -l app.kubernetes.io/name=traefik --tail=100"
log_command "./hcloud-k3s kubectl get secrets --all-namespaces -o json | jq '.items[] | select(.type == \"kubernetes.io/tls\") | .metadata.name'"
log_command "./hcloud-k3s kubectl get ingressroutes --all-namespaces"
log_command "./hcloud-k3s kubectl get pods -n cert-manager"
log_command "./hcloud-k3s kubectl logs -n cert-manager -l app=cert-manager"
log_command "./hcloud-k3s kubectl get certificate -n cert-manager storyteller-plus-wildcard -o yaml"
log_command "./hcloud-k3s kubectl get events -n cert-manager --sort-by='.lastTimestamp'"
log_command "./hcloud-k3s kubectl get ingressroute -n kube-system kubernetes-api-tls -o yaml"
log_command "./hcloud-k3s kubectl get ingress -n kube-system kubernetes-api -o yaml"
log_command "./hcloud-k3s kubectl get service -n kube-system kubernetes-api-proxy -o yaml"
log_command "./hcloud-k3s kubectl get endpoints -n kube-system kubernetes-api-proxy -o yaml"
log_command "./hcloud-k3s kubectl logs -n traefik -l app.kubernetes.io/name=traefik"
log_command "./hcloud-k3s kubectl get events -n kube-system"
log_command "./hcloud-k3s kubectl get service -n traefik"
log_command "./hcloud-k3s kubectl run -it --rm --restart=Never curl-test --image=curlimages/curl -- curl https://10.43.0.1:6443 -k"
log_command "./hcloud-k3s kubectl get certificate -n cert-manager storyteller-plus-wildcard -o yaml"
log_command "./hcloud-k3s kubectl get nodes"
log_command "./hcloud-k3s kubectl get pods --all-namespaces"
log_command "./hcloud-k3s kubectl cluster-info"
log_command "./hcloud-k3s kubectl describe deployment traefik -n traefik"
log_command "./hcloud-k3s kubectl describe service traefik -n traefik"

log_command "curl -k https://k3s.storyteller.plus/version"

echo "Debugging information has been saved to debug.log"