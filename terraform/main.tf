terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# EKS Cluster
resource "aws_eks_cluster" "user_app" {
  name     = var.cluster_name
  role_arn = aws_iam_role.eks_cluster.arn
  version  = var.kubernetes_version

  vpc_config {
    subnet_ids              = aws_subnet.private[*].id
    endpoint_private_access = true
    endpoint_public_access  = true
    public_access_cidrs     = ["0.0.0.0/0"]
  }

  depends_on = [
    aws_iam_role_policy_attachment.eks_cluster_AmazonEKSClusterPolicy,
  ]
}

# EKS Node Group
resource "aws_eks_node_group" "user_app" {
  cluster_name    = aws_eks_cluster.user_app.name
  node_group_name = "user-app-nodes"
  node_role_arn   = aws_iam_role.eks_node.arn
  subnet_ids      = aws_subnet.private[*].id

  scaling_config {
    desired_size = 2
    max_size     = 4
    min_size     = 1
  }

  instance_types = ["t3.medium"]

  depends_on = [
    aws_iam_role_policy_attachment.eks_node_AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.eks_node_AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.eks_node_AmazonEC2ContainerRegistryReadOnly,
  ]
}

# VPC
resource "aws_vpc" "user_app" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "user-app-vpc"
  }
}

# Internet Gateway
resource "aws_internet_gateway" "user_app" {
  vpc_id = aws_vpc.user_app.id

  tags = {
    Name = "user-app-igw"
  }
}

# Public Subnets
resource "aws_subnet" "public" {
  count = 2

  vpc_id                  = aws_vpc.user_app.id
  cidr_block              = "10.0.${count.index + 1}.0/24"
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name = "user-app-public-${count.index + 1}"
    "kubernetes.io/role/elb" = "1"
  }
}

# Private Subnets
resource "aws_subnet" "private" {
  count = 2

  vpc_id            = aws_vpc.user_app.id
  cidr_block        = "10.0.${count.index + 10}.0/24"
  availability_zone = data.aws_availability_zones.available.names[count.index]

  tags = {
    Name = "user-app-private-${count.index + 1}"
    "kubernetes.io/role/internal-elb" = "1"
  }
}

# Route Tables
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.user_app.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.user_app.id
  }

  tags = {
    Name = "user-app-public-rt"
  }
}

resource "aws_route_table" "private" {
  count = 2

  vpc_id = aws_vpc.user_app.id

  tags = {
    Name = "user-app-private-rt-${count.index + 1}"
  }
}

# Route Table Associations
resource "aws_route_table_association" "public" {
  count = 2

  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "private" {
  count = 2

  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private[count.index].id
}

# IAM Roles
resource "aws_iam_role" "eks_cluster" {
  name = "eks-cluster-role"

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "eks.amazonaws.com"
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role" "eks_node" {
  name = "eks-node-role"

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
    Version = "2012-10-17"
  })
}

# IAM Role Policy Attachments
resource "aws_iam_role_policy_attachment" "eks_cluster_AmazonEKSClusterPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.eks_cluster.name
}

resource "aws_iam_role_policy_attachment" "eks_node_AmazonEKSWorkerNodePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.eks_node.name
}

resource "aws_iam_role_policy_attachment" "eks_node_AmazonEKS_CNI_Policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.eks_node.name
}

resource "aws_iam_role_policy_attachment" "eks_node_AmazonEC2ContainerRegistryReadOnly" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.eks_node.name
}

# Data sources
data "aws_availability_zones" "available" {
  state = "available"
}

# RDS Subnet Group
resource "aws_db_subnet_group" "user_app" {
  count = var.use_rds ? 1 : 0
  
  name       = "user-app-db-subnet-group"
  subnet_ids = aws_subnet.private[*].id

  tags = {
    Name = "user-app-db-subnet-group"
  }
}

# RDS Security Group
resource "aws_security_group" "rds" {
  count = var.use_rds ? 1 : 0
  
  name_prefix = "user-app-rds-"
  vpc_id      = aws_vpc.user_app.id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.user_app.cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "user-app-rds-sg"
  }
}

# RDS Instance
resource "aws_db_instance" "user_app" {
  count = var.use_rds ? 1 : 0
  
  identifier = "user-app-db"
  
  engine         = "postgres"
  engine_version = "16.1"
  instance_class = var.db_instance_class
  
  allocated_storage     = var.db_allocated_storage
  max_allocated_storage = 100
  storage_type          = "gp2"
  storage_encrypted      = true
  
  db_name  = "user_app_db"
  username = "postgres"
  password = "your-secure-password-change-in-production"
  
  vpc_security_group_ids = [aws_security_group.rds[0].id]
  db_subnet_group_name   = aws_db_subnet_group.user_app[0].name
  
  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "sun:04:00-sun:05:00"
  
  skip_final_snapshot = true
  deletion_protection = false
  
  tags = {
    Name = "user-app-db"
  }
}

# Kubernetes Provider
provider "kubernetes" {
  host                   = aws_eks_cluster.user_app.endpoint
  cluster_ca_certificate = base64decode(aws_eks_cluster.user_app.certificate_authority[0].data)

  exec {
    api_version = "client.authentication.k8s.io/v1beta1"
    command     = "aws"
    args        = ["eks", "get-token", "--cluster-name", aws_eks_cluster.user_app.name]
  }
}

# Deploy Kubernetes manifests
resource "kubernetes_namespace" "user_app" {
  metadata {
    name = "user-app"
    labels = {
      name = "user-app"
    }
  }
}

