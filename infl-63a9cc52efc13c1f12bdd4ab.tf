resource "aws_iam_role" "PileusRole-ab8" {
  assume_role_policy = jsonencode({
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::932213950603:root"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "REDACTED-BY-FIREFLY:484139c2f3315d016544750e0f1d2a5bc5f046502eaa611707cc4e6c67cd235a:sha256"
        }
      }
    }
  ]
  })
  inline_policy {
    name   = "PileusPolicy"
    policy = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Action\":\"ec2:Describe*\",\"Resource\":\"*\",\"Effect\":\"Allow\"},{\"Action\":\"organizations:ListAccounts\",\"Resource\":\"*\",\"Effect\":\"Allow\"},{\"Action\":\"elasticloadbalancing:Describe*\",\"Resource\":\"*\",\"Effect\":\"Allow\"},{\"Action\":[\"s3:ListBucket\",\"s3:GetBucketLocation\",\"s3:ListBucketVersions\",\"s3:GetBucketVersioning\",\"s3:GetLifecycleConfiguration\",\"s3:GetEncryptionConfiguration\",\"s3:ListAllMyBuckets\",\"s3:ListBucketMultipartUploads\",\"s3:ListMultipartUploadParts\"],\"Resource\":\"*\",\"Effect\":\"Allow\"},{\"Action\":[\"cloudwatch:ListMetrics\",\"cloudwatch:GetMetricStatistics\",\"cloudwatch:GetMetricData\",\"logs:DescribeLogGroups\",\"logs:GetQueryResults\"],\"Resource\":\"*\",\"Effect\":\"Allow\"},{\"Action\":[\"logs:StartQuery\"],\"Resource\":[\"arn:aws:logs:*:*:log-group:/aws/containerinsights/*/performance\",\"arn:aws:logs:*:*:log-group:/aws/containerinsights/*/performance:*\",\"arn:aws:logs:*:*:log-group:/aws/containerinsights/*/performance:*:*\"],\"Effect\":\"Allow\"},{\"Action\":\"autoscaling:Describe*\",\"Resource\":\"*\",\"Effect\":\"Allow\"},{\"Action\":[\"eks:ListFargateProfiles\",\"eks:DescribeNodegroup\",\"eks:ListNodegroups\",\"eks:DescribeFargateProfile\",\"eks:ListTagsForResource\",\"eks:ListUpdates\",\"eks:DescribeUpdate\",\"eks:DescribeCluster\",\"eks:ListClusters\"],\"Resource\":\"*\",\"Effect\":\"Allow\"},{\"Action\":[\"dynamodb:Describe*\",\"dynamodb:List*\",\"tag:GetResources\",\"rds:DescribeDBInstances\",\"rds:DescribeDBClusters\",\"rds:ListTagsForResource\",\"ecs:DescribeClusters\",\"redshift:DescribeClusters\",\"es:ListDomainNames\",\"es:DescribeElasticsearchDomains\",\"elasticache:DescribeCacheClusters\",\"kinesis:ListStreams\",\"kinesis:DescribeStream\",\"cloudTrail:DescribeTrails\",\"kms:ListKeys\",\"kms:DescribeKey\",\"kms:ListResourceTags\"],\"Resource\":\"*\",\"Effect\":\"Allow\"},{\"Action\":[\"ce:GetRightsizingRecommendation\",\"ce:GetReservationUtilization\",\"ce:GetSavingsPlansUtilizationDetails\",\"ce:GetSavingsPlansUtilization\",\"ce:GetSavingsPlansCoverage\",\"ce:GetTags\",\"ce:GetCostAndUsage\",\"aws-portal:ViewBilling\"],\"Resource\":\"*\",\"Effect\":\"Allow\"}]}"
  }
  name = "PileusRole"
  tags = {
    AssetType = "CloudOpsTools"
    Project   = "Pileus"
    Team      = "CloudOps"
  }
}




resource "aws_ssm_parameter" "second-61f" {
  name  = "second"
  type  = "String"
  value = "REDACTED-BY-FIREFLY:a441b15fe9a3cf56661190a0b93b9dec7d04127288cc87250967cf3b52894d11:sha256"
}




resource "aws_ssm_parameter" "first-52b" {
  name  = "first"
  type  = "String"
  value = "REDACTED-BY-FIREFLY:5e61da5a7e6f3ce4c8f7e91859629e6deae61f9e71d36c7a2e18af9d746cd2a3:sha256"
}

