package main

import (
	"embed"

	"github.com/maxgio92/kamajictl/cmd"
)

//go:embed manifests/*.yaml
var embeddedManifests embed.FS

func main() {
	cmd.Execute(&embeddedManifests)
}