resource "kubernetes_config_map" "user_app_config" {
  metadata {
    name      = "user-app-config"
    namespace = kubernetes_namespace.user_app.metadata[0].name
  }

  data = {
    DB_HOST                = "postgres-service"
    DB_PORT                = "5432"
    DB_NAME                = "user_app_db"
    JWT_EXPIRATION         = "24h"
    JWT_REFRESH_EXPIRATION = "168h"
  }
}

resource "kubernetes_secret" "user_app_secrets" {
  metadata {
    name      = "user-app-secrets"
    namespace = kubernetes_namespace.user_app.metadata[0].name
  }

  data = {
    POSTGRES_USER   = base64encode("local")
    POSTGRES_PASSWORD = base64encode("local")
    JWT_SECRET_KEY  = base64encode("your-secret-key-change-in-production")
  }

  type = "Opaque"
}

# PostgreSQL Deployment
resource "kubernetes_deployment" "postgres" {
  metadata {
    name      = "postgres"
    namespace = kubernetes_namespace.user_app.metadata[0].name
    labels = {
      app = "postgres"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "postgres"
      }
    }

    template {
      metadata {
        labels = {
          app = "postgres"
        }
      }

      spec {
        container {
          name  = "postgres"
          image = "postgis/postgis:16-3.4"
          port {
            container_port = 5432
          }

          env {
            name = "POSTGRES_USER"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.user_app_secrets.metadata[0].name
                key  = "POSTGRES_USER"
              }
            }
          }

          env {
            name = "POSTGRES_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.user_app_secrets.metadata[0].name
                key  = "POSTGRES_PASSWORD"
              }
            }
          }

          env {
            name  = "POSTGRES_DB"
            value = "user_app_db"
          }

          env {
            name  = "POSTGRES_HOST_AUTH_METHOD"
            value = "password"
          }

          liveness_probe {
            exec {
              command = ["pg_isready", "-U", "$(POSTGRES_USER)", "-d", "$(POSTGRES_DB)"]
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }

          readiness_probe {
            exec {
              command = ["pg_isready", "-U", "$(POSTGRES_USER)", "-d", "$(POSTGRES_DB)"]
            }
            initial_delay_seconds = 5
            period_seconds        = 5
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "postgres" {
  metadata {
    name      = "postgres-service"
    namespace = kubernetes_namespace.user_app.metadata[0].name
  }

  spec {
    selector = {
      app = "postgres"
    }

    port {
      port        = 5432
      target_port = 5432
    }

    type = "ClusterIP"
  }
}

# Application Deployment
resource "kubernetes_deployment" "user_app" {
  metadata {
    name      = "user-app"
    namespace = kubernetes_namespace.user_app.metadata[0].name
    labels = {
      app = "user-app"
    }
  }

  spec {
    replicas = 2

    selector {
      match_labels = {
        app = "user-app"
      }
    }

    template {
      metadata {
        labels = {
          app = "user-app"
        }
      }

      spec {
        container {
          name  = "user-app"
          image = "user-app:latest"
          port {
            container_port = 8080
          }

          env {
            name = "DB_HOST"
            value_from {
              config_map_key_ref {
                name = kubernetes_config_map.user_app_config.metadata[0].name
                key  = "DB_HOST"
              }
            }
          }

          env {
            name = "DB_PORT"
            value_from {
              config_map_key_ref {
                name = kubernetes_config_map.user_app_config.metadata[0].name
                key  = "DB_PORT"
              }
            }
          }

          env {
            name = "DB_NAME"
            value_from {
              config_map_key_ref {
                name = kubernetes_config_map.user_app_config.metadata[0].name
                key  = "DB_NAME"
              }
            }
          }

          env {
            name = "DB_USER"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.user_app_secrets.metadata[0].name
                key  = "POSTGRES_USER"
              }
            }
          }

          env {
            name = "DB_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.user_app_secrets.metadata[0].name
                key  = "POSTGRES_PASSWORD"
              }
            }
          }

          env {
            name = "JWT_SECRET_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.user_app_secrets.metadata[0].name
                key  = "JWT_SECRET_KEY"
              }
            }
          }

          env {
            name = "JWT_EXPIRATION"
            value_from {
              config_map_key_ref {
                name = kubernetes_config_map.user_app_config.metadata[0].name
                key  = "JWT_EXPIRATION"
              }
            }
          }

          env {
            name = "JWT_REFRESH_EXPIRATION"
            value_from {
              config_map_key_ref {
                name = kubernetes_config_map.user_app_config.metadata[0].name
                key  = "JWT_REFRESH_EXPIRATION"
              }
            }
          }

          liveness_probe {
            http_get {
              path = "/api/v1/health"
              port = 8080
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }

          readiness_probe {
            http_get {
              path = "/api/v1/health"
              port = 8080
            }
            initial_delay_seconds = 5
            period_seconds        = 5
          }

          resources {
            requests = {
              memory = "128Mi"
              cpu    = "100m"
            }
            limits = {
              memory = "512Mi"
              cpu    = "500m"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "user_app" {
  metadata {
    name      = "user-app-service"
    namespace = kubernetes_namespace.user_app.metadata[0].name
  }

  spec {
    selector = {
      app = "user-app"
    }

    port {
      port        = 80
      target_port = 8080
    }

    type = "LoadBalancer"
  }
}

# Outputs
output "cluster_endpoint" {
  value = aws_eks_cluster.user_app.endpoint
}

output "cluster_security_group_id" {
  value = aws_eks_cluster.user_app.vpc_config[0].cluster_security_group_id
}

output "cluster_certificate_authority_data" {
  value = aws_eks_cluster.user_app.certificate_authority[0].data
}
