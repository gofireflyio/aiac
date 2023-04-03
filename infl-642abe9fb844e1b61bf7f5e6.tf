resource "aws_security_group" "k8s-iaccode-iaccoded-b365a797c1-d24" {
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
  name = "k8s-iaccode-iaccoded-b365a797c1"
  tags = {
    "elbv2.k8s.aws/cluster"    = "stag-init"
    "ingress.k8s.aws/resource" = "ManagedLBSecurityGroup"
    "ingress.k8s.aws/stack"    = "iac-code/iac-code-deleter-ingress"
  }
  vpc_id = "vpc-08f6ecf508dd3bf50"
}

