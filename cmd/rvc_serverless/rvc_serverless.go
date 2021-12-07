package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/rapidsai/rvc/pkg/rvc"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	path := request.Path
	version := request.PathParameters["version"]
	var matchingVersion string
	var err error

	if strings.HasPrefix(path, "/rapids/") {
		matchingVersion, err = rvc.GetUcxPyFromRapids(version)
	} else if strings.HasPrefix(path, "/ucx-py/") {
		matchingVersion, err = rvc.GetRapidsFromUcxPy(version)
	} else {
		errorMsg := fmt.Sprintf("Unexpected path \"%v\"", path)
		return events.APIGatewayProxyResponse{Body: errorMsg, StatusCode: 400}, nil
	}

	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	return events.APIGatewayProxyResponse{Body: matchingVersion, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
