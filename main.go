package main

import (
	"context"
	"log"
	"os"

	"github.com/yaninyzwitty/aws-resource-auditor-go/cmd"
)

func main() {
	if err := cmd.NewCliCommand().Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
