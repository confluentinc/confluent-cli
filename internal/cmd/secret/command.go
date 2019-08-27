package secret

import (
	"github.com/spf13/cobra"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/config"
	secret "github.com/confluentinc/cli/internal/pkg/secret"
)

type command struct {
	*cobra.Command
	config *config.Config
	prompt pcmd.Prompt
	resolv pcmd.FlagResolver
	plugin secret.PasswordProtection
}

// New returns the default command object for Password Protection
func New(prerunner pcmd.PreRunner, config *config.Config, prompt pcmd.Prompt, resolv pcmd.FlagResolver, plugin secret.PasswordProtection) *cobra.Command {
	cmd := &command{
		Command: &cobra.Command{
			Use:   "secret",
			Short: "Manage secrets for Confluent Platform.",
		},
		config: config,
		prompt: prompt,
		resolv: resolv,
		plugin: plugin,
	}
	cmd.init()
	return cmd.Command
}

func (c *command) init() {
	c.AddCommand(NewMasterKeyCommand(c.config, c.prompt, c.resolv, c.plugin))
	c.AddCommand(NewFileCommand(c.config, c.prompt, c.resolv, c.plugin))
}