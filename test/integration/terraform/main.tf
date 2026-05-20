terraform {
  required_version = ">= 1.5"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

provider "aws" {
  region = var.region
}

variable "region" {
  default = "us-east-1"
}

variable "prefix" {
  default = "amg-migrator-integ"
}

# IAM role for AMG workspaces
resource "aws_iam_role" "grafana" {
  name = "${var.prefix}-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "grafana.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "grafana_cloudwatch" {
  role       = aws_iam_role.grafana.name
  policy_arn = "arn:aws:iam::aws:policy/CloudWatchReadOnlyAccess"
}

# Source workspace: AMG v9.4
resource "aws_grafana_workspace" "src_v9" {
  name                     = "${var.prefix}-src-v9"
  account_access_type      = "CURRENT_ACCOUNT"
  authentication_providers = ["AWS_SSO"]
  permission_type          = "SERVICE_MANAGED"
  role_arn                 = aws_iam_role.grafana.arn
  grafana_version          = "9.4"
  data_sources             = ["CLOUDWATCH"]
}

resource "aws_grafana_workspace_service_account" "src_v9" {
  name         = "migrator-sa"
  grafana_role = "ADMIN"
  workspace_id = aws_grafana_workspace.src_v9.id
}

# Source workspace: AMG v10.4
resource "aws_grafana_workspace" "src_v10" {
  name                     = "${var.prefix}-src-v10"
  account_access_type      = "CURRENT_ACCOUNT"
  authentication_providers = ["AWS_SSO"]
  permission_type          = "SERVICE_MANAGED"
  role_arn                 = aws_iam_role.grafana.arn
  grafana_version          = "10.4"
  data_sources             = ["CLOUDWATCH"]
}

resource "aws_grafana_workspace_service_account" "src_v10" {
  name         = "migrator-sa"
  grafana_role = "ADMIN"
  workspace_id = aws_grafana_workspace.src_v10.id
}

# Destination workspace: AMG v12.4
resource "aws_grafana_workspace" "dst_v12" {
  name                     = "${var.prefix}-dst-v12"
  account_access_type      = "CURRENT_ACCOUNT"
  authentication_providers = ["AWS_SSO"]
  permission_type          = "SERVICE_MANAGED"
  role_arn                 = aws_iam_role.grafana.arn
  grafana_version          = "12.4"
  data_sources             = ["CLOUDWATCH"]
}

resource "aws_grafana_workspace_service_account" "dst_v12" {
  name         = "migrator-sa"
  grafana_role = "ADMIN"
  workspace_id = aws_grafana_workspace.dst_v12.id
}

# Enable unified alerting on all workspaces via AWS CLI
resource "null_resource" "enable_alerting" {
  depends_on = [
    aws_grafana_workspace.src_v9,
    aws_grafana_workspace.src_v10,
    aws_grafana_workspace.dst_v12,
  ]

  provisioner "local-exec" {
    command = <<-EOT
      for ws in ${aws_grafana_workspace.src_v9.id} ${aws_grafana_workspace.src_v10.id} ${aws_grafana_workspace.dst_v12.id}; do
        aws grafana update-workspace-configuration \
          --workspace-id "$ws" \
          --configuration '{"unifiedAlerting":{"enabled":true}}' \
          --region ${var.region}
        echo "Enabled alerting on $ws"
      done
      # Wait for all workspaces to finish updating
      for ws in ${aws_grafana_workspace.src_v9.id} ${aws_grafana_workspace.src_v10.id} ${aws_grafana_workspace.dst_v12.id}; do
        while [ "$(aws grafana describe-workspace --workspace-id $ws --region ${var.region} --query 'workspace.status' --output text)" != "ACTIVE" ]; do
          sleep 10
        done
        echo "$ws is ACTIVE"
      done
    EOT
  }
}

output "src_v9_endpoint" {
  value = aws_grafana_workspace.src_v9.endpoint
}

output "src_v9_service_account_id" {
  value = aws_grafana_workspace_service_account.src_v9.service_account_id
}

output "src_v10_endpoint" {
  value = aws_grafana_workspace.src_v10.endpoint
}

output "src_v10_service_account_id" {
  value = aws_grafana_workspace_service_account.src_v10.service_account_id
}

output "dst_v12_endpoint" {
  value = aws_grafana_workspace.dst_v12.endpoint
}

output "dst_v12_service_account_id" {
  value = aws_grafana_workspace_service_account.dst_v12.service_account_id
}
