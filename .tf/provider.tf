locals {
  endpoint = "http://localstack:4566"
}

provider "aws" {
  access_key                  = "mock_access_key"
  region                      = "us-east-1"
  s3_force_path_style         = true
  secret_key                  = "mock_secret_key"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    apigateway     = local.endpoint
    cloudformation = local.endpoint
    cloudwatch     = local.endpoint
    dynamodb       = local.endpoint
    es             = local.endpoint
    firehose       = local.endpoint
    iam            = local.endpoint
    kinesis        = local.endpoint
    lambda         = local.endpoint
    route53        = local.endpoint
    redshift       = local.endpoint
    s3             = local.endpoint
    secretsmanager = local.endpoint
    ses            = local.endpoint
    sns            = local.endpoint
    sqs            = local.endpoint
    ssm            = local.endpoint
    stepfunctions  = local.endpoint
    sts            = local.endpoint
  }
}
