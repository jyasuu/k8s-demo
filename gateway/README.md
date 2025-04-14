

```yaml
---
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: cont-ap-gateway
  namespace: cont-ap
spec:
  gatewayClassName: istio
  listeners:
  - name: cont-ap-http
    port: 80
    protocol: HTTP
    allowedRoutes:
      namespaces:
        from: Same
  - name: cont-ap-https
    port: 443
    protocol: HTTPS
    hostname: "*.example.com"
    tls:
      certificateRefs:
      - kind: Secret
        name: gateway-tls
    allowedRoutes:
      namespaces:
        from: Same
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: route-whoami-test
  namespace: cont-ap
spec:
  parentRefs:
  - name: cont-ap-gateway
    sectionName: cont-ap-https
    namespace: cont-ap
  - name: cont-ap-gateway
    sectionName: cont-ap-http
    namespace: cont-ap
  hostnames:
  - "gateway.example.com"
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: whoami-svc
      port: 80
  - matches:
    - path:
        type: PathPrefix
        value: /api
    backendRefs:
    - name: spring-boot-service
      port: 80
---
apiVersion: v1
kind: Secret
metadata:
  name: gateway-tls
  namespace: cont-ap
data:
  tls.crt: "{ tls.crt file encode base64 }"
  tls.key: "{ tls.key file encode base64 }"
type: kubernetes.io/tls

---

---
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: cont-ap
  name: postgres
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:14.9
        env:
        - name: POSTGRES_DB
          value: postgres
        - name: POSTGRES_USER
          value: postgres
        - name: POSTGRES_PASSWORD
          value: postgres
        resources:
          requests:
            memory: "256Mi"  # Minimum memory
            cpu: "500m"      # Minimum CPU
          limits:
            memory: "512Mi"  # Maximum memory
            cpu: "1"         # Maximum CPU
        ports:
        - containerPort: 5432
---
apiVersion: v1
kind: Service
metadata:
  namespace: cont-ap
  name: postgres-service
spec:
  ports:
  - port: 5432
  selector:
    app: postgres
---

apiVersion: v1
kind: ConfigMap
metadata:
  namespace: cont-ap
  name: spring-boot-config
data:
  application.properties: |
    # Your external configuration goes here
    server.port=8080
    spring.datasource.driver-class-name=org.postgresql.Driver
    spring.datasource.url=jdbc:postgresql://postgres-service:5432/postgres
    spring.datasource.username=postgres
    spring.datasource.password=postgres
    spring.datasource.continue-on-error=true
    spring.jpa.generate-ddl=true
    spring.jpa.hibernate.ddl-auto=update
    #spring.jpa.hibernate.naming.physical-strategy=org.hibernate.boot.model.naming.PhysicalNamingStrategyStandardImpl
    spring.jpa.defer-datasource-initialization=true
    spring.jpa.database-platform=org.hibernate.dialect.PostgreSQLDialect

---
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: cont-ap
  name: spring-boot-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: spring-boot-app
  template:
    metadata:
      labels:
        app: spring-boot-app
    spec:
      imagePullSecrets:
      - name: spring-demo-secret
      containers:
      - name: spring-boot-container
        image: jyasu/spring-demo:latest
        ports:
        - containerPort: 8080
        volumeMounts:
        - name: config-volume
          mountPath: /config
        - mountPath: /data
          name: test-nfs
        env:
        - name: SPRING_CONFIG_LOCATION
          value: /config/application.properties
        resources:
          requests:
            memory: "256Mi"  # Minimum memory
            cpu: "500m"      # Minimum CPU
          limits:
            memory: "512Mi"  # Maximum memory
            cpu: "1"         # Maximum CPU
      volumes:
      - name: config-volume
        configMap:
          name: spring-boot-config
      - name: test-nfs
        nfs:
          server: nfs.example.com
          path: /mnt/nfs_share

---
apiVersion: v1
kind: Service
metadata:
  namespace: cont-ap
  name: spring-boot-service
spec:
  type: NodePort
  ports:
  - port: 80
    targetPort: 8080
  selector:
    app: spring-boot-app
```
