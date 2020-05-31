package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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

	var (
		profile, region string

		sess *session.Session

		err error
	)

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

	if os.Getenv(loadConfig) == "true" {
		sess, err = session.NewSession(&aws.Config{
			Credentials: credentials.NewSharedCredentials("", profile),
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		sess, err = session.NewSession(&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewSharedCredentials("", profile),
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	svc := ecr.New(sess)
	input := &ecr.GetAuthorizationTokenInput{}

	resp, err := svc.GetAuthorizationToken(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				fmt.Println(ecr.ErrCodeServerException, aerr.Error())
			case ecr.ErrCodeInvalidParameterException:
				fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(aerr.Error())
		}

	}

	auth := resp.AuthorizationData
	decode, err := base64.StdEncoding.DecodeString(*auth[0].AuthorizationToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	token := strings.SplitN(string(decode), ":", 2)
	user := token[0]
	password := token[1]
	endpoint := *auth[0].ProxyEndpoint

	cmd := fmt.Sprintf(dockerLogin, user, password, endpoint)
	login := exec.Command("bash", "-c", cmd)
	loginErr := login.Run()
	if loginErr != nil {
		fmt.Println(loginErr)
		return
	}

	fmt.Printf("Login Succeeded.\n\nThe registry endpoint is:\n%s\n",
		endpoint)

}
