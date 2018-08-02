package kafka

import (
	"github.com/spf13/cobra"

	"github.com/confluentinc/cli/shared"
)

type topicCommand struct {
	*cobra.Command
	config *shared.Config
	kafka  Kafka
}

// NewTopicCommand returns the Cobra clusterCommand for Kafka Cluster.
func NewTopicCommand(config *shared.Config, kafka Kafka) *cobra.Command {
	cmd := &topicCommand{
		Command: &cobra.Command{
			Use:   "cluster",
			Short: "Manage kafka clusters.",
		},
		config: config,
		kafka:  kafka,
	}
	cmd.init()
	return cmd.Command
}

func (c *topicCommand) init() error {
	c.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List Kafka topics.",
		RunE:  c.list,
		Args: cobra.NoArgs,
	})
	c.AddCommand(&cobra.Command{
		Use:   "create TOPIC",
		Short: "Create a Kafka topic.",
		RunE:  c.create,
	})
	c.AddCommand(&cobra.Command{
		Use:   "describe TOPIC",
		Short: "Describe a Kafka topic.",
		RunE:  c.describe,
		Args:  cobra.ExactArgs(1),
	})
	c.AddCommand(&cobra.Command{
		Use:   "update TOPIC",
		Short: "Update a Kafka topic.",
		RunE:  c.update,
		Args:  cobra.ExactArgs(1),
	})
	c.AddCommand(&cobra.Command{
		Use:   "delete TOPIC",
		Short: "Delete a Kafka topic.",
		RunE:  c.delete,
		Args:  cobra.ExactArgs(1),
	})
	c.AddCommand(&cobra.Command{
		Use:   "produce TOPIC",
		Short: "Produce messages to a Kafka topic.",
		RunE:  c.produce,
		Args:  cobra.ExactArgs(1),
	})
	c.AddCommand(&cobra.Command{
		Use:   "consume TOPIC",
		Short: "Consume messages from a Kafka topic.",
		RunE:  c.consume,
		Args:  cobra.ExactArgs(1),
	})
	return nil
}

func (c *topicCommand) list(cmd *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}

func (c *topicCommand) create(cmd *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}

func (c *topicCommand) describe(cmd *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}

func (c *topicCommand) update(cmd *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}

func (c *topicCommand) delete(cmd *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}

func (c *topicCommand) produce(cmd *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}

func (c *topicCommand) consume(cmd *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}
