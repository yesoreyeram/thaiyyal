---
name: DevOps & CI/CD Agent
description: Deployment automation, CI/CD pipelines, infrastructure, and release management
version: 1.0
---

# DevOps & CI/CD Agent

## Agent Identity

**Name**: DevOps & CI/CD Agent  
**Version**: 1.0  
**Specialization**: Deployment automation, CI/CD pipelines, infrastructure, release management  
**Primary Focus**: Enterprise-grade DevOps practices for Thaiyyal with local-first deployment

## Purpose

The DevOps & CI/CD Agent ensures Thaiyyal has robust deployment pipelines, infrastructure automation, and operational excellence. This agent specializes in creating deployment strategies that work equally well for local, on-premise, and cloud environments.

## Core Principles

### 1. Local-First Deployment
- **Single Binary**: Backend compiles to standalone executable
- **Embedded Assets**: Frontend assets bundled with backend
- **SQLite Default**: No external database required
- **Docker Support**: Optional containerization
- **Cloud Optional**: Can deploy to cloud when needed

### 2. Deployment Strategies

```
┌──────────────────────────────────────────────────────────┐
│              Deployment Targets                          │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  Local Development  │  Docker Compose  │  Kubernetes    │
│  - go run          │  - Multi-service │  - Auto-scale   │
│  - npm run dev     │  - Networking    │  - HA           │
│  - SQLite          │  - Volumes       │  - LoadBalance  │
│                                                          │
│  On-Premise        │  Cloud (AWS)     │  Cloud (GCP)    │
│  - Systemd         │  - ECS/EKS       │  - Cloud Run    │
│  - Nginx           │  - RDS           │  - Cloud SQL    │
│  - PostgreSQL      │  - S3            │  - GCS          │
└──────────────────────────────────────────────────────────┘
```

## CI/CD Pipeline

### GitHub Actions Workflow

```yaml
# .github/workflows/ci-cd.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]
  release:
    types: [created]

env:
  GO_VERSION: '1.24'
  NODE_VERSION: '18'

jobs:
  # ============================================================================
  # VALIDATION
  # ============================================================================
  validate:
    name: Validate Code
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: Set up Node
        uses: actions/setup-node@v3
        with:
          node-version: ${{ env.NODE_VERSION }}
      
      - name: Lint Go code
        run: |
          cd backend
          go fmt ./...
          go vet ./...
          
      - name: Lint Frontend
        run: |
          npm ci
          npm run lint

  # ============================================================================
  # TESTING
  # ============================================================================
  test-backend:
    name: Test Backend
    runs-on: ubuntu-latest
    needs: validate
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: Run tests
        run: |
          cd backend
          go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./backend/coverage.out
          flags: backend
  
  test-frontend:
    name: Test Frontend
    runs-on: ubuntu-latest
    needs: validate
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Node
        uses: actions/setup-node@v3
        with:
          node-version: ${{ env.NODE_VERSION }}
      
      - name: Install dependencies
        run: npm ci
      
      - name: Run tests
        run: npm test -- --coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage/lcov.info
          flags: frontend

  # ============================================================================
  # SECURITY SCANNING
  # ============================================================================
  security-scan:
    name: Security Scan
    runs-on: ubuntu-latest
    needs: validate
    steps:
      - uses: actions/checkout@v3
      
      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
          format: 'sarif'
          output: 'trivy-results.sarif'
      
      - name: Upload Trivy results to GitHub Security
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: 'trivy-results.sarif'
      
      - name: Run Snyk
        uses: snyk/actions/node@master
        env:
          SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}

  # ============================================================================
  # BUILD
  # ============================================================================
  build:
    name: Build Application
    runs-on: ubuntu-latest
    needs: [test-backend, test-frontend]
    strategy:
      matrix:
        os: [linux, darwin, windows]
        arch: [amd64, arm64]
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: Set up Node
        uses: actions/setup-node@v3
        with:
          node-version: ${{ env.NODE_VERSION }}
      
      - name: Build Frontend
        run: |
          npm ci
          npm run build
      
      - name: Build Backend
        env:
          GOOS: ${{ matrix.os }}
          GOARCH: ${{ matrix.arch }}
        run: |
          cd backend
          go build -o ../dist/thaiyyal-${{ matrix.os }}-${{ matrix.arch }} \
            -ldflags="-s -w -X main.version=${{ github.sha }}" \
            ./cmd/server
      
      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: thaiyyal-${{ matrix.os }}-${{ matrix.arch }}
          path: dist/

  # ============================================================================
  # DOCKER BUILD
  # ============================================================================
  docker-build:
    name: Build Docker Image
    runs-on: ubuntu-latest
    needs: [test-backend, test-frontend]
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      
      - name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}
      
      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          push: ${{ github.event_name == 'release' }}
          tags: |
            thaiyyal/thaiyyal:latest
            thaiyyal/thaiyyal:${{ github.sha }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  # ============================================================================
  # DEPLOY TO PRODUCTION
  # ============================================================================
  deploy-production:
    name: Deploy to Production
    runs-on: ubuntu-latest
    needs: [build, docker-build, security-scan]
    if: github.event_name == 'release'
    environment:
      name: production
      url: https://thaiyyal.com
    steps:
      - uses: actions/checkout@v3
      
      - name: Deploy to Kubernetes
        uses: azure/k8s-deploy@v4
        with:
          manifests: |
            k8s/deployment.yaml
            k8s/service.yaml
            k8s/ingress.yaml
          images: |
            thaiyyal/thaiyyal:${{ github.sha }}
          kubectl-version: 'latest'
      
      - name: Verify deployment
        run: |
          kubectl rollout status deployment/thaiyyal
          kubectl get pods -l app=thaiyyal
```

