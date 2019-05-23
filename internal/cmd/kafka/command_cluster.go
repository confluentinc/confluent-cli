package kafka

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/confluentinc/ccloud-sdk-go"
	kafkav1 "github.com/confluentinc/ccloudapis/kafka/v1"
	"github.com/confluentinc/go-printer"

	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/config"
	"github.com/confluentinc/cli/internal/pkg/errors"
)

var (
	listFields      = []string{"Id", "Name", "ServiceProvider", "Region", "Durability", "Status"}
	listLabels      = []string{"Id", "Name", "Provider", "Region", "Durability", "Status"}
	describeFields  = []string{"Id", "Name", "NetworkIngress", "NetworkEgress", "Storage", "ServiceProvider", "Region", "Status", "Endpoint", "ApiEndpoint", "PricePerHour"}
	describeRenames = map[string]string{"NetworkIngress": "Ingress", "NetworkEgress": "Egress", "ServiceProvider": "Provider"}
)

type clusterCommand struct {
	*cobra.Command
	config *config.Config
	client ccloud.Kafka
	ch     *pcmd.ConfigHelper
}

// NewClusterCommand returns the Cobra command for Kafka cluster.
func NewClusterCommand(config *config.Config, client ccloud.Kafka, ch *pcmd.ConfigHelper) *cobra.Command {
	cmd := &clusterCommand{
		Command: &cobra.Command{
			Use:   "cluster",
			Short: "Manage Kafka clusters",
		},
		config: config,
		client: client,
		ch:     ch,
	}
	cmd.init()
	return cmd.Command
}

func (c *clusterCommand) init() {
	c.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List Kafka clusters",
		RunE:  c.list,
		Args:  cobra.NoArgs,
	})

	createCmd := &cobra.Command{
		Use:   "create NAME",
		Short: "Create a Kafka cluster",
		RunE:  c.create,
		Args:  cobra.ExactArgs(1),
	}
	createCmd.Flags().String("cloud", "", "Choose aws or gcp")
	createCmd.Flags().String("region", "", "A valid region in the given cloud")
	// default to smallest size allowed
	createCmd.Flags().Int32("ingress", 1, "Network ingress in MB/s")
	createCmd.Flags().Int32("egress", 1, "Network egress in MB/s")
	createCmd.Flags().Int32("storage", 500, "Total usable data storage in GB")
	createCmd.Flags().Bool("multizone", false, "Use multiple zones for high availability")
	check(createCmd.MarkFlagRequired("cloud"))
	check(createCmd.MarkFlagRequired("region"))
	createCmd.Flags().SortFlags = false
	createCmd.Hidden = true
	c.AddCommand(createCmd)

	c.AddCommand(&cobra.Command{
		Use:   "describe ID",
		Short: "Describe a Kafka cluster",
		RunE:  c.describe,
		Args:  cobra.ExactArgs(1),
	})

	updateCmd := &cobra.Command{
		Use:   "update ID",
		Short: "Update a Kafka cluster",
		RunE:  c.update,
		Args:  cobra.ExactArgs(1),
	}
	updateCmd.Hidden = true
	c.AddCommand(updateCmd)

	deleteCmd := &cobra.Command{
		Use:   "delete ID",
		Short: "Delete a Kafka cluster",
		RunE:  c.delete,
		Args:  cobra.ExactArgs(1),
	}
	deleteCmd.Hidden = true
	c.AddCommand(deleteCmd)
	c.AddCommand(&cobra.Command{
		Use:   "use ID",
		Short: "Make the Kafka cluster active for use in other commands",
		RunE:  c.use,
		Args:  cobra.ExactArgs(1),
	})
}

func (c *clusterCommand) list(cmd *cobra.Command, args []string) error {
	environment, err := pcmd.GetEnvironment(cmd, c.config)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	req := &kafkav1.KafkaCluster{AccountId: environment}
	clusters, err := c.client.List(context.Background(), req)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	currCtx, err := c.config.Context()
	if err != nil && err != errors.ErrNoContext {
		return err
	}
	var data [][]string
	for _, cluster := range clusters {
		if cluster.Id == currCtx.Kafka {
			cluster.Id = fmt.Sprintf("* %s", cluster.Id)
		} else {
			cluster.Id = fmt.Sprintf("  %s", cluster.Id)
		}
		data = append(data, printer.ToRow(cluster, listFields))
	}
	printer.RenderCollectionTable(data, listLabels)
	return nil
}

func (c *clusterCommand) create(cmd *cobra.Command, args []string) error {
	if true {
		return errors.ErrNotImplemented
	}

	cloud, err := cmd.Flags().GetString("cloud")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	region, err := cmd.Flags().GetString("region")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	ingress, err := cmd.Flags().GetInt32("ingress")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	egress, err := cmd.Flags().GetInt32("egress")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	storage, err := cmd.Flags().GetInt32("storage")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	multizone, err := cmd.Flags().GetBool("multizone")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	environment, err := pcmd.GetEnvironment(cmd, c.config)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	durability := kafkav1.Durability_LOW
	if multizone {
		durability = kafkav1.Durability_HIGH
	}
	cfg := &kafkav1.KafkaClusterConfig{
		AccountId:       environment,
		Name:            args[0],
		ServiceProvider: cloud,
		Region:          region,
		NetworkIngress:  ingress,
		NetworkEgress:   egress,
		Storage:         storage,
		Durability:      durability,
	}
	cluster, err := c.client.Create(context.Background(), cfg)
	if err != nil {
		// TODO: don't swallow validation errors (reportedly separately)
		return errors.HandleCommon(err, cmd)
	}
	return printer.RenderTableOut(cluster, describeFields, describeRenames, os.Stdout)
}

func (c *clusterCommand) describe(cmd *cobra.Command, args []string) error {
	environment, err := pcmd.GetEnvironment(cmd, c.config)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	req := &kafkav1.KafkaCluster{AccountId: environment, Id: args[0]}
	cluster, err := c.client.Describe(context.Background(), req)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	return printer.RenderTableOut(cluster, describeFields, describeRenames, os.Stdout)
}

func (c *clusterCommand) update(cmd *cobra.Command, args []string) error {
	return errors.ErrNotImplemented
}

func (c *clusterCommand) delete(cmd *cobra.Command, args []string) error {
	if true {
		return errors.ErrNotImplemented
	}

	environment, err := pcmd.GetEnvironment(cmd, c.config)

	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	req := &kafkav1.KafkaCluster{AccountId: environment, Id: args[0]}
	err = c.client.Delete(context.Background(), req)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	pcmd.Printf(cmd, "The Kafka cluster %s has been deleted.\n", args[0])
	return nil
}

func (c *clusterCommand) use(cmd *cobra.Command, args []string) error {
	clusterID := args[0]

	cfg, err := c.config.Context()
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	// This ensures that the clusterID actually exists or throws an error
	environment, err := pcmd.GetEnvironment(cmd, c.ch.Config)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	_, err = c.ch.KafkaClusterConfig(clusterID, environment)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	cfg.Kafka = clusterID
	return c.config.Save()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}