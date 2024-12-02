resource "aws_cloudwatch_log_group" "rapids_logs" {
  name              = "/aws/lambda/rvc-${var.environment}-rapids"
  retention_in_days = 30
}

resource "aws_cloudwatch_log_group" "ucx_py_logs" {
  name              = "/aws/lambda/rvc-${var.environment}-ucx-py"
  retention_in_days = 30
}
