# GitHub Actions Workflows

This directory contains GitHub Actions workflows for building and pushing Docker images to Docker Hub.

## Prerequisites

Before using these workflows, ensure you have:

1. A Docker Hub account (username: `yashodperera`)
2. A Personal Access Token (PAT) from Docker Hub
3. Added the PAT as a GitHub secret named `DOCKER_PAT`

### Setting up Docker Hub PAT

1. Log in to [Docker Hub](https://hub.docker.com)
2. Go to Account Settings → Security → Access Tokens
3. Create a new access token with "Read, Write, Delete" permissions
4. Copy the token and add it to your GitHub repository:
   - Go to Settings → Secrets and variables → Actions
   - Click "New repository secret"
   - Name: `DOCKER_PAT`
   - Value: Your Docker Hub access token

## Workflows Overview

### 1. Reusable Go Build (`reusable-go-build.yml`)

A reusable workflow for building Go services. Used by other workflows to avoid duplication.

**Features:**
- Docker Buildx for multi-platform builds
- GitHub Actions cache for faster builds
- Tags images with both `latest` and commit SHA

### 2. Cell A Build (`cell-a-build.yml`)

Builds all Cell A services when changes are detected in the `cell-a-apps/` directory.

**Services built:**
- `yashodperera/cell-a-auth-svc`
- `yashodperera/cell-a-user-svc`
- `yashodperera/cell-a-frontend-svc`

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main`
- Manual workflow dispatch

### 3. Cell B Build (`cell-b-build.yml`)

Builds all Cell B services when changes are detected in the `cell-b-apps/` directory.

**Services built:**
- `yashodperera/cell-b-analytics-frontend`
- `yashodperera/cell-b-report-generator`

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main`
- Manual workflow dispatch

### 4. React Client Build (`react-client-build.yml`)

Builds the React client application when changes are detected in the `react-client/` directory.

**Image built:**
- `yashodperera/cell-architecture-react-client`

**Features:**
- Multi-stage Docker build
- Build arguments for API URLs
- Optimized for production

### 5. Build All Services (`build-all.yml`)

A comprehensive workflow that builds all services in the monorepo.

**Features:**
- Triggered on push to `main` branch
- Manual workflow dispatch with options to select which services to build
- Parallel builds for faster execution
- Build summary with all image tags

**Manual Dispatch Options:**
- `build-cell-a`: Build Cell A services (default: true)
- `build-cell-b`: Build Cell B services (default: true)
- `build-react-client`: Build React client (default: true)

## Usage Examples

### Manual Build of Specific Services

1. Go to Actions tab in GitHub
2. Select "Build All Services" workflow
3. Click "Run workflow"
4. Uncheck services you don't want to build
5. Click "Run workflow" button

### Automatic Builds

The workflows automatically trigger when you:
- Push changes to specific directories
- Create pull requests
- Push to `main` or `develop` branches

## Docker Images

All images are pushed to Docker Hub under the `yashodperera` namespace:

| Service | Image Name | Tags |
|---------|------------|------|
| Auth Service | `yashodperera/cell-a-auth-svc` | `latest`, `<commit-sha>` |
| User Service | `yashodperera/cell-a-user-svc` | `latest`, `<commit-sha>` |
| Frontend Service | `yashodperera/cell-a-frontend-svc` | `latest`, `<commit-sha>` |
| Analytics Frontend | `yashodperera/cell-b-analytics-frontend` | `latest`, `<commit-sha>` |
| Report Generator | `yashodperera/cell-b-report-generator` | `latest`, `<commit-sha>` |
| React Client | `yashodperera/cell-architecture-react-client` | `latest`, `<commit-sha>` |

## Monitoring Builds

- Check the Actions tab for build status
- View build logs for debugging
- Check Docker Hub for pushed images
- Review the build summary in completed workflows

## Troubleshooting

### Build Failures

1. **Authentication errors**: Verify `DOCKER_PAT` secret is set correctly
2. **Docker build errors**: Check Dockerfile syntax and build context
3. **Path filtering issues**: Ensure file changes match the configured paths

### Common Issues

- **Rate limiting**: Docker Hub has pull rate limits for anonymous users
- **Cache misses**: First builds take longer due to no cache
- **Concurrent builds**: May hit Docker Hub push limits if too many parallel builds

## Best Practices

1. Use specific tags for production deployments (commit SHA)
2. Keep `latest` tag for development/testing
3. Monitor Docker Hub storage usage
4. Regularly update base images for security
5. Use workflow dispatch for manual control when needed