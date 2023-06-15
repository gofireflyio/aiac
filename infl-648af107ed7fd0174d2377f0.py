from aws_cdk import (
    aws_s3 as s3,
    aws_sqs as sqs,
    aws_sns as sns,
    aws_sns_subscriptions as subs,
    core
)

app = core.App()

# S3 Bucket
s3_bucket = s3.Bucket(
    app, "S3Bucket",
    bucket_name="prod-fetched-resources",
    versioned=True,
    encryption=s3.BucketEncryption.KMS_MANAGED,
    block_public_access=s3.BlockPublicAccess.BLOCK_ALL,
    removal_policy=core.RemovalPolicy.RETAIN
)

# SQS Queue
sqs_queue = sqs.Queue(
    app, "SQSQueue",
    queue_name="prod-iac-ci-worker-sqs",
    encryption=sqs.QueueEncryption.KMS_MANAGED,
    retention_period=core.Duration.days(2),
    visibility_timeout=core.Duration.hours(12),
    content_based_deduplication=False,
    fifo_queue=False,
    max_message_size=262144,
    receive_wait_time=core.Duration.seconds(0),
    dead_letter_queue=sqs.DeadLetterQueue(
        max_receive_count=3,
        queue=sqs.Queue(
            app, "DeadLetterQueue",
            queue_name="prod-iac-ci-worker-sqs-dl"
        )
    ),
    policy=sqs.QueuePolicy(
        app, "QueuePolicy",
        queues=[sqs_queue],
        policy_document={
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {"Service": "s3.amazonaws.com"},
                    "Action": "sqs:SendMessage",
                    "Resource": sqs_queue.queue_arn,
                    "Condition": {
                        "StringEquals": {"aws:SourceAccount": "094724549126"},
                        "ArnLike": {"aws:SourceArn": s3_bucket.bucket_arn}
                    }
                }
            ]
        }
    ),
    tags={
        "environment": "prod",
        "project": "flywheel"
    }
)

# SNS Topic
sns_topic = sns.Topic(
    app, "SNSTopic",
    topic_name="prod-iac-ci-worker-sns",
    display_name="Prod IAC CI Worker SNS",
    encryption=sns.TopicEncryption.KMS_MANAGED,
    master_key=sqs_queue.encryption_master_key,
    fifo=False,
    content_based_deduplication=False,
    delivery_policy=sns.DeliveryPolicy(
        max_retry_delay=core.Duration.seconds(300),
        http_retry_delay=core.Duration.seconds(5),
        num_retries=3
    ),
    tags={
        "environment": "prod",
        "project": "flywheel"
    }
)

# SNS Subscription
sns_subscription = subs.SqsSubscription(sqs_queue)

sns_topic.add_subscription(sns_subscription)

app.synth()