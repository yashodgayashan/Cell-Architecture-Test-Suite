# Cell Architecture Test Suite

A demonstration of cell-based architecture patterns using microservices, designed for testing scale-to-zero capabilities in Kubernetes environments.

## Overview

This project implements a cell-based architecture with two independent cells (Cell A and Cell B) that demonstrate:
- Service chaining within cells
- Cross-cell communication patterns
- Scale-to-zero testing capabilities
- Kubernetes-native service discovery
- React-based UI for testing

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         React Client                             │
│                    (nginx reverse proxy)                         │
└────────────────────┬───────────────────┬────────────────────────┘
                     │                   │
                     ▼                   ▼
        ┌────────────────────┐  ┌─────────────────────┐
        │       Cell A       │  │       Cell B        │
        ├────────────────────┤  ├─────────────────────┤
        │   frontend-svc     │  │  analytics-frontend │
        │         ↓          │  │          ↓          │
        │    user-svc        │  │   report-generator  │
        │         ↓          │  │          ↓          │
        │    auth-svc        │  │    (calls Cell A    │
        │                    │  │      user-svc)      │
        └────────────────────┘  └─────────────────────┘
```

### Cell A Services

1. **frontend-svc**: Entry point for Cell A
   - Exposes HTTP endpoint on port 8080
   - Calls user-svc to process requests
   - Path: `/cell-a-apps/frontend-svc`

2. **user-svc**: User management service
   - Validates requests with auth-svc
   - Returns user information with authentication status
   - Path: `/cell-a-apps/user-svc`

3. **auth-svc**: Authentication service
   - Provides authentication verification
   - Base service in the chain
   - Path: `/cell-a-apps/auth-svc`

### Cell B Services

1. **analytics-frontend**: Entry point for Cell B
   - Exposes HTTP endpoint on port 8080
   - Calls report-generator service
   - Path: `/cell-b-apps/analytics-frontend`

2. **report-generator**: Report generation service
   - Demonstrates cross-cell communication
   - Calls Cell A's user-svc for user data
   - Path: `/cell-b-apps/report-generator`

### React Client

- Modern React 18 application
- Two-button interface to test both cells
- Nginx reverse proxy configuration for API routing
- Path: `/react-client`

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Kubernetes cluster (for full deployment)
- Go 1.19+ (for local development)
- Node.js 16+ (for React development)

### Local Development

1. **Build all services**:
   ```bash
   # Build Cell A services
   cd cell-a-apps
   docker build -t cell-a-frontend ./frontend-svc
   docker build -t cell-a-user ./user-svc
   docker build -t cell-a-auth ./auth-svc

   # Build Cell B services
   cd ../cell-b-apps
   docker build -t cell-b-analytics ./analytics-frontend
   docker build -t cell-b-report ./report-generator

   # Build React client
   cd ../react-client
   docker build -t react-client .
   ```

2. **Run with Docker Compose** (create a docker-compose.yml):
   ```yaml
   version: '3.8'
   services:
     # Cell A services
     auth-svc:
       image: cell-a-auth
       ports:
         - "8081:8080"
     
     user-svc:
       image: cell-a-user
       ports:
         - "8082:8080"
       environment:
         - AUTH_SVC_URL=http://auth-svc:8080
     
     frontend-svc:
       image: cell-a-frontend
       ports:
         - "8083:8080"
       environment:
         - USER_SVC_URL=http://user-svc:8080
     
     # Cell B services
     report-generator:
       image: cell-b-report
       ports:
         - "8084:8080"
       environment:
         - USER_SVC_URL=http://user-svc:8080
     
     analytics-frontend:
       image: cell-b-analytics
       ports:
         - "8085:8080"
       environment:
         - REPORT_GENERATOR_URL=http://report-generator:8080
     
     # React client
     react-client:
       image: react-client
       ports:
         - "3000:80"
       environment:
         - REACT_APP_CELL_A_URL=http://frontend-svc:8080
         - REACT_APP_CELL_B_URL=http://analytics-frontend:8080
   ```

### Kubernetes Deployment

1. **Create namespaces**:
   ```bash
   kubectl create namespace cell-a
   kubectl create namespace cell-b
   kubectl create namespace frontend
   ```

2. **Deploy Cell A**:
   ```bash
   kubectl apply -f cell-a-apps/k8s/ -n cell-a
   ```

3. **Deploy Cell B**:
   ```bash
   kubectl apply -f cell-b-apps/k8s/ -n cell-b
   ```

4. **Deploy React Client**:
   ```bash
   kubectl apply -f react-client/k8s/ -n frontend
   ```

## Configuration

### Environment Variables

Each service can be configured using environment variables:

- **Cell A Services**:
  - `AUTH_SVC_URL`: URL for auth service (default: `http://auth-svc.cell-a.svc.cluster.local:8080`)
  - `USER_SVC_URL`: URL for user service (default: `http://user-svc.cell-a.svc.cluster.local:8080`)

