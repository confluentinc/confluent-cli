package kafka

import (
	"github.com/confluentinc/cli/shared"
	"github.com/spf13/cobra"
)

type Command struct {
	*cobra.Command
	config *shared.Config
}

func New(config *shared.Config) *cobra.Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:   "kafka",
			Short: "Manage kafka clusters.",
		},
		config: config,
	}
	cmd.init()
	return cmd.Command
}

func (c *Command) init() {
	c.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a Kafka cluster.",
		RunE:  c.create,
	})
	c.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List Kafka clusters.",
		RunE:  c.list,
	})
	c.AddCommand(&cobra.Command{
		Use:   "describe",
		Short: "Describe a Kafka cluster.",
		RunE:  c.describe,
	})
	c.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete a Kafka cluster.",
		RunE:  c.delete,
	})
	c.AddCommand(&cobra.Command{
		Use:   "update",
		Short: "Update a Kafka cluster.",
		RunE:  c.update,
	})
	c.AddCommand(&cobra.Command{
		Use:   "auth",
		Short: "Auth a Kafka cluster.",
		RunE:  c.auth,
	})
}

func (c *Command) create(Command *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}

func (c *Command) list(Command *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}

func (c *Command) describe(Command *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}

func (c *Command) update(Command *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}

func (c *Command) delete(Command *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}

func (c *Command) auth(Command *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}