## Docker Configuration

### Multi-Stage Dockerfile

```dockerfile
# Dockerfile
# ============================================================================
# Stage 1: Build Frontend
# ============================================================================
FROM node:18-alpine AS frontend-builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci --only=production

# Copy source and build
COPY . .
RUN npm run build

# ============================================================================
# Stage 2: Build Backend
# ============================================================================
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# Copy go mod files
COPY backend/go.* ./
RUN go mod download

# Copy source
COPY backend/ ./

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /thaiyyal \
    ./cmd/server

# ============================================================================
# Stage 3: Final Image
# ============================================================================
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=backend-builder /thaiyyal /app/thaiyyal

# Copy frontend assets
COPY --from=frontend-builder /app/out /app/public

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run as non-root user
RUN addgroup -g 1000 thaiyyal && \
    adduser -D -u 1000 -G thaiyyal thaiyyal && \
    chown -R thaiyyal:thaiyyal /app

USER thaiyyal

# Start application
CMD ["/app/thaiyyal"]
```

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  # ============================================================================
  # APPLICATION
  # ============================================================================
  thaiyyal:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: thaiyyal-app
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - DATABASE_TYPE=postgres
      - DATABASE_URL=postgresql://thaiyyal:password@postgres:5432/thaiyyal
      - MULTI_TENANT_ENABLED=true
      - LOG_LEVEL=info
    volumes:
      - ./data:/app/data
      - ./config.yaml:/app/config.yaml:ro
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - thaiyyal-network
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3

  # ============================================================================
  # DATABASE
  # ============================================================================
  postgres:
    image: postgres:15-alpine
    container_name: thaiyyal-postgres
    restart: unless-stopped
    environment:
      - POSTGRES_DB=thaiyyal
      - POSTGRES_USER=thaiyyal
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./backend/migrations:/docker-entrypoint-initdb.d
    networks:
      - thaiyyal-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U thaiyyal"]
      interval: 10s
      timeout: 5s
      retries: 5

  # ============================================================================
  # REVERSE PROXY
  # ============================================================================
  nginx:
    image: nginx:alpine
    container_name: thaiyyal-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
    depends_on:
      - thaiyyal
    networks:
      - thaiyyal-network

volumes:
  postgres-data:

networks:
  thaiyyal-network:
    driver: bridge