- **Cell B Services**:
  - `REPORT_GENERATOR_URL`: URL for report generator (default: `http://report-generator.cell-b.svc.cluster.local:8080`)
  - `USER_SVC_URL`: URL for Cell A's user service (cross-cell communication)

- **React Client**:
  - `REACT_APP_CELL_A_URL`: URL for Cell A frontend service
  - `REACT_APP_CELL_B_URL`: URL for Cell B analytics frontend

## API Endpoints

### Cell A
- `GET /`: Returns a chain of responses from frontend → user → auth services

### Cell B
- `GET /`: Returns analytics data with embedded user information from Cell A

### Response Format
All services return JSON responses with the following structure:
```json
{
  "message": "Service-specific message",
  "service": "service-name",
  "hostname": "pod-hostname",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    // Service-specific data
  }
}
```

## Testing Scale-to-Zero

This architecture is designed to test Kubernetes scale-to-zero capabilities:

1. **Deploy with HPA and KEDA**:
   ```yaml
   apiVersion: keda.sh/v1alpha1
   kind: ScaledObject
   metadata:
     name: service-scaler
   spec:
     scaleTargetRef:
       name: your-deployment
     minReplicaCount: 0
     maxReplicaCount: 10
     triggers:
     - type: prometheus
       metadata:
         serverAddress: http://prometheus:9090
         metricName: http_requests_per_second
         threshold: '1'
   ```

2. **Test with the React client**:
   - Services scale to zero when idle
   - First request wakes up the service chain
   - Subsequent requests are served normally

## Development

### Running Go Services Locally

```bash
# Cell A - Auth Service
cd cell-a-apps/auth-svc
go run main.go

# Cell A - User Service
cd cell-a-apps/user-svc
AUTH_SVC_URL=http://localhost:8081 go run main.go

# Cell A - Frontend Service
cd cell-a-apps/frontend-svc
USER_SVC_URL=http://localhost:8082 go run main.go
```

### Running React Client Locally

```bash
cd react-client
npm install
npm start
```

## Architecture Patterns Demonstrated

1. **Cell Isolation**: Each cell operates independently with its own namespace and services
2. **Service Chaining**: Demonstrates synchronous service-to-service communication
3. **Cross-Cell Communication**: Shows how cells can depend on services in other cells
4. **Resilience**: Services handle downstream failures gracefully
5. **Observability**: Each service includes hostname in responses for debugging
6. **Configuration Management**: Environment-based configuration for flexibility

## Troubleshooting

### Common Issues

1. **Services not reachable**: Ensure Kubernetes DNS is working and services are in the correct namespaces
2. **Cross-cell communication fails**: Verify network policies allow communication between namespaces
3. **Scale-to-zero not working**: Check KEDA installation and metrics server configuration

### Debugging

- Check service logs: `kubectl logs -n <namespace> <pod-name>`
- Verify service discovery: `kubectl exec -n <namespace> <pod> -- nslookup <service-name>`
- Test endpoints directly: `kubectl port-forward -n <namespace> svc/<service-name> 8080:8080`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Designed for demonstrating cell-based architecture patterns
- Optimized for Kubernetes environments
- Built with cloud-native principles in mind