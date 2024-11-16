# k8s-demo


```sh
kubectl get pods
kubectl get replicasets.apps
kubectl get deployments.apps
kubectl describe deployments.apps $DEPLOYMENT | grep -i image
kubectl describe pod $POD

```

```yaml
# deployment-definition-1.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-1
spec:
  replicas: 2
  selector:
    matchLabels:
      app: busybox-pod
  template:
    metadata:
      labels:
        app: busybox-pod
    spec:
      containers:
        - name: busybox888
          image: busybox888
          command:
            - "/bin/sh"
            - "-c"
            - "while true; do date; sleep 5; done"
          
```


```sh
kubectl apply -f deployment-definition-1.yaml
kubectl describe deployments.apps deployment-1 | grep -i image

kubectl create deployment httpd-frontend --image=httpd:2.4-alpine
kubectl scale deployment --replicas=3 httpd-frontend
kubectl get deployments.apps httpd-frontend
kubectl describe deployments.apps httpd-frontend
kubectl get pods
kubectl logs $POD -f
```

```yaml
# namespace-dev.yaml
---
apiVersion: v1
kind: Namespace
metadata:
  name: dev
```

```yaml
# deployment-definition-1.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-1
  namespace: dev
spec:
  replicas: 2
  selector:
    matchLabels:
      app: busybox-pod
  template:
    metadata:
      labels:
        app: busybox-pod
    spec:
      containers:
        - name: busybox888
          image: busybox888
          command:
            - "/bin/sh"
            - "-c"
            - "while true; do date; sleep 5; done"
          
```

```sh
kubectl apply -f namespace-dev.yaml
kubectl apply -f deployment-definition-1.yaml
kubectl get namespaces
kubectl describe namespaces
kubectl config set-context $(kubectl config current-context) --namespace=dev
kubectl get pods
kubectl get replicasets.apps
kubectl get deployments.apps
```


```yaml
# compute-quota-dev.yaml
---
apiVersion: v1
kind: ResourceQuota
metadata:
  name: compute-quota
  namespace: dev
spec:
  hard:
    pods: "10"
    requests.cpu: "4"
    requests.memory: 5Gi
    limits.cpu: "10"
    limits.memory: 10Gi

```

```sh
kubectl apply -f compute-quota-dev.yaml
kubectl describe namespaces dev
```
