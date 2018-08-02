package kafka

import (
	"context"
	"os"

	"github.com/codyaray/go-printer"
	"github.com/spf13/cobra"

	schedv1 "github.com/confluentinc/cc-structs/kafka/scheduler/v1"
	"github.com/confluentinc/cli/command/common"
	"github.com/confluentinc/cli/shared"
)

var (
	listFields      = []string{"Id", "Name", "ServiceProvider", "Region", "Durability", "Status"}
	listLabels      = []string{"Id", "Name", "Provider", "Region", "Durability", "Status"}
	describeFields  = []string{"Id", "Name", "NetworkIngress", "NetworkEgress", "Storage", "ServiceProvider", "Region", "Status", "Endpoint", "PricePerHour"}
	describeRenames = map[string]string{"NetworkIngress": "Ingress", "NetworkEgress": "Egress", "ServiceProvider": "Provider"}
)

type clusterCommand struct {
	*cobra.Command
	config *shared.Config
	kafka  Kafka
}

// NewClusterCommand returns the Cobra clusterCommand for Kafka Cluster.
func NewClusterCommand(config *shared.Config, kafka Kafka) *cobra.Command {
	cmd := &clusterCommand{
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

func (c *clusterCommand) init() {
	c.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List Kafka clusters.",
		RunE:  c.list,
	})
	c.AddCommand(&cobra.Command{
		Use:   "create NAME",
		Short: "Create a Kafka cluster.",
		RunE:  c.create,
	})
	c.AddCommand(&cobra.Command{
		Use:   "describe ID",
		Short: "Describe a Kafka cluster.",
		RunE:  c.describe,
		Args:  cobra.ExactArgs(1),
	})
	c.AddCommand(&cobra.Command{
		Use:   "update ID",
		Short: "Update a Kafka cluster.",
		RunE:  c.update,
		Args:  cobra.ExactArgs(1),
	})
	c.AddCommand(&cobra.Command{
		Use:   "delete ID",
		Short: "Delete a Kafka cluster.",
		RunE:  c.delete,
		Args:  cobra.ExactArgs(1),
	})
	c.AddCommand(&cobra.Command{
		Use:   "auth",
		Short: "Auth a Kafka cluster.",
		RunE:  c.auth,
	})
	c.AddCommand(&cobra.Command{
		Use:   "use ID",
		Short: "Set a Kafka cluster in the CLI context.",
		RunE:  c.use,
		Args:  cobra.ExactArgs(1),
	})
}

func (c *clusterCommand) list(cmd *cobra.Command, args []string) error {
	req := &schedv1.KafkaCluster{AccountId: c.config.Auth.Account.Id}
	clusters, err := c.kafka.List(context.Background(), req)
	if err != nil {
		return common.HandleError(err)
	}
	var data [][]string
	for _, cluster := range clusters {
		data = append(data, printer.ToRow(cluster, listFields))
	}
	printer.RenderCollectionTable(data, listLabels)
	return nil
}

func (c *clusterCommand) create(cmd *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}

func (c *clusterCommand) describe(cmd *cobra.Command, args []string) error {
	req := &schedv1.KafkaCluster{AccountId: c.config.Auth.Account.Id, Id: args[0]}
	cluster, err := c.kafka.Describe(context.Background(), req)
	if err != nil {
		return common.HandleError(err)
	}
	printer.RenderTableOut(cluster, describeFields, describeRenames, os.Stdout)
	return nil
}

func (c *clusterCommand) update(cmd *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}

func (c *clusterCommand) delete(cmd *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}

func (c *clusterCommand) auth(cmd *cobra.Command, args []string) error {
	return shared.ErrNotImplemented
}

func (c *clusterCommand) use(cmd *cobra.Command, args []string) error {
	ctx, err := c.config.Context()
	if err != nil {
		return common.HandleError(err)
	}
	ctx.Kafka = args[0]
	return c.config.Save()
}
