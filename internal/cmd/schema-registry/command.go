package schema_registry

import (
	"context"
	ccsdk "github.com/confluentinc/ccloud-sdk-go"
	srv1 "github.com/confluentinc/ccloudapis/schemaregistry/v1"
	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/config"
	"github.com/confluentinc/cli/internal/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

type command struct {
	*cobra.Command
	config   *config.Config
	ccClient ccsdk.SchemaRegistry
}

func New(prerunner pcmd.PreRunner, config *config.Config, ccloudClient ccsdk.SchemaRegistry) *cobra.Command {
	cmd := &command{
		Command: &cobra.Command{
			Use:               "schema-registry",
			Short:             `Manage Schema Registry.`,
			PersistentPreRunE: prerunner.Authenticated(),
		},
		config:   config,
		ccClient: ccloudClient,
	}
	cmd.init()
	return cmd.Command
}

func (c *command) init() {
	createCmd := &cobra.Command{
		Use:     "enable",
		Short:   `Enable Schema Registry for this account.`,
		Example: `ccloud schema-registry enable --cloud gcp`,
		RunE:    c.enable,
		Args:    cobra.NoArgs,
	}
	createCmd.Flags().String("cluster", "", "Kafka cluster ID.")
	createCmd.Flags().String("cloud", "", "Cloud provider ('aws', 'azure', or 'gcp').")
	_ = createCmd.MarkFlagRequired("cloud")
	createCmd.Flags().String("geo", "", "Either 'us', 'eu', or 'apac' (only applies to Enterprise accounts).")
	createCmd.Flags().SortFlags = false
	c.AddCommand(createCmd)

}

func (c *command) enable(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Collect the parameters
	accountId, err := pcmd.GetEnvironment(cmd, c.config)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	serviceProvider, err := cmd.Flags().GetString("cloud")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	locationFlag, err := cmd.Flags().GetString("geo")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	// Trust the API will handle CCP/CCE and whether geo is required
	location := srv1.GlobalSchemaRegistryLocation(srv1.GlobalSchemaRegistryLocation_value[strings.ToUpper(locationFlag)])

	kafkaClusterId, err := cmd.Flags().GetString("cluster")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	// Build the SR instance
	clusterConfig := &srv1.SchemaRegistryClusterConfig{
		KafkaClusterId:  kafkaClusterId,
		AccountId:       accountId,
		Location:        location,
		ServiceProvider: serviceProvider,
		// Name is a special string that everyone expects. Originally, this field was added to support
		// multiple SR instances, but for now there's a contract between our services that it will be
		// this hardcoded string constant
		Name: "account schema-registry",
	}
	pcmd.Println(cmd, kafkaClusterId+accountId+serviceProvider)
	newCluster, err := c.ccClient.CreateSchemaRegistryCluster(ctx, clusterConfig)
	if err != nil {
		// If it already exists, return the existing one
		existingClusters, getExistingErr := c.ccClient.GetSchemaRegistryClusters(ctx, &srv1.SchemaRegistryCluster{
			AccountId: accountId,
		})
		if getExistingErr != nil {
			return errors.HandleCommon(getExistingErr, cmd)
		}
		if len(existingClusters) > 0 {
			pcmd.Println(cmd, "Cluster already exists:")
			for _, cluster := range existingClusters {
				pcmd.Println(cmd, "Cluster ID: "+cluster.Id)
				pcmd.Println(cmd, "Endpoint: "+cluster.Endpoint)
			}
			return nil
		}
		return errors.HandleCommon(err, cmd)
	}

	pcmd.Println(cmd, "Cluster already exists:")
	pcmd.Println(cmd, "Cluster ID: "+newCluster.Id)
	pcmd.Println(cmd, "Endpoint: "+newCluster.Endpoint)

	return nil
}