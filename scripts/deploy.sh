#!/bin/bash

# User App Deployment Script
# This script helps deploy the user app to different environments

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
ENVIRONMENT="local"
REGISTRY=""
NAMESPACE="user-app"
IMAGE_TAG="latest"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -e, --environment ENV    Deployment environment (local|k8s|aws) [default: local]"
    echo "  -r, --registry REGISTRY Docker registry URL"
    echo "  -t, --tag TAG          Docker image tag [default: latest]"
    echo "  -n, --namespace NS     Kubernetes namespace [default: user-app]"
    echo "  -h, --help             Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 --environment local"
    echo "  $0 --environment k8s --registry your-registry.com --tag v1.0.0"
    echo "  $0 --environment aws --registry 123456789.dkr.ecr.us-west-2.amazonaws.com"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--environment)
            ENVIRONMENT="$2"
            shift 2
            ;;
        -r|--registry)
            REGISTRY="$2"
            shift 2
            ;;
        -t|--tag)
            IMAGE_TAG="$2"
            shift 2
            ;;
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option $1"
            show_usage
            exit 1
            ;;
    esac
done

# Validate environment
if [[ ! "$ENVIRONMENT" =~ ^(local|k8s|aws)$ ]]; then
    print_error "Invalid environment. Must be one of: local, k8s, aws"
    exit 1
fi

print_status "Starting deployment for environment: $ENVIRONMENT"

# Function to deploy locally with Docker Compose
deploy_local() {
    print_status "Deploying locally with Docker Compose..."
    
    # Stop existing containers
    docker-compose down 2>/dev/null || true
    
    # Build and start services
    docker-compose up --build -d
    
    print_success "Local deployment completed!"
    print_status "Application is available at: http://localhost:8080"
    print_status "API documentation: http://localhost:8080/swagger/"
    
    # Show logs
    print_status "Showing application logs (Ctrl+C to exit):"
    docker-compose logs -f app
}

# Function to deploy to Kubernetes
deploy_k8s() {
    print_status "Deploying to Kubernetes..."
    
    # Check if kubectl is available
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl is not installed or not in PATH"
        exit 1
    fi
    
    # Check if we can connect to cluster
    if ! kubectl cluster-info &> /dev/null; then
        print_error "Cannot connect to Kubernetes cluster"
        exit 1
    fi
    
    # Build and push image if registry is provided
    if [[ -n "$REGISTRY" ]]; then
        print_status "Building and pushing Docker image..."
        IMAGE_NAME="${REGISTRY}/user-app:${IMAGE_TAG}"
        docker build -t "$IMAGE_NAME" .
        docker push "$IMAGE_NAME"
        
        # Update image in deployment
        sed -i.bak "s|image: user-app:latest|image: $IMAGE_NAME|g" terraform/k8s/app-deployment.yaml
    fi
    
    # Apply Kubernetes manifests
    print_status "Applying Kubernetes manifests..."
    kubectl apply -f terraform/k8s/namespace.yaml
    kubectl apply -f terraform/k8s/configmap.yaml
    kubectl apply -f terraform/k8s/secret.yaml
    kubectl apply -f terraform/k8s/postgres-deployment.yaml
    
    # Wait for database to be ready
    print_status "Waiting for database to be ready..."
    kubectl wait --for=condition=available --timeout=300s deployment/postgres -n "$NAMESPACE"
    
    # Run database migrations
    print_status "Running database migrations..."
    kubectl apply -f terraform/k8s/migration-configmap.yaml
    kubectl apply -f terraform/k8s/migration-job.yaml
    kubectl wait --for=condition=complete --timeout=300s job/db-migration -n "$NAMESPACE"
    
    # Deploy application
    print_status "Deploying application..."
    kubectl apply -f terraform/k8s/app-deployment.yaml
    kubectl apply -f terraform/k8s/ingress.yaml
    
    # Wait for deployments to be ready
    print_status "Waiting for deployments to be ready..."
    kubectl wait --for=condition=available --timeout=300s deployment/user-app -n "$NAMESPACE"
    kubectl wait --for=condition=available --timeout=300s deployment/postgres -n "$NAMESPACE"
    
    print_success "Kubernetes deployment completed!"
    
    # Show service information
    print_status "Getting service information..."
    kubectl get services -n "$NAMESPACE"
    kubectl get ingress -n "$NAMESPACE"
}

# Function to deploy to AWS EKS
deploy_aws() {
    print_status "Deploying to AWS EKS..."
    
    # Check if required tools are available
    for tool in aws terraform kubectl; do
        if ! command -v "$tool" &> /dev/null; then
            print_error "$tool is not installed or not in PATH"
            exit 1
        fi
    done
    
    # Initialize and apply Terraform
    print_status "Initializing Terraform..."
    cd terraform
    terraform init
    
    print_status "Planning Terraform deployment..."
    terraform plan -out=tfplan
    
    print_status "Applying Terraform configuration..."
    terraform apply tfplan
    
    # Get cluster information
    CLUSTER_NAME=$(terraform output -raw cluster_name)
    CLUSTER_ENDPOINT=$(terraform output -raw cluster_endpoint)
    
    print_status "Configuring kubectl for EKS cluster: $CLUSTER_NAME"
    aws eks update-kubeconfig --region us-west-2 --name "$CLUSTER_NAME"
    
    # Deploy application to EKS
    cd ..
    deploy_k8s
}

# Main deployment logic
case $ENVIRONMENT in
    local)
        deploy_local
        ;;
    k8s)
        deploy_k8s
        ;;
    aws)
        deploy_aws
        ;;
esac

print_success "Deployment completed successfully!"
