module "aws-lambda-function" {
  source = "github.com/nissim-infra/engine/devops/modules/aws-lambda-function"

  create_event_source = true
  event_source_arn    = "${data.aws_sqs_queue.stag-firefly-engine-queue-e70.arn}"
}

data "aws_sqs_queue" "stag-firefly-engine-queue-e70" {
  name = "stag-firefly-engine-queue"
}

