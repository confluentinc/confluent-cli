package schema_registry

import (
	"github.com/spf13/cobra"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	v3 "github.com/confluentinc/cli/internal/pkg/config/v3"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/confluentinc/cli/internal/pkg/output"
	srsdk "github.com/confluentinc/schema-registry-sdk-go"
)

type subjectCommand struct {
	*pcmd.AuthenticatedCLICommand
	srClient *srsdk.APIClient
}

// NewSubjectCommand returns the Cobra command for Schema Registry subject list
func NewSubjectCommand(config *v3.Config, prerunner pcmd.PreRunner, srClient *srsdk.APIClient) *cobra.Command {
	cliCmd := pcmd.NewAuthenticatedCLICommand(
		&cobra.Command{
			Use:   "subject",
			Short: "Manage Schema Registry subjects.",
		},
		config, prerunner)
	subjectCmd := &subjectCommand{
		AuthenticatedCLICommand: cliCmd,
		srClient:                srClient,
	}
	subjectCmd.init()
	return subjectCmd.Command
}

func (c *subjectCommand) init() {

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List subjects.",
		Example: FormatDescription(`
Retrieve all subjects available in a Schema Registry

::
		config.CLIName schema-registry subject list
`, c.Config.CLIName),
		RunE: c.list,
		Args: cobra.NoArgs,
	}
	listCmd.Flags().StringP(output.FlagName, output.ShortHandFlag, output.DefaultValue, output.Usage)
	listCmd.Flags().SortFlags = false
	c.AddCommand(listCmd)
	// Update
	updateCmd := &cobra.Command{
		Use:   "update <subjectname> [--compatibility <compatibility>] [--mode <mode>] ",
		Short: "Update subject compatibility or mode.",
		Example: FormatDescription(`
Update subject level compatibility or mode of schema registry.

::
		config.CLIName schema-registry subject update <subjectname> --compatibility=BACKWARD
		config.CLIName schema-registry subject update <subjectname> --mode=READWRITE
`, c.Config.CLIName),
		RunE: c.update,
		Args: cobra.ExactArgs(1),
	}
	updateCmd.Flags().String("compatibility", "", "Can be BACKWARD, BACKWARD_TRANSITIVE, FORWARD, FORWARD_TRANSITIVE, FULL, FULL_TRANSITIVE, or NONE.")
	updateCmd.Flags().String("mode", "", "Can be READWRITE, READ, OR WRITE.")
	updateCmd.Flags().SortFlags = false
	c.AddCommand(updateCmd)

	// Describe
	describeCmd := &cobra.Command{
		Use:   "describe <subjectname>",
		Short: "Describe subject versions and compatibility.",
		Example: FormatDescription(`
Retrieve all versions registered under a given subject and its compatibility level.

::
		config.CLIName schema-registry subject describe <subjectname>
`, c.Config.CLIName),
		RunE: c.describe,
		Args: cobra.ExactArgs(1),
	}
	describeCmd.Flags().StringP(output.FlagName, output.ShortHandFlag, output.DefaultValue, output.Usage)
	describeCmd.Flags().SortFlags = false
	c.AddCommand(describeCmd)
}

func (c *subjectCommand) update(cmd *cobra.Command, args []string) error {
	compat, err := cmd.Flags().GetString("compatibility")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	if compat != "" {
		return c.updateCompatibility(cmd, args)
	}
	mode, err := cmd.Flags().GetString("mode")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	if mode != "" {
		return c.updateMode(cmd, args)
	}
	return errors.New("flag --compatibility or --mode is required.")
}
func (c *subjectCommand) updateCompatibility(cmd *cobra.Command, args []string) error {
	srClient, ctx, err := GetApiClient(cmd, c.srClient, c.Config, c.Version)
	if err != nil {
		return err
	}
	compat, err := cmd.Flags().GetString("compatibility")
	if err != nil {
		return err
	}
	updateReq := srsdk.ConfigUpdateRequest{Compatibility: compat}
	_, _, err = srClient.DefaultApi.UpdateSubjectLevelConfig(ctx, args[0], updateReq)
	if err != nil {
		return err
	}
	pcmd.Println(cmd, "Successfully updated")
	return nil
}

func (c *subjectCommand) updateMode(cmd *cobra.Command, args []string) error {
	mode, err := cmd.Flags().GetString("mode")
	if err != nil {
		return err
	}
	srClient, ctx, err := GetApiClient(cmd, c.srClient, c.Config, c.Version)
	if err != nil {
		return err
	}
	updatedMode, _, err := srClient.DefaultApi.UpdateMode(ctx, args[0], srsdk.ModeUpdateRequest{Mode: mode})
	if err != nil {
		return err
	}
	pcmd.Println(cmd, "Successfully updated Subject level Mode: "+updatedMode.Mode)
	return nil
}

func (c *subjectCommand) list(cmd *cobra.Command, args []string) error {
	type listDisplay struct {
		Subject string
	}
	srClient, ctx, err := GetApiClient(cmd, c.srClient, c.Config, c.Version)
	if err != nil {

		return err
	}
	list, _, err := srClient.DefaultApi.List(ctx)
	if err != nil {
		return err
	}
	if len(list) > 0 {
		outputWriter, err := output.NewListOutputWriter(cmd, []string{"Subject"}, []string{"Subject"}, []string{"subject"})
		if err != nil {
			return errors.HandleCommon(err, cmd)
		}
		for _, l := range list {
			outputWriter.AddElement(&listDisplay{
				Subject: l,
			})
		}
		return outputWriter.Out()
	} else {
		pcmd.Println(cmd, "No subjects")
	}
	return nil
}

func (c *subjectCommand) describe(cmd *cobra.Command, args []string) error {
	srClient, ctx, err := GetApiClient(cmd, c.srClient, c.Config, c.Version)
	if err != nil {
		return err
	}
	versions, _, err := srClient.DefaultApi.ListVersions(ctx, args[0])
	if err != nil {
		return err
	}
	outputOption, err := cmd.Flags().GetString(output.FlagName)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	if outputOption == output.Human.String() {
		PrintVersions(versions)
	} else {
		structuredOutput := &struct {
			Version []int32
		}{
			versions,
		}
		fields := []string{"Version"}
		structuredRenames := map[string]string{"Version": "version"}
		return output.DescribeObject(cmd, structuredOutput, fields, map[string]string{}, structuredRenames)
	}
	return nil
}