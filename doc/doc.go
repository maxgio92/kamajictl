package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra/doc"
	"k8s.io/kubectl/pkg/cmd/util"

	"github.com/maxgio92/kamajicli/cmd"
	"github.com/maxgio92/kamajicli/internal/options"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := doc.GenMarkdownTree(cmd.NewRootCommand(options.NewCommonOptions(
		options.WithLogger(log.New()),
	)), "./doc"); err != nil {
		fmt.Println(err)
		os.Exit(util.DefaultErrorExitCode)
	}
}
