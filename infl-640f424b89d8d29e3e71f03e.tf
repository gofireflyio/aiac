resource "aws_iam_role" "exampleCI1-ff5" {
  assume_role_policy  = jsonencode({
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
})
  managed_policy_arns = ["arn:aws:iam::aws:policy/AmazonS3FullAccess"]
  name                = "exampleCI1"
}

