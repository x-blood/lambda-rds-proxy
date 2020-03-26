package main


import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var dbHost = os.Getenv("DB_HOST_PROXY")
var dbUser = os.Getenv("DB_USER")
var dbPass = os.Getenv("DB_PASS")
var dbName = os.Getenv("DB_NAME")
var dbSource = dbUser + ":" + dbPass + "@tcp(" + dbHost + ":3306)/" + dbName

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	var buf bytes.Buffer

	cnn, err := sql.Open("mysql", dbSource)
	if err != nil {
		return Response{StatusCode: 500}, err
	}
	defer cnn.Close()

	id := int(1)
	val := ""

	err =
		cnn.QueryRow(
			"SELECT * FROM test WHERE id = ?", id).Scan(&id, &val)
	if err != nil {
		fmt.Println(err)
		return Response{StatusCode: 500}, err
	}

	resultMessage := fmt.Sprintf("Hello Success!! id : %d, val : %s", id, val)

	body, err := json.Marshal(map[string]interface{}{
		"message": resultMessage,
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