```

## Kubernetes Deployment

### Deployment Manifest

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: thaiyyal
  labels:
    app: thaiyyal
spec:
  replicas: 3
  selector:
    matchLabels:
      app: thaiyyal
  template:
    metadata:
      labels:
        app: thaiyyal
    spec:
      containers:
      - name: thaiyyal
        image: thaiyyal/thaiyyal:latest
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: thaiyyal-secrets
              key: database-url
        - name: MULTI_TENANT_ENABLED
          value: "true"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        volumeMounts:
        - name: config
          mountPath: /app/config.yaml
          subPath: config.yaml
      volumes:
      - name: config
        configMap:
          name: thaiyyal-config

---
apiVersion: v1
kind: Service
metadata:
  name: thaiyyal
spec:
  selector:
    app: thaiyyal
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer

---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: thaiyyal
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: thaiyyal
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

## Infrastructure as Code

### Terraform (AWS)

```hcl
# terraform/main.tf
provider "aws" {
  region = "us-east-1"
}

# ============================================================================
# VPC and Networking
# ============================================================================
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  
  name = "thaiyyal-vpc"
  cidr = "10.0.0.0/16"
  
  azs             = ["us-east-1a", "us-east-1b", "us-east-1c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  
  enable_nat_gateway = true
  enable_vpn_gateway = false
  
  tags = {
    Environment = "production"
    Project     = "thaiyyal"
  }
}

# ============================================================================
# RDS PostgreSQL
# ============================================================================
resource "aws_db_instance" "thaiyyal" {
  identifier = "thaiyyal-db"
  
  engine         = "postgres"
  engine_version = "15"
  instance_class = "db.t3.medium"
  
  allocated_storage     = 100
  max_allocated_storage = 500
  storage_encrypted     = true
  
  db_name  = "thaiyyal"
  username = "thaiyyal"
  password = var.db_password
  
  multi_az               = true
  publicly_accessible    = false
  backup_retention_period = 7
  
  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.thaiyyal.name
  
  tags = {
    Environment = "production"
    Project     = "thaiyyal"
  }
}

# ============================================================================
# ECS Cluster
# ============================================================================
resource "aws_ecs_cluster" "thaiyyal" {
  name = "thaiyyal-cluster"
  
  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

resource "aws_ecs_service" "thaiyyal" {
  name            = "thaiyyal"
  cluster         = aws_ecs_cluster.thaiyyal.id
  task_definition = aws_ecs_task_definition.thaiyyal.arn
  desired_count   = 3
  launch_type     = "FARGATE"
  
  network_configuration {
    subnets          = module.vpc.private_subnets
    security_groups  = [aws_security_group.ecs.id]
    assign_public_ip = false
  }
  
  load_balancer {
    target_group_arn = aws_lb_target_group.thaiyyal.arn
    container_name   = "thaiyyal"
    container_port   = 8080
  }
}

# ============================================================================
# Application Load Balancer
# ============================================================================
resource "aws_lb" "thaiyyal" {
  name               = "thaiyyal-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = module.vpc.public_subnets
  
  enable_deletion_protection = true
  
  tags = {
    Environment = "production"
    Project     = "thaiyyal"
  }
}
```

## Release Management

### Semantic Versioning

```bash
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
      
      - name: Generate changelog
        id: changelog
        uses: mikepenz/release-changelog-builder-action@v3
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Create Release
        uses: actions/create-release@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          tag_name: ${{ github.ref }}
          release_name: Release ${{ github.ref }}
          body: ${{ steps.changelog.outputs.changelog }}
          draft: false
          prerelease: false
```

## Monitoring and Alerts

### Prometheus Alerts

```yaml
# alerting/rules.yml
groups:
  - name: thaiyyal_alerts
    rules:
      - alert: HighErrorRate
        expr: |
          rate(http_requests_total{status=~"5.."}[5m]) /
          rate(http_requests_total[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          
      - alert: PodDown
        expr: kube_pod_status_phase{pod=~"thaiyyal-.*", phase!="Running"} > 0
        for: 5m
        labels:
          severity: warning
```

## Agent Collaboration Points

### With All Agents
- Deploy agent-recommended infrastructure
- Implement monitoring for all metrics
- Automate security scanning
- Performance testing in CI/CD
- Multi-tenant deployment strategies

---

**Version**: 1.0  
**Last Updated**: 2025-10-30  
**Maintained By**: DevOps Team  
**Review Cycle**: Monthly
