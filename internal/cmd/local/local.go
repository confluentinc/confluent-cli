package local

import (
	"os"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/errors"

	"github.com/spf13/cobra"
)

const longDescription = `You can test Confluent Platform by running a single-node instance locally on
your laptop or desktop. THESE LOCAL COMMANDS ARE NOT INTENDED FOR PRODUCTION SETUP.

The CLI "local" commands help you manage and interact with this installation
for exploring, testing, experimenting, and otherwise familiarizing yourself
with Confluent Platform.

LOCAL COMMANDS ARE NOT INTENDED TO SETUP OR MANAGE CONFLUENT PLATFORM IN PRODUCTION.
`

type command struct {
	*cobra.Command
	shell ShellRunner
}

// New returns the Cobra command for `local`.
func New(prerunner pcmd.PreRunner, shell ShellRunner) *cobra.Command {
	localCmd := &command{
		Command: &cobra.Command{
			Use:               "local",
			Short:             "Manage local Confluent Platform development environment",
			Long:              longDescription,
			Args:              cobra.ArbitraryArgs,
			PersistentPreRunE: prerunner.Anonymous(),
		},
		shell: shell,
	}
	localCmd.Command.RunE = localCmd.run
	// possibly we should make this an arg and/or move it to env var
	localCmd.Flags().String("path", "", "Path to Confluent Platform install directory")
	_ = localCmd.MarkFlagRequired("path")
	localCmd.Flags().SortFlags = false
	return localCmd.Command
}

func (c *command) run(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	c.shell.Init(os.Stdout, os.Stderr)
	c.shell.Export("CONFLUENT_HOME", path)
	err = c.shell.Source("cp_cli/confluent.sh", Asset)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	_, err = c.shell.Run("main", args)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	return nil
}