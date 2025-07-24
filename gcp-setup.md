# GCP Setup Guide

## Prerequisites

1. **GCP Project**: Create or select a GCP project
2. **APIs**: Enable required APIs:
   ```bash
   gcloud services enable cloudbuild.googleapis.com
   gcloud services enable run.googleapis.com
   gcloud services enable artifactregistry.googleapis.com
   ```

3. **Artifact Registry**: Create repository
   ```bash
   gcloud artifacts repositories create taskservice-repo \
     --repository-format=docker \
     --location=us-central1
   ```

## Environment Setup

### 1. MongoDB
For production, use MongoDB Atlas or Google Cloud MongoDB:

```bash
# MongoDB Atlas connection string example:
# mongodb+srv://username:password@cluster.mongodb.net/taskdb

# Or use Google Cloud Marketplace MongoDB
```

### 2. Message Queue
Replace RabbitMQ with Google Pub/Sub for cloud-native messaging:

```bash
# Create Pub/Sub topic
gcloud pubsub topics create task-events

# Create subscription
gcloud pubsub subscriptions create task-events-sub \
  --topic=task-events
```

## Deployment Commands

### 1. Initial Setup
```bash
# Set project
gcloud config set project YOUR_PROJECT_ID

# Configure Docker for Artifact Registry
gcloud auth configure-docker us-central1-docker.pkg.dev
```

### 2. Manual Build and Deploy
```bash
# Build and push manually
docker build -f Dockerfile.backend -t us-central1-docker.pkg.dev/YOUR_PROJECT_ID/taskservice-repo/taskservice-backend .
docker build -f Dockerfile.frontend -t us-central1-docker.pkg.dev/YOUR_PROJECT_ID/taskservice-repo/taskservice-frontend .

docker push us-central1-docker.pkg.dev/YOUR_PROJECT_ID/taskservice-repo/taskservice-backend
docker push us-central1-docker.pkg.dev/YOUR_PROJECT_ID/taskservice-repo/taskservice-frontend
```

### 3. Cloud Build Deploy
```bash
# Submit build
gcloud builds submit \
  --config=cloudbuild.yaml \
  --substitutions=_MONGO_URI="YOUR_MONGO_URI",_RABBITMQ_ENABLED="false"
```

### 4. Auto-deploy from GitHub
```bash
# Connect repository
gcloud builds triggers create github \
  --repo-name=go-tasks-microservice \
  --repo-owner=YOUR_GITHUB_USERNAME \
  --branch-pattern="^main$" \
  --build-config=cloudbuild.yaml
```

## Environment Variables

Update these in `cloudbuild.yaml` or set them during deployment:

```yaml
substitutions:
  _MONGO_URI: 'mongodb+srv://user:pass@cluster.mongodb.net/taskdb'
  _RABBITMQ_ENABLED: 'false'  # Use Pub/Sub instead
```

## Access URLs

After deployment, get service URLs:

```bash
# Backend URL
gcloud run services describe taskservice-backend \
  --region=us-central1 \
  --format='value(status.url)'

# Frontend URL  
gcloud run services describe taskservice-frontend \
  --region=us-central1 \
  --format='value(status.url)'
```

## Monitoring

### 1. Cloud Logging
```bash
# View logs
gcloud logs read "resource.type=cloud_run_revision AND resource.labels.service_name=taskservice-backend"
```

### 2. Cloud Monitoring
- Set up uptime checks
- Configure alerting policies
- Monitor resource usage

## Security

### 1. IAM Roles
```bash
# Create service account
gcloud iam service-accounts create taskservice-sa

# Grant necessary permissions
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:taskservice-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/cloudsql.client"
```

### 2. VPC Connector (if needed)
```bash
# Create VPC connector for private resources
gcloud compute networks vpc-access connectors create taskservice-connector \
  --region=us-central1 \
  --subnet=default \
  --subnet-project=YOUR_PROJECT_ID \
  --min-instances=2 \
  --max-instances=3
```

## Cost Optimization

1. **Cloud Run**: Pay-per-use, automatic scaling to zero
2. **MongoDB Atlas**: Use shared clusters for development
3. **Pub/Sub**: Pay per message
4. **Artifact Registry**: Delete old images regularly

## Scaling Considerations

1. **Backend**: Configure appropriate memory and CPU limits
2. **Frontend**: Use CDN for static assets
3. **Database**: Use connection pooling
4. **Message Queue**: Monitor Pub/Sub subscription backlog