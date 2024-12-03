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

### SecurityContext

```yaml
# https://kubernetes.io/docs/tasks/configure-pod-container/security-context/#set-the-security-context-for-a-pod
---
apiVersion: v1
kind: Pod
metadata:
  name: security-context-demo
spec:
  securityContext:
    runAsUser: 1000
    runAsGroup: 3000
    fsGroup: 2000
    supplementalGroups: [4000]
  volumes:
  - name: sec-ctx-vol
    emptyDir: {}
  containers:
  - name: sec-ctx-demo
    image: busybox:1.28
    command: [ "sh", "-c", "sleep 1h" ]
    volumeMounts:
    - name: sec-ctx-vol
      mountPath: /data/demo
    securityContext:
      allowPrivilegeEscalation: false

```

```sh
kubectl exec -it security-context-demo -- sh
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

## Ingress

```yaml
# https://kubernetes.io/zh-cn/docs/concepts/services-networking/ingress/#the-ingress-resource
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: minimal-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx-example
  rules:
  - http:
      paths:
      - path: /testpath
        pathType: Prefix
        backend:
          service:
            name: test
            port:
              number: 80

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

1. canary deployment
2. 
```sh
kubectl create ns goshawk
kubectl create deployment current-chipmunk-deployment --image=nginx --replicas=5 --dry-run=client -o yaml -n goshawk > current-chipmunk-deployment.yaml
kubectl expose -f current-chipmunk-deployment.yaml --port=80 -n goshawk --dry-run=client -o yaml > chipmunk-service.yaml

```


```yaml
# current-chipmunk-deployment.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: chipmunk
    track: stable
  name: current-chipmunk-deployment
  namespace: goshawk
spec:
  replicas: 6
  selector:
    matchLabels:
      app: chipmunk
      track: stable
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: chipmunk
        track: stable
    spec:
      containers:
      - image: nginx
        name: nginx
        resources: {}
status: {}
```

```yaml
# chipmunk-service.yaml
---
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  labels:
    app: chipmunk
  name: chipmunk-service
  namespace: goshawk
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: chipmunk
status:
  loadBalancer: {}

```

```sh

kubectl apply -f current-chipmunk-deployment.yaml
kubectl apply -f chipmunk-service.yaml
kubectl scale deployment current-chipmunk-deployment --replicas=10 -n goshawk

```

```yaml
# canary-chipmunk-deployment.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: chipmunk
    track: canary
  name: canary-chipmunk-deployment
  namespace: goshawk
spec:
  replicas: 4
  selector:
    matchLabels:
      app: chipmunk
      track: canary
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: chipmunk
        track: canary
    spec:
      containers:
      - image: nginx:latest
        name: nginx-canary
        resources: {}
status: {}
```

```sh
kubectl get svc -n goshawk -o wide
kubectl get ep -n goshawk
kubectl edit ep -n goshawk
kubectl get ep -n goshawk -o yaml | grep deployment
kubectl describe ep -n goshawk

```


1. ingress canary

```yaml
# base-app-svc.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: base-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: base-app
  template:
    metadata:
      labels:
        app: base-app
    spec:
      containers:
      - name: base-app
        image: anvesh35/echo-pod-name
        ports:
          - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: base-svc
  labels:
    app: base-app
spec:
  type: ClusterIP
  selector:
    app: base-app
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
```

```yaml
# canary-app-svc.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: canary-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: canary-app
  template:
    metadata:
      labels:
        app: canary-app
    spec:
      containers:
      - name: canary-app
        image: anvesh35/echo-pod-name
        ports:
          - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: canary-svc
  labels:
    app: canary-app
spec:
  type: ClusterIP
  selector:
    app: canary-app
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
```

```yaml
# base-ingress.yaml
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: base-ingress
spec:
  ingressClassName: nginx
  rules:
  - host: canary.echo.pod.name.com
    http:
      paths:
      - pathType: Prefix
        path: /
        backend:
          service:
            name: base-svc
            port:
              number: 80
```

```yaml
# canary-ingress.ymal
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: canary-ingress
  annotations:
    nginx.ingress.kubernetes.io/canary: "true"
    nginx.ingress.kubernetes.io/canary-weight: "30"
spec:
  ingressClassName: nginx
  rules:
  - host: canary.echo.pod.name.com
    http:
      paths:
      - pathType: Prefix
        path: /
        backend:
          service:
            name: canary-svc
            port:
              number: 80
```


