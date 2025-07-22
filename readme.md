# Go API Deployment Guide

## 📁 Project Structure

```
go-api-app/
├── cmd/
│   └── api/
│       └── main.go          # Main application code
├── static/                  # Static files (optional)
├── Dockerfile              # Docker build configuration
├── flylike.toml            # Flylike deployment configuration
├── go.mod                  # Go module dependencies
├── go.sum                  # Go module checksums
└── README.md               # Project documentation
```

## 🚀 Quick Deployment

### 1. Create the Application
```bash
# Create new app for user_andi
./bin/flylike --user=user_andi apps create go-api-app --port=9090
```

### 2. Prepare ConfigMap and Secrets
```bash
# Create configuration ConfigMap
./bin/flylike --user=user_andi configmap create go-api-config \
  --from-literal=database_host=postgres \
  --from-literal=redis_host=redis \
  --from-literal=debug_mode=false

# Create secrets
./bin/flylike --user=user_andi secret create go-api-secrets \
  --from-literal=jwt_secret=your-super-secret-jwt-key \
  --from-literal=database_password=your-db-password \
  --from-literal=api_key=your-api-key
```

### 3. Deploy from git
```bash
# Build and deploy from source
./bin/flylike --user=user_andi apps deploy go-api-app \
  --dockerfile=Dockerfile \
  --config=flylike.toml
```

### 4. Deploy with Pre-built Image
```bash
# Deploy using pre-built image
./bin/flylike --user=user_andi apps deploy go-api-app \
  --image=your-registry/go-api-app:latest \
  --port=9090 \
  --config=flylike.toml
```

## 🔧 Configuration Options

### Environment Variables
The application reads these environment variables:
- `PORT`: Server port (default: 9090)
- `GIN_MODE`: Gin framework mode (release/debug)
- `LOG_LEVEL`: Logging level (info/debug/warn/error)
- `APP_ENV`: Application environment (production/development)
- `DATABASE_URL`: Database connection string
- `REDIS_URL`: Redis connection string
- `JWT_SECRET`: JWT signing secret

### Volume Mounts
- `/app/data`: Persistent storage for user data
- `/app/logs`: Application logs
- `/app/config`: Configuration files from ConfigMap
- `/app/secrets`: Secret files (JWT keys, etc.)
- `/app/cache`: Temporary cache storage

### Health Check
The application exposes a health check endpoint at `/health` that returns:
```json
{
  "status": "ok",
  "message": "Go API is healthy",
  "data": {
    "timestamp": "2025-07-12T10:30:00Z",
    "version": "1.0.0",
    "uptime": "2h30m45s"
  }
}
```

## 📊 API Endpoints

### User Management API
- `GET /api/v1/users` - Get all users
- `POST /api/v1/users` - Create new user
- `GET /api/v1/users/:id` - Get specific user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Example API Usage
```bash
# Get health status
curl http://localhost:9090/health

# Create a new user
curl -X POST http://localhost:9090/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'

# Get all users
curl http://localhost:9090/api/v1/users

# Get specific user
curl http://localhost:9090/api/v1/users/user-id-here
```

## 🔍 Monitoring and Debugging

### Check Application Status
```bash
# Check app status
./bin/flylike --user=user_andi apps status go-api-app

# View application logs
./bin/flylike --user=user_andi apps logs go-api-app

# Check resource usage
./bin/flylike --user=user_andi usage
```

### Scale Application
```bash
# Scale to 3 replicas
./bin/flylike --user=user_andi apps scale go-api-app 3

# Scale back to 1 replica
./bin/flylike --user=user_andi apps scale go-api-app 1
```

### Storage Management
```bash
# Check persistent storage usage
./bin/flylike --user=user_andi storage usage

# List persistent volumes
./bin/flylike --user=user_andi storage list

# Backup application data
./bin/flylike --user=user_andi storage backup go-api-app
```

## 🛠️ Development Workflow

### Local Development
```bash
# Run locally for development
go mod tidy
go run cmd/api/main.go

# Build locally
go build -o bin/go-api-app cmd/api/main.go
```

### Build and Test Docker Image
```bash
# Build Docker image
docker build -t go-api-app:latest .

# Run locally with Docker
docker run -p 9090:9090 \
  -e PORT=9090 \
  -e GIN_MODE=debug \
  -e LOG_LEVEL=debug \
  go-api-app:latest
```

### Update Deployment
```bash
# Update application with new image
./bin/flylike --user=user_andi apps update go-api-app \
  --image=your-registry/go-api-app:v1.1.0 \
  --config=flylike.toml

# Restart application
./bin/flylike --user=user_andi apps restart go-api-app
```

## 📋 Troubleshooting

### Common Issues
1. **Application won't start**: Check logs for port conflicts or missing environment variables
2. **Health check fails**: Verify the `/health` endpoint is accessible
3. **Storage issues**: Check PVC creation and mounting
4. **Configuration errors**: Verify ConfigMap and Secret creation

### Debug Commands
```bash
# Check pods in user namespace
kubectl get pods -n user-user_andi

# Describe pod for detailed info
kubectl describe pod -n user-user_andi deployment/go-api-app

# Check persistent volumes
kubectl get pvc -n user-user_andi

# Access pod shell for debugging
kubectl exec -it -n user-user_andi deployment/go-api-app -- /bin/sh
```

## 🎯 Production Checklist

- [ ] Set `GIN_MODE=release` in production
- [ ] Configure proper logging levels
- [ ] Set up persistent storage for data
- [ ] Configure health checks
- [ ] Set resource limits and requests
- [ ] Enable backup for critical data
- [ ] Configure secrets properly
- [ ] Set up monitoring and alerting
- [ ] Test scaling and failover
- [ ] Configure proper CORS settings

## 📈 Performance Optimization

### Resource Limits
```toml
# Add to flylike.toml
[resources]
cpu_limit = "500m"
memory_limit = "512Mi"
cpu_request = "100m"
memory_request = "128Mi"
```

### Storage Optimization
- Use appropriate storage classes (SSD for performance)
- Enable compression for logs
- Implement log rotation
- Regular cleanup of temporary files

### Scaling Strategy
- Start with 2 replicas for high availability
- Use horizontal pod autoscaling based on CPU/memory
- Configure readiness and liveness probes
- Use rolling updates for zero-downtime deployments