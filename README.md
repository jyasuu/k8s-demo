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

```sh
kubectl api-resources
kubectl explain namespaces
kubectl explain deployments
kubectl explain replicasets
kubectl explain configmaps
kubectl explain secrets
kubectl explain services
kubectl explain cronjobs
```

## Cronjob

```yaml
# https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#example
---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "* * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox:1.28
            imagePullPolicy: IfNotPresent
            command:
            - /bin/sh
            - -c
            - date; echo Hello from the Kubernetes cluster
          restartPolicy: OnFailure
```

### Task

- schedule: Uses a Cron expression to schedule the task to run every 5 minutes.
- containers: Defines the container and the command to be executed in the Pod.
- restartPolicy: Set to Never, ensuring the Pod does not restart upon failure.
- backoffLimit: Prevents retries when the job fails.
- successfulJobsHistoryLimit and failedJobsHistoryLimit: Limit the number of retained successful and failed jobs.
- ttlSecondsAfterFinished: Ensures the Pod is terminated 8 seconds after completion.

```yaml
# cronjob-task-1..yaml
---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: ppi
spec:
  schedule: "*/5 * * * *" # Run every 5 minutes
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: pi
            image: perl:5
            command: ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"]
          restartPolicy: Never # Never restart the Pod
      backoffLimit: 0 # Do not retry on failure
      activeDeadlineSeconds: 8 # Terminate Pod 8 seconds after start 
      ttlSecondsAfterFinished: 8 # Terminate Pod 8 seconds after completion 
  successfulJobsHistoryLimit: 2 # Retain 2 successful jobs
  failedJobsHistoryLimit: 4 # Retain 4 failed jobs

```

```sh
kubectl apply -f cronjob-task-1.yaml
kubectl get cronjob
kubectl get jobs
kubectl logs $POD

```
