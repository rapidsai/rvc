resource "aws_cloudwatch_log_group" "rapids_logs" {
  name              = "/aws/lambda/rvc-rapids"
  retention_in_days = 30
}

resource "aws_cloudwatch_log_group" "ucx_py_logs" {
  name              = "/aws/lambda/rvc-ucx-py"
  retention_in_days = 30
}
