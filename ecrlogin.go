package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

const (
	usage = `ECR Login
Provide the profile to login to ECR.
The region is optional and if not set the ~/.aws/config will be loaded.

ecrlogin [PROFILE] [REGION]
`
	loadConfig = "AWS_SDK_LOAD_CONFIG"

	dockerLogin = "docker login -u %s -p %s %s"
)

func main() {

	var profile, region string

	args := os.Args[1:]

	switch len(args) {
	case 1:
		profile = args[0]
		os.Setenv(loadConfig, "true")
	case 2:
		profile, region = args[0], args[1]
	default:
		fmt.Printf(usage)
		return
	}

	sess, err := AWSSession(loadConfig, profile, region)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	auth, err := GetECRAuth(sess)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	user, password, endpoint, err := GetECRInfo(auth)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	err = ECRLogin(user, password, endpoint)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	fmt.Printf("Login Succeeded.\n\nThe registry endpoint is:\n%s\n",
		endpoint)

}

// AWSSession will create the AWS session by getting the credentials and region
func AWSSession(loadConfig, profile, region string) (sess *session.Session, err error) {

	if os.Getenv(loadConfig) == "true" {
		sess, err = session.NewSession(&aws.Config{
			Credentials: credentials.NewSharedCredentials("", profile),
		})
	} else {
		sess, err = session.NewSession(&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewSharedCredentials("", profile),
		})
	}

	return
}

// GetECRAuth will retrive the ECR authorization data from the AWS session
func GetECRAuth(sess *session.Session) (auth []*ecr.AuthorizationData, err error) {
	svc := ecr.New(sess)
	input := &ecr.GetAuthorizationTokenInput{}
	resp, err := svc.GetAuthorizationToken(input)
	auth = resp.AuthorizationData

	return
}

// GetECRInfo will get the ECR endpoint, username, and password from the ECR authorization data
func GetECRInfo(auth []*ecr.AuthorizationData) (user, password, endpoint string, err error) {
	decode, err := base64.StdEncoding.DecodeString(*auth[0].AuthorizationToken)
	token := strings.SplitN(string(decode), ":", 2)
	user = token[0]
	password = token[1]
	endpoint = *auth[0].ProxyEndpoint

	return
}

// ECRLogin will login to the ECR using Docker and passing username, password, and endpoint
func ECRLogin(user, password, endpoint string) (err error) {
	cmd := fmt.Sprintf(dockerLogin, user, password, endpoint)
	login := exec.Command("bash", "-c", cmd)
	err = login.Run()

	return
}
