output "cluster_name" {
  description = "EKS cluster name"
  value       = aws_eks_cluster.user_app.name
}

output "cluster_endpoint" {
  description = "Endpoint for EKS control plane"
  value       = aws_eks_cluster.user_app.endpoint
}

output "cluster_security_group_id" {
  description = "Security group ids attached to the cluster control plane"
  value       = aws_eks_cluster.user_app.vpc_config[0].cluster_security_group_id
}

output "cluster_certificate_authority_data" {
  description = "Base64 encoded certificate data required to communicate with the cluster"
  value       = aws_eks_cluster.user_app.certificate_authority[0].data
}

output "cluster_arn" {
  description = "The Amazon Resource Name (ARN) of the cluster"
  value       = aws_eks_cluster.user_app.arn
}

output "cluster_version" {
  description = "The Kubernetes server version for the EKS cluster"
  value       = aws_eks_cluster.user_app.version
}
