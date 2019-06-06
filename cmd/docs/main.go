package main

import (
	"fmt"
	"path"

	"github.com/confluentinc/cli/internal/cmd"
	"github.com/confluentinc/cli/internal/pkg/config"
	"github.com/confluentinc/cli/internal/pkg/doc"
	"github.com/confluentinc/cli/internal/pkg/log"
	"github.com/confluentinc/cli/internal/pkg/version"
)

var (
	// Injected from linker flags like `go build -ldflags "-X main.cliName=$NAME"`
	cliName = "confluent"
)

// See https://github.com/spf13/cobra/blob/master/doc/rest_docs.md
func main() {
	emptyStr := func(filename string) string { return "" }
	sphinxRef := func(name, ref string) string { return fmt.Sprintf(":ref:`%s`", ref) }
	confluent, err := cmd.NewConfluentCommand(cliName, &config.Config{}, &version.Version{}, log.New())
	if err != nil {
		panic(err)
	}
	err = doc.GenReSTTreeCustom(confluent, path.Join(".", "docs", cliName), emptyStr, sphinxRef)
	if err != nil {
		panic(err)
	}

	indexHeader := func(filename string) string {
		return `.. _ccloud-ref:

|ccloud| CLI Command Reference
==============================

The available |ccloud| CLI commands are documented here.

`
	}
	err = doc.GenReSTIndex(confluent, path.Join(".", "docs", cliName, "index.rst"), indexHeader, sphinxRef)
	if err != nil {
		panic(err)
	}
}