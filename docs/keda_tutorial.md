

## CONCEPTS

### 1. Kubernetes (K8s) Core
- Pod & Deployment: Container packaging, lifecycle, and replica management.
- K8s Architecture: Control Plane, Worker Nodes, and kubectl CLI.
- HPA (Horizontal Pod Autoscaler): Default resource-based (CPU/RAM) scaling mechanism.

### 2. KEDA (Kubernetes Event-driven Autoscaling)
- KEDA Architecture: Operator pattern, Metrics Adapter, and Controller.
- RabbitMQ Scaler: Direct connection to RabbitMQ Management API to monitor queueLength.
- ScaledObject CRD: Custom Resource Definition for target deployments, triggers, thresholds, min/max replicas.

## TO-DO IMPLEMENTATION
### Current phase (completed): Application-level scalability
- Outbox Pattern
- Idempotency handling
- Using Worker pool pattern
- Consumer ACK/NACK, Prefetch
- Retry mechanism

### Step 1: Local Environment Setup
- K8s Cluster: Install Kind.
- Package Manager: Install Helm (for installing KEDA and infrastructure tools).
- CLI Tooling: Install kubectl.

### Step 2: Containerization & Infra Deployment
- Dockerize Consumer: Write multi-stage Dockerfile for the Go consumer app.
- Deploy RabbitMQ: Install via Helm chart inside the local K8s cluster.
- Install KEDA: Deploy KEDA operator to the cluster using Helm.

### Step 3: K8s & KEDA Manifests Configuration
- Deployment YAML: Define Go consumer container specs, environment variables, and initial replica count (replicas: 1).
- ScaledObject YAML: Configure the RabbitMQ trigger, authentication credentials, and threshold (e.g., target queueLength: 100).

### Step 4: Load Testing & Verification
- Simulate Spike: Run publisher to flood RabbitMQ with 10,000+ messages.
- Monitor Scale-Up: Run kubectl get pods -w to watch automatic consumer pods generation.
- Monitor Scale-Down: Stop traffic, clear the queue, and verify pods terminate back to minimum count.