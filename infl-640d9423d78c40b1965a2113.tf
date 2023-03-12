resource "aws_launch_template" "Karpenter-stag-init-6049714188429450335-539" {
  block_device_mappings {
    device_name = "/dev/xvda"
    ebs {
      encrypted   = "true"
      volume_size = 20
      volume_type = "gp3"
    }
  }
  default_version         = 1
  disable_api_termination = false
  iam_instance_profile {
    name = "stag-init-KarpenterNodeInstanceProfile"
  }
  image_id = "ami-0eba34316a915ee9f"
  metadata_options {
    http_put_response_hop_limit = 2
    http_tokens                 = "required"
  }
  name = "Karpenter-stag-init-6049714188429450335"
  tag_specifications {
    resource_type = "network-interface"
    tags = {
      Name                            = "karpenter.sh/provisioner-name/general-jobs"
      "karpenter.sh/provisioner-name" = "general-jobs"
    }
  }
  tags = {
    Name                            = "karpenter.sh/provisioner-name/general-jobs"
    "karpenter.k8s.aws/cluster"     = "stag-init"
    "karpenter.sh/provisioner-name" = "general-jobs"
  }
  user_data              = "REDACTED-BY-FIREFLY:aaa1a3fc3f8c71ed3639150b79ee8a0e0bba160666833a1ba457a2e3377169fc:sha256"
  vpc_security_group_ids = ["sg-04f87a0508fa463c0"]
}