```sh
for i in $(seq 1 10); do curl -s --resolve canary.echo.pod.name.com:80:<Ingress-Controller-IP> canary.echo.pod.name.com; done
```

```sh
kubectl create deployment api --image=nginx:1.16 --replicas=6 -n ckad00014 --dry-run=client -o yaml > api.yaml
```

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: api
  name: api
  namespace: ckad00014
spec:
  replicas: 6
  selector:
    matchLabels:
      app: api
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: api
    spec:
      containers:
      - image: nginx:1.16
        name: nginx
        ports:
        - containerPort: 80
        env:
        - name: NGINX_PORT
          value: "8000"
        resources: {}
status: {}

```

```sh
kubectl exec -n ckad00014 -it api-59f9d5cf5b-b22wm -- sh -c 'env'
```

```sh
kubectl create ns gorilla
kubectl taint nodes --all node-role.kubernetes.io/master-
kubectl taint nodes --all node-role.kubernetes.io/control-plane-
```

```yaml
# honeybee-deployment.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: honeybee-deployment
  namespace: gorilla
spec:
  replicas: 1
  selector:
    matchLabels:
      role: honeybee-deployment
  template:
    metadata:
      labels:
        role: honeybee-deployment
    spec:
      containers:
      - name: honeybee
        image: bitnami/kubectl:latest
        command: ["/bin/sh"]
        args: ["-c", "while true; do kubectl get sa; sleep 10;done"]
        tty: true
        securityContext:
          privileged: true
        volumeMounts:
        - name: kube-config
          mountPath: /root/.kube/
      volumes:
      - name: kube-config
        hostPath:
          path: /home/$USER/.kube/
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: node-role.kubernetes.io/control-plane
                operator: Exists
      tolerations:
      - effect: NoSchedule
        key: node-role.kubernetes.io/control-plane
        operator: Exists

```

```sh
kubectl apply -f honeybee-deployment.yaml
kubectl delete -f honeybee-deployment.yaml
kubectl create role gorilla --verb=get --verb=list --resource=serviceaccounts -n gorilla
kubectl create serviceaccount gorilla -n gorilla
kubectl delete rolebinding gorilla -n gorilla
kubectl create rolebinding gorilla --role=gorilla --serviceaccount=gorilla:gorilla --namespace=gorilla
kubectl set serviceaccount deployments honeybee-deployment gorilla -n gorilla
kubectl edit deployment honeybee-deployment -n gorilla
```



### 1. Create a PersistentVolume (PV):
You need to create a PersistentVolume named `earth-project-earthflower-pv` with the following requirements:
- Capacity of 2Gi
- Access mode `ReadWriteOnce`
- HostPath of `/Volumes/Data`
- No storageClassName defined.

Here’s how you can define the PersistentVolume (PV) in a YAML file:

```yaml
# earth-project-earthflower-pv.yaml
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: earth-project-earthflower-pv
spec:
  capacity:
    storage: 2Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /Volumes/Data
  persistentVolumeReclaimPolicy: Retain
  storageClassName: ""
```

### 2. Create a PersistentVolumeClaim (PVC):
Next, create a PersistentVolumeClaim named `earth-project-earthflower-pvc` in the `earth` namespace with the following requirements:
- Request for 2Gi of storage
- Access mode `ReadWriteOnce`
- No `storageClassName` defined.

Here’s the YAML definition for the PersistentVolumeClaim (PVC):

```yaml
# earth-project-earthflower-pvc.yaml
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: earth-project-earthflower-pvc
  namespace: earth
spec:
  resources:
    requests:
      storage: 2Gi
  accessModes:
    - ReadWriteOnce
  storageClassName: ""
```

### 3. Create a Deployment:
Finally, create a Deployment named `project-earthflower` in the `earth` namespace that mounts the PVC to `/tmp/project-data` and uses the image `httpd:2.4.41-alpine`.

Here’s the YAML definition for the Deployment:

```yaml
# project-earthflower-deployment.yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: project-earthflower
  name: project-earthflower
  namespace: earth
