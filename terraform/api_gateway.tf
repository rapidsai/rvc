resource "aws_api_gateway_rest_api" "rvc" {
  name = "rvc"
  
  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

# Rapids endpoint resources
resource "aws_api_gateway_resource" "rapids" {
  rest_api_id = aws_api_gateway_rest_api.rvc.id
  parent_id   = aws_api_gateway_rest_api.rvc.root_resource_id
  path_part   = "rapids"
}

resource "aws_api_gateway_resource" "rapids_version" {
  rest_api_id = aws_api_gateway_rest_api.rvc.id
  parent_id   = aws_api_gateway_resource.rapids.id
  path_part   = "{version}"
}

# UCX-Py endpoint resources
resource "aws_api_gateway_resource" "ucx_py" {
  rest_api_id = aws_api_gateway_rest_api.rvc.id
  parent_id   = aws_api_gateway_rest_api.rvc.root_resource_id
  path_part   = "ucx-py"
}

resource "aws_api_gateway_resource" "ucx_py_version" {
  rest_api_id = aws_api_gateway_rest_api.rvc.id
  parent_id   = aws_api_gateway_resource.ucx_py.id
  path_part   = "{version}"
}

# Method and integration for Rapids
resource "aws_api_gateway_method" "rapids_get" {
  rest_api_id   = aws_api_gateway_rest_api.rvc.id
  resource_id   = aws_api_gateway_resource.rapids_version.id
  http_method   = "GET"
  authorization = "NONE"

  request_parameters = {
    "method.request.path.version" = true
  }
}

resource "aws_api_gateway_integration" "rapids_lambda" {
  rest_api_id = aws_api_gateway_rest_api.rvc.id
  resource_id = aws_api_gateway_resource.rapids_version.id
  http_method = aws_api_gateway_method.rapids_get.http_method
  
  integration_http_method = "POST"
  type                   = "AWS_PROXY"
  uri                    = aws_lambda_function.rvc_rapids.invoke_arn
  credentials            = aws_iam_role.api_gateway_executor.arn
}

# Method and integration for UCX-Py
resource "aws_api_gateway_method" "ucx_py_get" {
  rest_api_id   = aws_api_gateway_rest_api.rvc.id
  resource_id   = aws_api_gateway_resource.ucx_py_version.id
  http_method   = "GET"
  authorization = "NONE"

  request_parameters = {
    "method.request.path.version" = true
  }
}

resource "aws_api_gateway_integration" "ucx_py_lambda" {
  rest_api_id = aws_api_gateway_rest_api.rvc.id
  resource_id = aws_api_gateway_resource.ucx_py_version.id
  http_method = aws_api_gateway_method.ucx_py_get.http_method
  
  integration_http_method = "POST"
  type                   = "AWS_PROXY"
  uri                    = aws_lambda_function.rvc_ucx_py.invoke_arn
  credentials            = aws_iam_role.api_gateway_executor.arn
}

# Domain name and mapping
resource "aws_api_gateway_domain_name" "rvc" {
  domain_name     = var.domain_name
  regional_certificate_arn = data.aws_acm_certificate.domain_cert.arn
  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_base_path_mapping" "rvc" {
  api_id      = aws_api_gateway_rest_api.rvc.id
  stage_name  = aws_api_gateway_stage.rvc.stage_name
  domain_name = aws_api_gateway_domain_name.rvc.domain_name
}

# Deployment and stage
resource "aws_api_gateway_deployment" "rvc" {
  rest_api_id = aws_api_gateway_rest_api.rvc.id
  
  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.rapids.id,
      aws_api_gateway_resource.ucx_py.id,
      aws_api_gateway_method.rapids_get.id,
      aws_api_gateway_method.ucx_py_get.id,
      aws_api_gateway_integration.rapids_lambda.id,
      aws_api_gateway_integration.ucx_py_lambda.id
    ]))
  }
}

resource "aws_api_gateway_stage" "rvc" {
  deployment_id = aws_api_gateway_deployment.rvc.id
  rest_api_id  = aws_api_gateway_rest_api.rvc.id
  stage_name   = "prod"
}

data "aws_acm_certificate" "domain_cert" {
  domain      = var.certificate_name
  statuses    = ["ISSUED"]
  most_recent = true
}

resource "aws_route53_record" "domain" {
  name    = var.domain_name
  type    = "A"
  zone_id = data.aws_route53_zone.domain.zone_id

  alias {
    name                   = aws_api_gateway_domain_name.rvc.regional_domain_name
    zone_id                = aws_api_gateway_domain_name.rvc.regional_zone_id
    evaluate_target_health = false
  }
}

data "aws_route53_zone" "domain" {
  name = "gpuci.io"
}
