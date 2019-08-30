# crd-conversion-webhook

## Requirements
* Kubernetes 1.15 (Latest version of `minikube` should work)   
* Kubernetes 1.13/1.14 will work if `CustomResourceWebhookConversion` [Feature Gate](https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/) is enabled
* [jq](https://stedolan.github.io/jq/download/)
     
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
## Exercise: Simulate the CRD upgrade process

### Create the initial `v1beta` CRD

```
kind: CustomResourceDefinition
metadata:
  name: crontabs.stable.example.com
spec:
  group: stable.example.com
  versions:
  - name: v1beta1
    served: true
    storage: true
  scope: Namespaced
  names:
    plural: crontabs
    singular: crontab
    kind: CronTab
    shortNames:
    - ct
  validation:
    openAPIV3Schema:
      type: object
      properties:
        spec:
          type: object
          properties:
            hostPort:
              type: string
```

```
kubectl create -f deploy/crd1.yaml
```

### Create the `v1beta` CR

```
apiVersion: stable.example.com/v1beta1
kind: CronTab
metadata:
  name: cr1
  namespace: crd-conversion-webhook
spec:
  hostPort: "localhost:8080"
```

```
kubectl create -f deploy/cr1.yaml
```

### Confirm you can get the `v1beta1` CR

```
kubectl get --raw /apis/stable.example.com/v1beta1/namespaces/crd-conversion-webhook/crontabs/cr1 | jq
```

```
apiVersion: stable.example.com/v1beta1
kind: CronTab
metadata:
  creationTimestamp: 2019-08-30T02:47:31Z
  generation: 1
  name: cr1
  namespace: crd-conversion-webhook
  resourceVersion: "237787"
  selfLink: /apis/stable.example.com/v1beta1/namespaces/crd-conversion-webhook/crontabs/cr1
  uid: f4b7e2b3-cd57-492f-b177-a7dd9df3b1e5
spec:
  hostPort: localhost:8080
```

### Bump the CRD to `v1`
### It contains an updated schema along with webhook conversion info.

```
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: crontabs.stable.example.com
spec:
  preserveUnknownFields: false
  group: stable.example.com
  versions:
  - name: v1beta1
    served: true
    storage: false
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              hostPort:
                type: string
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              host:
                type: string
              port:
                type: string
  conversion:
    strategy: Webhook
    webhookClientConfig:
      service:
        namespace: crd-conversion-webhook
        name: crd-conversion-webhook
        path: /crdconvert
      caBundle: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1ekNDQWMrZ0F3SUJBZ0lCQVRBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwdGFXNXAKYTNWaVpVTkJNQjRYRFRFNU1ERXdPVEU1TkRNd01Gb1hEVEk1TURFd056RTVORE13TUZvd0ZURVRNQkVHQTFVRQpBeE1LYldsdWFXdDFZbVZEUVRDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTnhECmE3NEptc2NVSGRxSWVOQ29nQjZHaUllMnhZai9nWUl2WU82eGVCSmFISlF0NXJZdUVMcWpEbU9qK1R4QisxQUUKK2UyRnNkNXplME94WWx2V3FOTWVrWm5USVE2ZCtqRnZwc2JBNHMrbW5wQkVuR04vUXA1WXdEdDBRZnEyd0x1QwpCaW0rWHMwTVNROHEyWUZDRUpsUUhUbU5YeGRlOFcySVJaK1R0RTV2U0V4VW5oL2M4dDgzd3VqS2ZwTXV0MXJIClhFOEJzVjBjcm5sOVBmR0lsbFJ4VHp5ZFMzWUlIZXBJaDQrb251Q1dSTHhEOU1Lci9GRUV1ZWJObloxYng0cmQKd0JHU1p1MjQ2Q1BlYmhmckxhOTdZdEtZVnRCTm1FRlJpdHRRV2RqZXBLWDRnS2trSnNTeU1xSVI4VWlCRkFJZApneldDc2pvKzluR25XVk5TMEtNQ0F3RUFBYU5DTUVBd0RnWURWUjBQQVFIL0JBUURBZ0trTUIwR0ExVWRKUVFXCk1CUUdDQ3NHQVFVRkJ3TUNCZ2dyQmdFRkJRY0RBVEFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxR1NJYjMKRFFFQkN3VUFBNElCQVFBY3lScmtEWVhybDlLVE55dzc4TVY0K1dDUEhUVEEwa3kvZE9vdFJnV0pRTXl2Yk5CNworZ3hvRGxSR0pNZmVDU09uL0NoaVhoWFNldk5EWVk4UmZ6Zm0yWjFMNVJuRmRuS1ZaaG5ZWnFjaWZ5WlFtZTI1Ck5zLzEwUHlrTkplWFpmdEhvTnNqbFc0dndvWkJsdmxrVEtrOExmTVBUSTVERUxTa3ZJK0ZxNTNhM3REdXlvanEKYlQ5cEYwTGxacDR2Rk9SOGVYbmtaVTN5ZDVQbDFBMWhqcDJVTG9aV05JY2pEbWVhQ3dsMTdZdlpLUlpHWCt6bApUNktOS0pEYm04a0pZSzhrTCtpT3VwaGZaRUpFaUJwcEwwOXNIeHU1Y0FNRVp4c3JSakhNMDdrclpjdnBiUU5sCmU2b1hndmUvZEJrM1FVdFVRbStETTRkSTdEa1BjRGI1R3Y3VQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
  scope: Namespaced
  names:
    plural: crontabs
    singular: crontab
    kind: CronTab
    shortNames:
    - ct
 ```
 
 ```
 kubectl replace -f deploy/crd2.yaml
 ```
 
### `cr1` CR is stored as `v1beta1`, but it can now serve as `v1` and reflect a `v1` schema.

```
kubectl get --raw /apis/stable.example.com/v1/namespaces/crd-conversion-webhook/crontabs/cr1 | jq
```

```
apiVersion: stable.example.com/v1
kind: CronTab
metadata:
  creationTimestamp: 2019-08-30T02:47:31Z
  generation: 1
  name: cr1
  namespace: crd-conversion-webhook
  resourceVersion: "237787"
  selfLink: /apis/stable.example.com/v1/namespaces/crd-conversion-webhook/crontabs/cr1
  uid: f4b7e2b3-cd57-492f-b177-a7dd9df3b1e5
spec:
  host: localhost
  port: "8080"
```

### Create a `v1` CR that will be stored in etcd as `v1`

```
apiVersion: stable.example.com/v1
kind: CronTab
metadata:
  name: cr2
  namespace: crd-conversion-webhook
spec:
  host: "localhost"
  port: "8080"
```

```
kubectl create -f deploy/cr2.yaml
```

### Confirm you can get the `v1` CR

```
kubectl get --raw apis/stable.example.com/v1/namespaces/crd-conversion-webhook/crontabs/cr2 | jq
```

```
apiVersion: stable.example.com/v1
kind: CronTab
metadata:
  creationTimestamp: 2019-08-30T02:51:07Z
  generation: 1
  name: cr2
  namespace: crd-conversion-webhook
  resourceVersion: "238052"
  selfLink: /apis/stable.example.com/v1/namespaces/crd-conversion-webhook/crontabs/cr2
  uid: 553ae7ab-0e4d-4255-83cb-760c8b1b730b
spec:
  host: localhost
  port: "8080"
```

**Note**: At this point etcd is storing cr1 (`v1beta1`) and cr2 (`v1`). When you are ready, you can bump all existing stored `v1beta1` versions to `v1` by following the instructions [here](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definition-versioning/#upgrade-existing-objects-to-a-new-stored-version)


### If you want to downgrade, update the CRD to no longer serve and store `v1`

```
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: crontabs.stable.example.com
spec:
  preserveUnknownFields: false
  group: stable.example.com
  versions:
  - name: v1beta1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              hostPort:
                type: string
  - name: v1
    served: false
    storage: false
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              host:
                type: string
              port:
                type: string
  conversion:
    strategy: Webhook
    webhookClientConfig:
      service:
        namespace: crd-conversion-webhook
        name: crd-conversion-webhook
        path: /crdconvert
      caBundle: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1ekNDQWMrZ0F3SUJBZ0lCQVRBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwdGFXNXAKYTNWaVpVTkJNQjRYRFRFNU1ERXdPVEU1TkRNd01Gb1hEVEk1TURFd056RTVORE13TUZvd0ZURVRNQkVHQTFVRQpBeE1LYldsdWFXdDFZbVZEUVRDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTnhECmE3NEptc2NVSGRxSWVOQ29nQjZHaUllMnhZai9nWUl2WU82eGVCSmFISlF0NXJZdUVMcWpEbU9qK1R4QisxQUUKK2UyRnNkNXplME94WWx2V3FOTWVrWm5USVE2ZCtqRnZwc2JBNHMrbW5wQkVuR04vUXA1WXdEdDBRZnEyd0x1QwpCaW0rWHMwTVNROHEyWUZDRUpsUUhUbU5YeGRlOFcySVJaK1R0RTV2U0V4VW5oL2M4dDgzd3VqS2ZwTXV0MXJIClhFOEJzVjBjcm5sOVBmR0lsbFJ4VHp5ZFMzWUlIZXBJaDQrb251Q1dSTHhEOU1Lci9GRUV1ZWJObloxYng0cmQKd0JHU1p1MjQ2Q1BlYmhmckxhOTdZdEtZVnRCTm1FRlJpdHRRV2RqZXBLWDRnS2trSnNTeU1xSVI4VWlCRkFJZApneldDc2pvKzluR25XVk5TMEtNQ0F3RUFBYU5DTUVBd0RnWURWUjBQQVFIL0JBUURBZ0trTUIwR0ExVWRKUVFXCk1CUUdDQ3NHQVFVRkJ3TUNCZ2dyQmdFRkJRY0RBVEFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxR1NJYjMKRFFFQkN3VUFBNElCQVFBY3lScmtEWVhybDlLVE55dzc4TVY0K1dDUEhUVEEwa3kvZE9vdFJnV0pRTXl2Yk5CNworZ3hvRGxSR0pNZmVDU09uL0NoaVhoWFNldk5EWVk4UmZ6Zm0yWjFMNVJuRmRuS1ZaaG5ZWnFjaWZ5WlFtZTI1Ck5zLzEwUHlrTkplWFpmdEhvTnNqbFc0dndvWkJsdmxrVEtrOExmTVBUSTVERUxTa3ZJK0ZxNTNhM3REdXlvanEKYlQ5cEYwTGxacDR2Rk9SOGVYbmtaVTN5ZDVQbDFBMWhqcDJVTG9aV05JY2pEbWVhQ3dsMTdZdlpLUlpHWCt6bApUNktOS0pEYm04a0pZSzhrTCtpT3VwaGZaRUpFaUJwcEwwOXNIeHU1Y0FNRVp4c3JSakhNMDdrclpjdnBiUU5sCmU2b1hndmUvZEJrM1FVdFVRbStETTRkSTdEa1BjRGI1R3Y3VQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
  scope: Namespaced
  names:
    plural: crontabs
    singular: crontab
    kind: CronTab
    shortNames:
    - ct
```

```
kubectl apply -f deploy/crd3-downgrade.yaml
```

### `cr2` should now appear to serve from `v1beta1` and reflect the `v1beta1` schema

```
kubectl get --raw /apis/stable.example.com/v1beta1/namespaces/crd-conversion-webhook/crontabs/cr2 | jq
```

```
apiVersion: stable.example.com/v1beta1
kind: CronTab
metadata:
  creationTimestamp: 2019-08-30T02:51:07Z
  generation: 1
  name: cr2
  namespace: crd-conversion-webhook
  resourceVersion: "238052"
  selfLink: /apis/stable.example.com/v1beta1/namespaces/crd-conversion-webhook/crontabs/cr2
  uid: 553ae7ab-0e4d-4255-83cb-760c8b1b730b
spec:
  hostPort: localhost:8080
```
