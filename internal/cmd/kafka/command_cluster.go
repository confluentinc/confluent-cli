package kafka

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/confluentinc/ccloud-sdk-go"
	authv1 "github.com/confluentinc/ccloudapis/auth/v1"
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
}

// NewClusterCommand returns the Cobra command for Kafka cluster.
func NewClusterCommand(config *config.Config, client ccloud.Kafka) *cobra.Command {
	cmd := &clusterCommand{
		Command: &cobra.Command{
			Use:   "cluster",
			Short: "Manage Kafka clusters",
		},
		config: config,
		client: client,
	}
	cmd.init()
	return cmd.Command
}

func (c *clusterCommand) init() {
	// Promote this to command.go when/if topic/ACL commands also want this flag
	c.PersistentFlags().String("environment", "", "ID of the environment in which to run the command")

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
		Use:   "auth",
		Short: "Configure authorization for a Kafka cluster",
		RunE:  c.auth,
		Args:  cobra.NoArgs,
	})
	c.AddCommand(&cobra.Command{
		Use:   "use ID",
		Short: "Make the Kafka cluster active for use in other commands",
		RunE:  c.use,
		Args:  cobra.ExactArgs(1),
	})
}

func (c *clusterCommand) getEnvironment(cmd *cobra.Command) (string, error) {
	environment, err := cmd.Flags().GetString("environment")
	if err != nil {
		return "", err
	}
	if environment == "" {
		environment = c.config.Auth.Account.Id
	}
	return environment, nil
}

func (c *clusterCommand) list(cmd *cobra.Command, args []string) error {
	environment, err := c.getEnvironment(cmd)
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
	environment, err := c.getEnvironment(cmd)
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
	environment, err := c.getEnvironment(cmd)
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

	environment, err := c.getEnvironment(cmd)
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

func (c *clusterCommand) auth(cmd *cobra.Command, args []string) error {
	cfg, err := c.config.Context()

	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	if cfg.Kafka == "" {
		return fmt.Errorf("No cluster selected. See ccloud kafka cluster use for help ")
	}

	environment, err := c.getEnvironment(cmd)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	cluster, known := c.config.Platforms[cfg.Platform].KafkaClusters[cfg.Kafka]
	if known {
		pcmd.Printf(cmd, "Kafka Cluster: %s\n", cfg.Kafka)
		pcmd.Printf(cmd, "Bootstrap Servers: %s\n", cluster.Bootstrap)
		pcmd.Printf(cmd, "API Key: %s\n", cluster.APIKey)
		pcmd.Printf(cmd, "API Secret: %s\n", cluster.APISecret)
		return nil
	}

	userProvidingKey, err := userHasKey(cmd, cfg.Kafka)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	var key, secret string
	if userProvidingKey {
		key, secret, err = promptForKafkaCreds(cmd)
	} else {
		key, secret, err = c.createKafkaCreds(context.Background(), cmd, environment, cfg.Kafka)
	}
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	req := &kafkav1.KafkaCluster{AccountId: environment, Id: cfg.Kafka}
	kc, err := c.client.Describe(context.Background(), req)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	if c.config.Platforms[cfg.Platform].KafkaClusters == nil {
		c.config.Platforms[cfg.Platform].KafkaClusters = map[string]config.KafkaClusterConfig{}
	}
	c.config.Platforms[cfg.Platform].KafkaClusters[cfg.Kafka] = config.KafkaClusterConfig{
		Bootstrap:   strings.TrimPrefix(kc.Endpoint, "SASL_SSL://"),
		APIEndpoint: kc.ApiEndpoint,
		APIKey:      key,
		APISecret:   secret,
	}
	return c.config.Save()
}

func (c *clusterCommand) use(cmd *cobra.Command, args []string) error {
	clusterID := args[0]

	environment, err := c.getEnvironment(cmd)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	cfg, err := c.config.Context()
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	req := &kafkav1.KafkaCluster{AccountId: environment, Id: clusterID}
	kc, err := c.client.Describe(context.Background(), req)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	// TODO: we can't rely on "use" to do this if we want to support stateless usage with --cluster instead
	if c.config.Platforms[cfg.Platform].KafkaClusters == nil {
		c.config.Platforms[cfg.Platform].KafkaClusters = map[string]config.KafkaClusterConfig{}
	}
	c.config.Platforms[cfg.Platform].KafkaClusters[kc.Id] = config.KafkaClusterConfig{
		ID:          kc.Id,
		Bootstrap:   strings.TrimPrefix(kc.Endpoint, "SASL_SSL://"),
		APIEndpoint: kc.ApiEndpoint,
	}

	cfg.Kafka = clusterID
	return c.config.Save()
}

//
// Helper functions
//

func userHasKey(cmd *cobra.Command, kafkaClusterID string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	pcmd.Printf(cmd, "Do you have an API key for %s? [N/y] ", kafkaClusterID)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	r := strings.TrimSpace(response)
	return r == "" || r[0] == 'y' || r[0] == 'Y', nil
}

func promptForKafkaCreds(cmd *cobra.Command) (string, string, error) {
	reader := bufio.NewReader(os.Stdin)
	pcmd.Print(cmd, "API Key: ")
	key, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	pcmd.Print(cmd, "API Secret: ")
	// TODO: this should be using our internal prompt.ReadPassword() for testability
	byteSecret, err := terminal.ReadPassword(0)
	pcmd.Println(cmd)
	if err != nil {
		return "", "", err
	}
	secret := string(byteSecret)

	return strings.TrimSpace(key), strings.TrimSpace(secret), nil
}

func (c *clusterCommand) createKafkaCreds(ctx context.Context, cmd *cobra.Command, environment string, kafkaClusterID string) (string, string, error) {
	client := ccloud.NewClientWithJWT(ctx, c.config.AuthToken, c.config.AuthURL, c.config.Logger)
	key, err := client.APIKey.Create(ctx, &authv1.ApiKey{
		UserId: c.config.Auth.User.Id,
		LogicalClusters: []*authv1.ApiKey_Cluster{
			{Id: kafkaClusterID},
		},
		AccountId: environment,
	})

	if err != nil {
		return "", "", errors.ConvertAPIError(err)
	}
	pcmd.Println(cmd, "Okay, we've created an API key. If needed, you can see it with `ccloud kafka cluster auth`.")
	return key.Key, key.Secret, nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}