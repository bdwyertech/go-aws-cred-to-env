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
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var disableSharedConfig bool

func init() {
	if flag.Lookup("disable-shared-config") == nil {
		flag.BoolVar(&disableSharedConfig, "disable-shared-config", false, "Disable Shared Configuration (force use of EC2/ECS metadata, ignore AWS_PROFILE, etc.)")
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

	fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", creds.AccessKeyID)
	fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", creds.SecretAccessKey)
	fmt.Printf("export AWS_SESSION_TOKEN=%s\n", creds.SessionToken)
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
	}

	sess := session.Must(session.NewSessionWithOptions(sess_opts))

	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		log.Fatal(err)
	}

	return
}
