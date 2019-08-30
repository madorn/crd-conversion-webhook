# crd-conversion-webhook

## Requirements
* Kubernetes 1.15 (Latest version of `minikube` should work)   
* Kubernetes 1.13/1.14 will work if `CustomResourceWebhookConversion` Feature Gate is enabled
     
## Setup: Deploy the Webhook

### Create the dedicated namespace
```
kubectl create -f deploy/namespace.yaml
```
### Create the secret that contains the signed cert and private key
```
kubectl create -f deploy/sample-secret.yaml
```
### Create the webhook pod
```
kubectl create -f deploy/deployment.yaml
```
### Expose the webhook pod as a ClusterIP service
```
kubectl create -f deploy/service.yaml
```