spec:
  replicas: 1
  selector:
    matchLabels:
      app: project-earthflower
  template:
    metadata:
      labels:
        app: project-earthflower
    spec:
      containers:
        - name: httpd
          image: httpd:2.4.41-alpine
          volumeMounts:
            - mountPath: /tmp/project-data
              name: project-data-volume
      volumes:
        - name: project-data-volume
          persistentVolumeClaim:
            claimName: earth-project-earthflower-pvc
```

### Complete Workflow:

1. **Apply the PV YAML:**
   ```bash
   kubectl apply -f earth-project-earthflower-pv.yaml
   ```

2. **Apply the PVC YAML in the `earth` namespace:**
   ```bash
   kubectl apply -f earth-project-earthflower-pvc.yaml
   ```

3. **Apply the Deployment YAML in the `earth` namespace:**
   ```bash
   kubectl apply -f project-earthflower-deployment.yaml
   ```

This will create the required PV, PVC, and Deployment in your Kubernetes cluster, ensuring that the PersistentVolume is properly bound to the PersistentVolumeClaim and that the deployment mounts the volume at the correct path.




```sh
kubectl create configmap some-config --from-literal=key3=value4
kubectl run nginx-configmap --image=nginx:stable --dry-run=client -o yaml > nginx-configmap.yaml
kubectl exec -it nginx-configmap -- sh -c 'cat /some/path/key3'
```

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: nginx-configmap
  name: nginx-configmap
spec:
  containers:
  - image: nginx:stable
    name: nginx-configmap
    resources: {}
    volumeMounts:
    - name: some-config
      mountPath: "/some/path"
      readOnly: true
  volumes:
  - name: some-config
    configMap:
      name: some-config
  dnsPolicy: ClusterFirst
  restartPolicy: Always
status: {}

```


```sh
kubectl create secret generic another-secret --from-literal=key1=value2 
kubectl run nginx-secret --image=nginx:stable --dry-run=client -o yaml > nginx-secret.yaml
kubectl exec -it nginx-secret -- sh -c 'env'
```

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: nginx-secret
  name: nginx-secret
spec:
  containers:
  - image: nginx:1.16
    name: nginx-secret
    resources: {}
    env:
    - name: COOL_VARIABLE
      valueFrom:
        secretKeyRef:
          name: another-secret
          key: key1
          optional: false
  dnsPolicy: ClusterFirst
  restartPolicy: Always
status: {}

```

```sh
kubectl create ns ckad00015
```
```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
 name: webapp
 namespace: ckad00015
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
       image: lfccncf/nginx:1.12.2-alpine
       ports:
       - containerPort: 80
```

```sh
kubectl edit deployment/webapp -n ckad00015
kubectl set image deployment/webapp nginx=lfccncf/nginx:1.13.7 -n ckad00015
kubectl get deployment webapp -n ckad00015 -o yaml
kubectl rollout undo deployment/webapp -n ckad00015
kubectl rollout history deployment/webapp -n ckad00015  --revision=3
```

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "3"
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{},"labels":{"app":"nginx"},"name":"webapp","namespace":"ckad00015"},"spec":{"replicas":3,"selector":{"matchLabels":{"app":"nginx"}},"template":{"metadata":{"labels":{"app":"nginx"}},"spec":{"containers":[{"image":"lfccncf/nginx:1.12.2-alpine","name":"nginx","ports":[{"containerPort":80}]}]}}}}
  creationTimestamp: "2024-11-24T10:06:40Z"
  generation: 4
  labels:
    app: nginx
  name: webapp
  namespace: ckad00015
  resourceVersion: "5532"
  uid: e0444ae4-ea53-4e81-aac4-dfcde20b59f0
spec:
  progressDeadlineSeconds: 600
  replicas: 3
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: nginx
  strategy:
    rollingUpdate:
      maxSurge: 10%
      maxUnavailable: 4
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: nginx
    spec:
      containers:
      - image: lfccncf/nginx:1.12.2-alpine
        imagePullPolicy: IfNotPresent
        name: nginx
        ports:
        - containerPort: 80
          protocol: TCP
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
status:
  availableReplicas: 3
  conditions:
  - lastTransitionTime: "2024-11-24T10:06:50Z"
    lastUpdateTime: "2024-11-24T10:06:50Z"
    message: Deployment has minimum availability.
    reason: MinimumReplicasAvailable
    status: "True"
    type: Available
  - lastTransitionTime: "2024-11-24T10:06:40Z"
    lastUpdateTime: "2024-11-24T10:07:32Z"
    message: ReplicaSet "webapp-855b6758c6" has successfully progressed.
    reason: NewReplicaSetAvailable
    status: "True"
    type: Progressing
  observedGeneration: 4
  readyReplicas: 3
  replicas: 3
  updatedReplicas: 3

```




