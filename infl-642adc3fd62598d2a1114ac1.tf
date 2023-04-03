resource "aws_security_group" "k8s-iaccode-iaccoded-c5c33bcaed-052" {
  description = "[k8s] Managed SecurityGroup for LoadBalancer"
  egress {
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
  }
  ingress {
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 443
    protocol    = "tcp"
    to_port     = 443
  }
  name = "k8s-iaccode-iaccoded-c5c33bcaed"
  tags = {
    "elbv2.k8s.aws/cluster"    = "dev-init"
    "ingress.k8s.aws/resource" = "ManagedLBSecurityGroup"
    "ingress.k8s.aws/stack"    = "iac-code/iac-code-deleter-ingress"
  }
  vpc_id = "vpc-0bb7435a6a7f15e26"
}




resource "aws_elasticache_parameter_group" "defaultredis26-c52" {
  description = "Default parameter group for redis2.6"
  family      = "redis2.6"
  name        = "default.redis2.6"
}

