# Deployment Guide

This is a guide to deploy FrontService on your cluster. The contents are as follows.

* [Prerequisites](#prerequisites)
* [Deploy FrontService](#deploy-frontservice)
* [Set Ingress](#set-ingress-optional)

## Prerequisites
- [User Service](https://github.com/theraffle/userservice/blob/main/docs/deployment.md)
- [Project Service](https://github.com/theraffle/projectservice/blob/main/docs/deployment.md)
- Kubernetes Cluster
- (optional) Ingress Controller

## Deploy FrontService 
1. Apply `frontservice.yaml`
   ```bash
   kubectl apply -f ./kubernetes-manifests/release.yaml
   ```
2. Wait until `frontservice` pod is ready
   ```bash
   kubectl get pod -w
   ```

# Set Ingress (optional)
You can expose your API server easily by using ingress. The sample is like below.
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: gce
  name: example-ingress
  namespace: default
spec:
  defaultBackend:
    service:
      name: webfront
      port:
        number: 443
  rules:
  - host: api.sample.host
    http:
      paths:
      - backend:
          service:
            name: frontservice
            port:
              number: 80
        path: /
        pathType: Prefix
```
```bash
kubectl get ing   
NAME                          CLASS    HOSTS              ADDRESS         PORTS   AGE
example-ingress               <none>   api.sample.host    xx.xxx.xxx.xx   80      4d
```