```sh
kubectl create ns frontend
kubectl create deploy nginx --image=nginx -n frontend
kubectl create serviceaccount app -n frontend
kubectl set serviceaccount deployment nginx app -n frontend
```



```sh
kubectl create ns ckad00017
kubectl create deploy ckad00017-deployment --image=nginx -n ckad00017  --replicas=5 -o yaml --dry-run=client  | kubectl label --local -f - tier=dmz -o yaml > ckad00017-deployment.yml
kubectl expose deployment ckad00017-deployment --port=81 --target-port=80 --name=rover --type=NodePort  -n ckad00017 --dry-run=client -o yaml > ckad00017-deployment-svc.yml
kubectl get svc -n ckad00017
curl 10.108.207.7:81
```



```sh
kubectl create ns ckad00018
kubectl run ckad00018-newpod  --labels="front=access,db=access,run=ckad00018-newpod" -n ckad00018 --dry-run=client -o yaml --image=nginx:1.16 > ckad00018-newpod.yaml
kubectl run front --labels="run=front" -n ckad00018 --dry-run=client -o yaml --image=nginx:1.16 > ckad00018-front.yaml
kubectl run db --labels="run=db" -n ckad00018 --dry-run=client -o yaml --image=nginx:1.16 > ckad00018-db.yaml
kubectl run test -n ckad00018 --dry-run=client -o yaml --image=nginx:1.16 > ckad00018-test.yaml
kubectl get pods -n ckad00018 -o wide
kubectl exec -it ckad00018-newpod -n ckad00018 -- sh -c 'curl 192.168.1.4'



```

```yaml
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: access-front-policy
  namespace: ckad00018
spec:
  podSelector:
    matchLabels:
      front: access
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          run: front
    ports:
    - protocol: TCP
      port: 80
  egress:
  - to:
    - podSelector:
        matchLabels:
          run: front
    ports:
    - protocol: TCP
      port: 80
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: access-db-policy
  namespace: ckad00018
spec:
  podSelector:
    matchLabels:
      db: access
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          run: db
    ports:
    - protocol: TCP
      port: 80
  egress:
  - to:
    - podSelector:
        matchLabels:
          run: db
    ports:
    - protocol: TCP
      port: 80
```


```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-dm
  namespace: ingress-ckad
spec:
  replicas: 2
  selector:
    matchLabels:
      name: nginx-ing
  template:
    metadata:
      labels:
        name: nginx-ing
    spec:
      containers:
      - name: nginx
        image: vicuu/nginx:hello81
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 81
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-ing-svc
  namespace: ingress-ckad
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 81
  selector:
    name: nginx-ing
status:
  loadBalancer: {}
```
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: nginx-ingress-test
  namespace: ingress-ckad
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
  - http:
      paths:
      - path: /hello
        pathType: Prefix
        backend:
          service:
            name: nginx-ing-svc
            port:
              number: 80
```


```sh
kubectl create deployment webapp --image=nginx1
kubectl set image deployment/webapp nginx1=nginx:latest
```


```sh
# https://github.com/kubernetes-sigs/metrics-server
# https://medium.com/@cloudspinx/fix-error-metrics-api-not-available-in-kubernetes-aa10766e1c2f
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
kubectl patch deployment metrics-server -n kube-system --type='json' -p='[
{
"op": "add",
"path": "/spec/template/spec/hostNetwork",
"value": true
},
{
"op": "replace",
"path": "/spec/template/spec/containers/0/args",
"value": [
"--cert-dir=/tmp",
"--secure-port=4443",
"--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname",
"--kubelet-use-node-status-port",
"--metric-resolution=15s",
"--kubelet-insecure-tls"
]
},
{
"op": "replace",
"path": "/spec/template/spec/containers/0/ports/0/containerPort",
"value": 4443
}
]'
kubectl create ns cpu-stress
kubectl create deployment webapp --image=nginx --replicas=3 -n cpu-stress
kubectl top pod -n cpu-stress --sort-by=cpu | awk 'NR==2{print $1}'
```