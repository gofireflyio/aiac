resource "aws_lambda_function" "lambda_promtail-e04" {
  architectures = ["x86_64"]
  environment {
    variables = {
      BATCH_SIZE    = ""
      EXTRA_LABELS  = ""
      KEEP_STREAM   = "true"
      PASSWORD      = ""
      TENANT_ID     = ""
      USERNAME      = ""
      WRITE_ADDRESS = "https://loki.prod.infralight.cloud/loki/api/v1/push"
    }
  }
  function_name = "lambda_promtail"
  image_uri     = "094724549126.dkr.ecr.us-east-1.amazonaws.com/promtail_lambda:000006"
  package_type  = "Image"
  role          = "arn:aws:iam::094724549126:role/iam_for_lambda"
  timeout       = 60
  tracing_config {
    mode = "PassThrough"
  }
}

