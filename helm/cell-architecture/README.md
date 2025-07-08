# Cell Architecture Helm Chart

This Helm chart deploys a cell-based architecture test suite with network isolation using Cilium network policies.

## Architecture Overview

The chart deploys:
- **Cell A**: Auth service, User service, and Frontend service
- **Cell B**: Report generator and Analytics frontend
- **React Client**: Frontend application that can communicate with both cells
- **Network Policies**: Cilium policies for cell isolation and controlled communication

## Prerequisites

- Kubernetes cluster with Cilium CNI
- Helm 3.x
- Cilium network policies support enabled

## Installation

### Basic Installation

```bash
# Install with default values
helm install cell-architecture ./helm/cell-architecture

# Install with custom values
helm install cell-architecture ./helm/cell-architecture -f custom-values.yaml
```

### Configuration

Update the `values.yaml` file or create a custom values file:

```yaml
# Example custom values
global:
  registry: "your-registry.com"
  imageTag: "v1.0.0"

cellA:
  enabled: true
  userService:
    replicas: 0  # Scale to 0 for testing

networkPolicies:
  enabled: true
  crossCell:
    enabled: true  # Enable cross-cell communication for testing
```

## Testing Network Isolation

### 1. Deploy with Default Settings (Isolation Enabled)

```bash
helm install cell-architecture ./helm/cell-architecture
```

This deploys cells with network isolation. Cross-cell communication is blocked by default.

### 2. Test Cross-Cell Communication

```bash
# Scale user service to 0 to test scale-to-zero scenarios
kubectl scale deployment user-svc --replicas=0 -n cell-a

# Enable cross-cell communication
helm upgrade cell-architecture ./helm/cell-architecture --set networkPolicies.crossCell.enabled=true
```

### 3. Verify Network Policies

```bash
# Check Cilium network policies
kubectl get cnp -A

# Test connectivity between cells
kubectl exec -it <pod-name> -n cell-b -- curl http://user-svc.cell-a:8080/health
```

## Scaling Services

```bash
# Scale user service to 0 (scale-to-zero test)
kubectl scale deployment user-svc --replicas=0 -n cell-a

# Scale back up
kubectl scale deployment user-svc --replicas=1 -n cell-a
```

## Network Policy Management

### Enable Cross-Cell Communication

```bash
helm upgrade cell-architecture ./helm/cell-architecture --set networkPolicies.crossCell.enabled=true
```

### Disable Cross-Cell Communication

```bash
helm upgrade cell-architecture ./helm/cell-architecture --set networkPolicies.crossCell.enabled=false
```

## Customization

### Image Registry

```yaml
global:
  registry: "your-registry.com"
  imageTag: "latest"
```

### Service Configuration

```yaml
cellA:
  authService:
    replicas: 2
    resources:
      limits:
        memory: "256Mi"
        cpu: "200m"
```

### Network Policy Configuration

```yaml
networkPolicies:
  enabled: true
  cilium:
    enabled: true
  crossCell:
    enabled: false
  dns:
    namespace: "kube-system"
    app: "kube-dns"
```

## Uninstallation

```bash
helm uninstall cell-architecture
```

## Directory Structure

```
helm/cell-architecture/
├── Chart.yaml
├── values.yaml
├── README.md
├── templates/
│   ├── _helpers.tpl
│   ├── namespaces.yaml
│   ├── cell-a/
│   │   ├── auth-service.yaml
│   │   ├── user-service.yaml
│   │   └── frontend-service.yaml
│   ├── cell-b/
│   │   ├── report-generator.yaml
│   │   └── analytics-frontend.yaml
│   ├── react-client/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── ingress.yaml
│   └── network-policies/
│       ├── default-deny.yaml
│       ├── dns-allow.yaml
│       ├── intra-cell.yaml
│       ├── client-access.yaml
│       └── cross-cell.yaml
```

## Troubleshooting

### Check Pod Status

```bash
kubectl get pods -A
```

### Check Network Policies

```bash
kubectl get cnp -A
kubectl describe cnp <policy-name> -n <namespace>
```

### Debug Connectivity

```bash
# From Cell B to Cell A (should fail without cross-cell policy)
kubectl exec -it <cell-b-pod> -n cell-b -- curl http://user-svc.cell-a:8080/health

# Check Cilium connectivity
kubectl exec -it <cilium-pod> -n kube-system -- cilium policy get
```

## Contributing

1. Update the chart version in `Chart.yaml`
2. Test changes with `helm lint` and `helm template`
3. Update this README with any new features or configuration options