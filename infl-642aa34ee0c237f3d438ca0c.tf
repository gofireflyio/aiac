resource "aws_launch_template" "Karpenter-prod-external-services-8413523120259111216-2b1" {
  block_device_mappings {
    device_name = "/dev/xvda"
    ebs {
      delete_on_termination = "true"
      encrypted             = "false"
      iops                  = 8000
      throughput            = 125
      volume_size           = 108
      volume_type           = "gp3"
    }
  }
  default_version         = 1
  disable_api_termination = false
  iam_instance_profile {
    name = "KarpenterNodeInstanceProfile-prod-external-services"
  }
  image_id = "ami-0ab1687f626a22bf5"
  metadata_options {
    http_put_response_hop_limit = 2
    http_tokens                 = "required"
  }
  name = "Karpenter-prod-external-services-8413523120259111216"
  tag_specifications {
    resource_type = "network-interface"
    tags = {
      Name                            = "karpenter.sh/provisioner-name/general-jobs"
      "karpenter.sh/provisioner-name" = "general-jobs"
    }
  }
  tags = {
    Name                            = "karpenter.sh/provisioner-name/general-jobs"
    "karpenter.k8s.aws/cluster"     = "prod-external-services"
    "karpenter.sh/provisioner-name" = "general-jobs"
  }
  user_data              = "REDACTED-BY-FIREFLY:73e3d1a543898b8738ba384276245fb4760045b075b8d12e274cb34204fb0bf3:sha256"
  vpc_security_group_ids = ["sg-073a49d3f7c7dc971"]
}

