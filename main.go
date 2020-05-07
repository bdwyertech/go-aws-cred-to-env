// Encoding: UTF-8
//
// AWS Credential to Environment
//
// Copyright © 2020 Brian Dwyer - Intelligent Digital Services
//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

var disableSharedConfig, displayCallerIdentity bool

func init() {
	if flag.Lookup("disable-shared-config") == nil {
		flag.BoolVar(&disableSharedConfig, "disable-shared-config", false, "Disable Shared Configuration (force use of EC2/ECS metadata, ignore AWS_PROFILE, etc.)")
	}
	if flag.Lookup("display-caller-identity") == nil {
		flag.BoolVar(&displayCallerIdentity, "display-caller-identity", false, "Display STS Get-Caller-Identity Output")
	}
}

func main() {
	// Parse Flags
	flag.Parse()

	if versionFlag {
		showVersion()
		os.Exit(0)
	}

	// Handler to work around Hashicorp aws-sdk-go-base issues...
	// https://github.com/hashicorp/aws-sdk-go-base/pull/20
	// export AWS_CRED_CONTAINER_RELATIVE_URI=$AWS_CONTAINER_CREDENTIALS_RELATIVE_URI
	if containerUri := os.Getenv("AWS_CRED_CONTAINER_RELATIVE_URI"); containerUri != "" {
		os.Setenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", containerUri)
	}

	creds := getCredentials()

	if runtime.GOOS == "windows" {
		fmt.Printf("$env:AWS_ACCESS_KEY_ID='%s'\n", creds.AccessKeyID)
		fmt.Printf("$env:AWS_SECRET_ACCESS_KEY='%s'\n", creds.SecretAccessKey)
		fmt.Printf("$env:AWS_SESSION_TOKEN='%s'\n", creds.SessionToken)
	} else {
		fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", creds.AccessKeyID)
		fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", creds.SecretAccessKey)
		fmt.Printf("export AWS_SESSION_TOKEN=%s\n", creds.SessionToken)
	}
}

func getCredentials() (creds credentials.Value) {
	// AWS Session
	sess_opts := session.Options{
		// Config:            *aws.NewConfig().WithRegion("us-east-1"),
		Config:            *aws.NewConfig().WithCredentialsChainVerboseErrors(true),
		SharedConfigState: session.SharedConfigEnable,
	}

	if disableSharedConfig {
		sess_opts.SharedConfigState = session.SharedConfigDisable
		os.Unsetenv("AWS_ACCESS_KEY_ID")
		os.Unsetenv("AWS_PROFILE")
	}

	sess := session.Must(session.NewSessionWithOptions(sess_opts))

	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		log.Fatal(err)
	}

	// Validate the Credentials
	svc := sts.New(sess)
	result, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal(err)
	}
	if displayCallerIdentity {
		log.Print("STS Get-Caller-Identity:\n", result)
	}

	return
}
