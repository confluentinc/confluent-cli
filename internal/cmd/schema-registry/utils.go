package schema_registry

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/confluentinc/go-printer"
	srsdk "github.com/confluentinc/schema-registry-sdk-go"

	"github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/version"
)

const (
	SubjectUsage = "Subject of the schema."
)

func GetApiClient(cmd *cobra.Command, srClient *srsdk.APIClient, cfg *cmd.DynamicConfig, ver *version.Version) (*srsdk.APIClient, context.Context, error) {
	if srClient != nil {
		// Tests/mocks
		return srClient, nil, nil
	}
	srClient, ctx, err := SchemaRegistryClient(cmd, cfg, ver)
	if err != nil {
		return nil, nil, err
	}
	return srClient, ctx, nil
}

func PrintVersions(versions []int32) {
	titleRow := []string{"Version"}
	var entries [][]string
	for _, version := range versions {
		record := &struct{ Version int32 }{version}
		entries = append(entries, printer.ToRow(record, titleRow))
	}
	printer.RenderCollectionTable(entries, titleRow)
}

func RequireSubjectFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("subject", "S", "", SubjectUsage)
	_ = cmd.MarkFlagRequired("subject")
}

func getServiceProviderFromUrl(url string) string {
	if url == "" {
		return ""
	}
	//Endpoint url is of the form https://psrc-<id>.<location>.<service-provider>.<devel/stag/prod/env>.cpdev.cloud
	stringSlice := strings.Split(url, ".")
	if len(stringSlice) != 6 {
		return ""
	}
	return strings.Trim(stringSlice[2], ".")
}

func FormatDescription(description string, cliName string) string {
	return strings.ReplaceAll(description, "{{.CLIName}}", cliName)
}