package awsutil

import (
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Session is a mock session which is used to hit the mock server
var MockSession = func() *session.Session {
	// server is the mock server that simply writes a 200 status back to the client
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	return session.Must(session.NewSession(&aws.Config{
		DisableSSL:       aws.Bool(true),
		Endpoint:         aws.String(server.URL),
		Region:           aws.String("us-west-2"),
		Credentials:      credentials.NewStaticCredentials(*aws.String("mock"), *aws.String("mock"), *aws.String("token")),
		S3ForcePathStyle: aws.Bool(true),
	}))
}()

func NewAwsMockService() (*AWSService, error) {

	// c := Session.ClientConfig("Mock", cfgs...)
	svc := s3.New(MockSession)
	return &AWSService{client: svc, bucket: "mock"}, nil

}
