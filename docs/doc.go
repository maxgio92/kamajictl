package main

import (
	"fmt"
	"github.com/spf13/cobra/doc"
	"os"
	"path"
	"strings"

	"k8s.io/kubectl/pkg/cmd/util"

	"github.com/maxgio92/kamajictl/cmd"
	"github.com/maxgio92/kamajictl/internal/options"
	"github.com/maxgio92/kamajictl/internal/output"
)

const (
	cmdline      = "kamajictl"
	docsDir      = "docs"
	fileTemplate = `---
title: %s
---	

`
)

var (
	filePrepender = func(filename string) string {
		title := strings.TrimPrefix(strings.TrimSuffix(strings.ReplaceAll(filename, "_", " "), ".md"), fmt.Sprintf("%s/", docsDir))
		return fmt.Sprintf(fileTemplate, title)
	}
	linkHandler = func(filename string) string {
		if filename == cmdline+".md" {
			return "_index.md"
		}
		return filename
	}
)

func main() {
	if err := doc.GenMarkdownTreeCustom(
		cmd.NewRootCommand(options.NewCommonOptions(options.WithLogger(output.NewPrinter()))),
		docsDir,
		filePrepender,
		linkHandler,
	); err != nil {
		fmt.Println(err)
		os.Exit(util.DefaultErrorExitCode)
	}

	err := os.Rename(path.Join(docsDir, cmdline+".md"), path.Join(docsDir, "_index.md"))
	if err != nil {
		fmt.Println(err)
		os.Exit(util.DefaultErrorExitCode)
	}
}
