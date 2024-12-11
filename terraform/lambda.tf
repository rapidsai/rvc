resource "aws_lambda_function" "rvc_rapids" {
  filename         = "../bin/rvc_serverless"
  function_name    = "rvc-rapids"
  role            = aws_iam_role.lambda_role.arn
  handler         = "bin/rvc_serverless"
  source_code_hash = filebase64sha256("../bin/rvc_serverless")
  runtime         = "go1.x"
  memory_size     = 1024
  timeout         = 30

  tags = {
    Environment = "prod"
    Service     = "rvc"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_lambda_function" "rvc_ucx_py" {
  filename         = "../bin/rvc_serverless"
  function_name    = "rvc-ucx-py"
  role            = aws_iam_role.lambda_role.arn
  handler         = "bin/rvc_serverless"
  source_code_hash = filebase64sha256("../bin/rvc_serverless")
  runtime         = "go1.x"
  memory_size     = 1024
  timeout         = 30

  tags = {
    Environment = "prod"
    Service     = "rvc"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_lambda_permission" "apigw_rapids" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.rvc_rapids.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.rvc.execution_arn}/*/*"
}

resource "aws_lambda_permission" "apigw_ucx_py" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.rvc_ucx_py.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.rvc.execution_arn}/*/*"
}
