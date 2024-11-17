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

## ResourceQuota

```yaml
# quota.yaml
---
apiVersion: v1
kind: List
items:
- apiVersion: v1
  kind: ResourceQuota
  metadata:
    name: pods-high
  spec:
    hard:
      cpu: "1000"
      memory: 200Gi
      pods: "10"
    scopeSelector:
      matchExpressions:
      - operator : In
        scopeName: PriorityClass
        values: ["high"]
- apiVersion: v1
  kind: ResourceQuota
  metadata:
    name: pods-medium
  spec:
    hard:
      cpu: "10"
      memory: 20Gi
      pods: "10"
    scopeSelector:
      matchExpressions:
      - operator : In
        scopeName: PriorityClass
        values: ["medium"]
- apiVersion: v1
  kind: ResourceQuota
  metadata:
    name: pods-low
  spec:
    hard:
      cpu: "5"
      memory: 10Gi
      pods: "10"
    scopeSelector:
      matchExpressions:
      - operator : In
        scopeName: PriorityClass
        values: ["low"]
```

```yaml
# quote-pod.yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: high-priority
spec:
  containers:
  - name: high-priority
    image: ubuntu
    command: ["/bin/sh"]
    args: ["-c", "while true; do echo hello; sleep 10;done"]
    resources:
      requests:
        memory: "10Gi"
        cpu: "500m"
      limits:
        memory: "10Gi"
        cpu: "500m"
  priorityClassName: high
```

## Deployment

```yaml
# https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#creating-a-deployment
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```

## Pod

```yaml
# https://kubernetes.io/docs/concepts/workloads/pods/#using-pods
---
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
```

## LimitRange

```yaml
# https://kubernetes.io/docs/concepts/policy/limit-range/#limitrange-and-admission-checks-for-pods
---
apiVersion: v1
kind: LimitRange
metadata:
  name: cpu-resource-constraint
spec:
  limits:
  - default: # this section defines default limits
      cpu: 500m
    defaultRequest: # this section defines default requests
      cpu: 500m
    max: # max and min define the limit range
      cpu: "1"
    min:
      cpu: 100m
    type: Container

```

## Task

Create a CronJob named ppi that runs a single-container Pod with the following configuration:

```yaml
- name: pi
  image: perl:5
  command: ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"]
```
Configure the CronJob with the following properties:

1. Runs every 5 minutes.
1. Retain 2 completed jobs.
1. Retain 4 failed jobs.
1. Never restart the Pod.
1. Stop the Pod after 8 seconds.

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

- schedule: Uses a Cron expression to schedule the task to run every 5 minutes.
- containers: Defines the container and the command to be executed in the Pod.
- restartPolicy: Set to Never, ensuring the Pod does not restart upon failure.
- backoffLimit: Prevents retries when the job fails.
- successfulJobsHistoryLimit and failedJobsHistoryLimit: Limit the number of retained successful and failed jobs.
- ttlSecondsAfterFinished: Ensures the Pod is terminated 8 seconds after completion.


```sh
kubectl apply -f cronjob-task-1.yaml
kubectl get cronjob
kubectl get jobs
kubectl logs $POD
```



1. Use busybox:stable image and execute date command
1. The job should run in 10 seconds or until Kubernetes terminates it
1. CronJob name and namespace should be "hello"
1. Verify the Job execution at least once

```yaml
# cronjob-task-2.yaml
---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "* * * * *" 
  jobTemplate:
    spec:
      activeDeadlineSeconds: 10 # Terminate Pod 10 seconds after start 
      template:
        spec:
          containers:
          - name: hello
            image: busybox:stable
            command:
            - /bin/sh
            - -c
            - date
          restartPolicy: OnFailure
```


```sh
kubectl apply -f cronjob-task-2.yaml
kubectl get cronjob
kubectl get jobs
kubectl logs $POD
```

```sh
kubectl run nginx-resources --dry-run=client -o yaml --image=nginx:1.16 > task-4.yaml
```

```yaml
# task-4.yaml
---
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: nginx-resources
  name: nginx-resources
spec:
  containers:
  - image: nginx:1.16
    name: nginx-resources
    resources:
      requests:
        cpu: "40m"
        memory: "50Mi"
  dnsPolicy: ClusterFirst
  restartPolicy: Always
status: {}
```

```sh
kubectl apply -f task-4.yaml
kubectl describe pod nginx-resources
```



```sh
kubectl create ns haddock
kubectl -n haddock get limitrange -o yaml
```

```yaml
# cpu-mem-resource-constraint.yaml
---
apiVersion: v1
kind: LimitRange
metadata:
  name: cpu-mem-resource-constraint
  namespace: haddock
spec:
  limits:
  - default: # this section defines default limits
      cpu: 500m
      memory: 16Mi
    defaultRequest: # this section defines default requests
      cpu: 500m
      memory: 8Mi
    max: # max and min define the limit range
      cpu: "1"
      memory: 40Mi
    min:
      cpu: 100m
    type: Container
```

```sh
kubectl apply -f cpu-mem-resource-constraint.yaml
```

```yaml
# haddock-nosql-deployment.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: nosql
  name: nosql
  namespace: haddock
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nosql
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: nosql
    spec:
      containers:
        - name: nginx
          image: nginx
          resources:
            limits:
              memory: "20Mi"
            requests:
              memory: "15Mi"
```

```sh
kubectl apply -f haddock-nosql-deployment.yaml
```
