package ksql

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/confluentinc/ccloud-sdk-go"
	kafkav1 "github.com/confluentinc/ccloudapis/kafka/v1"
	ksqlv1 "github.com/confluentinc/ccloudapis/ksql/v1"
	"github.com/confluentinc/go-printer"

	"github.com/confluentinc/cli/internal/pkg/acl"
	pcmd "github.com/confluentinc/cli/internal/pkg/cmd"
	"github.com/confluentinc/cli/internal/pkg/config"
	"github.com/confluentinc/cli/internal/pkg/errors"
)

var (
	listFields      = []string{"Id", "Name", "OutputTopicPrefix", "KafkaClusterId", "Storage", "Endpoint", "Status"}
	listLabels      = []string{"Id", "Name", "Topic Prefix", "Kafka", "Storage", "Endpoint", "Status"}
	describeFields  = []string{"Id", "Name", "OutputTopicPrefix", "KafkaClusterId", "Storage", "Endpoint", "Status"}
	describeRenames = map[string]string{"KafkaClusterId": "Kafka", "OutputTopicPrefix": "Topic Prefix"}
	aclsDryRun      = false
)

type clusterCommand struct {
	*cobra.Command
	config      *config.Config
	client      ccloud.KSQL
	kafkaClient ccloud.Kafka
	userClient  ccloud.User
}

// NewClusterCommand returns the Cobra clusterCommand for Ksql Cluster.
func NewClusterCommand(config *config.Config, client ccloud.KSQL, kafkaClient ccloud.Kafka, userClient ccloud.User) *cobra.Command {
	cmd := &clusterCommand{
		Command: &cobra.Command{
			Use:   "app",
			Short: "Manage KSQL apps",
		},
		config: config,
		client: client,
		kafkaClient: kafkaClient,
		userClient: userClient,
	}
	cmd.init()
	return cmd.Command
}

func (c *clusterCommand) init() {
	c.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List KSQL apps",
		RunE:  c.list,
		Args:  cobra.NoArgs,
	})

	createCmd := &cobra.Command{
		Use:   "create NAME",
		Short: "Create a KSQL app",
		RunE:  c.create,
		Args:  cobra.ExactArgs(1),
	}
	createCmd.Flags().String("cluster", "", "Kafka Cluster ID")
	check(createCmd.MarkFlagRequired("cluster"))
	createCmd.Flags().Int32("storage", 50, "total usable data storage in GB")
	check(createCmd.MarkFlagRequired("storage"))
	createCmd.Flags().Int32("servers", 1, "number of servers in the cluster")
	c.AddCommand(createCmd)

	c.AddCommand(&cobra.Command{
		Use:   "describe ID",
		Short: "Describe a ksql app",
		RunE:  c.describe,
		Args:  cobra.ExactArgs(1),
	})

	c.AddCommand(&cobra.Command{
		Use:   "delete ID",
		Short: "Delete a ksql app",
		RunE:  c.delete,
		Args:  cobra.ExactArgs(1),
	})

	aclsCmd := &cobra.Command{
		Use: "configure-acls ID TOPICS...",
		Short: "Configure acls for a KSQL cluster",
		RunE: c.configureACLs,
		Args: cobra.MinimumNArgs(1),
	}
	aclsCmd.Flags().String("cluster", "", "Kafka Cluster ID")
	check(createCmd.MarkFlagRequired("cluster"))
	aclsCmd.Flags().BoolVar(&aclsDryRun, "dry-run", false, "If specified, print the acls that will be set and exit")
	c.AddCommand(aclsCmd)
}

func (c *clusterCommand) list(cmd *cobra.Command, args []string) error {
	req := &ksqlv1.KSQLCluster{AccountId: c.config.Auth.Account.Id}
	clusters, err := c.client.List(context.Background(), req)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	var data [][]string
	for _, cluster := range clusters {
		data = append(data, printer.ToRow(cluster, listFields))
	}
	printer.RenderCollectionTable(data, listLabels)
	return nil
}

func (c *clusterCommand) create(cmd *cobra.Command, args []string) error {
	kafkaCluster, err := pcmd.GetKafkaCluster(cmd, c.config)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	storage, err := cmd.Flags().GetInt32("storage")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	servers, err := cmd.Flags().GetInt32("servers")
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	cfg := &ksqlv1.KSQLClusterConfig{
		AccountId:      c.config.Auth.Account.Id,
		Name:           args[0],
		Servers:        servers,
		Storage:        storage,
		KafkaClusterId: kafkaCluster.Id,
	}
	cluster, err := c.client.Create(context.Background(), cfg)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	return printer.RenderTableOut(cluster, describeFields, describeRenames, os.Stdout)
}

func (c *clusterCommand) describe(cmd *cobra.Command, args []string) error {
	req := &ksqlv1.KSQLCluster{AccountId: c.config.Auth.Account.Id, Id: args[0]}
	cluster, err := c.client.Describe(context.Background(), req)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	return printer.RenderTableOut(cluster, describeFields, describeRenames, os.Stdout)
}

func (c *clusterCommand) delete(cmd *cobra.Command, args []string) error {
	req := &ksqlv1.KSQLCluster{AccountId: c.config.Auth.Account.Id, Id: args[0]}
	err := c.client.Delete(context.Background(), req)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	pcmd.Printf(cmd, "The ksql app %s has been deleted.\n", args[0])
	return nil
}

func (c *clusterCommand) createAcl(prefix string, patternType kafkav1.PatternTypes_PatternType, operation kafkav1.ACLOperations_ACLOperation, resource kafkav1.ResourceTypes_ResourceType, serviceAccountId string) *kafkav1.ACLBinding {
	binding := &kafkav1.ACLBinding{
		Entry: &kafkav1.AccessControlEntryConfig{
			Host: "*",
		},
		Pattern: &kafkav1.ResourcePatternConfig{},
	}
	binding.Entry.PermissionType = kafkav1.ACLPermissionTypes_ALLOW
	binding.Entry.Operation = operation
	binding.Entry.Principal = "User:" + serviceAccountId
	binding.Pattern.PatternType = patternType
	binding.Pattern.ResourceType = resource
	binding.Pattern.Name = prefix
	return binding
}

func (c *clusterCommand) createClusterAcl(operation kafkav1.ACLOperations_ACLOperation, serviceAccountId string) *kafkav1.ACLBinding {
	binding := &kafkav1.ACLBinding{
		Entry: &kafkav1.AccessControlEntryConfig{
			Host: "*",
		},
		Pattern: &kafkav1.ResourcePatternConfig{},
	}
	binding.Entry.PermissionType = kafkav1.ACLPermissionTypes_ALLOW
	binding.Entry.Operation = operation
	binding.Entry.Principal = "User:" + serviceAccountId
	binding.Pattern.PatternType = kafkav1.PatternTypes_LITERAL
	binding.Pattern.ResourceType = kafkav1.ResourceTypes_CLUSTER
	binding.Pattern.Name = "kafka-cluster"
	return binding
}

func (c *clusterCommand) buildACLBindings(serviceAccountId string, cluster *ksqlv1.KSQLCluster, topics []string) []*kafkav1.ACLBinding {
	bindings := []*kafkav1.ACLBinding{}
	for _, op := range []kafkav1.ACLOperations_ACLOperation{
		kafkav1.ACLOperations_DESCRIBE,
		kafkav1.ACLOperations_DESCRIBE_CONFIGS,
	} {
		bindings = append(bindings, c.createClusterAcl(op, serviceAccountId))
	}
	for _ ,op := range []kafkav1.ACLOperations_ACLOperation{
		kafkav1.ACLOperations_CREATE,
		kafkav1.ACLOperations_DESCRIBE,
		kafkav1.ACLOperations_ALTER,
		kafkav1.ACLOperations_DESCRIBE_CONFIGS,
		kafkav1.ACLOperations_ALTER_CONFIGS,
		kafkav1.ACLOperations_READ,
		kafkav1.ACLOperations_WRITE,
		kafkav1.ACLOperations_DELETE,
	} {
		bindings = append(bindings, c.createAcl(cluster.OutputTopicPrefix, kafkav1.PatternTypes_PREFIXED, op, kafkav1.ResourceTypes_TOPIC, serviceAccountId))
		bindings = append(bindings, c.createAcl("_confluent-ksql-" + cluster.OutputTopicPrefix, kafkav1.PatternTypes_PREFIXED, op, kafkav1.ResourceTypes_TOPIC, serviceAccountId))
		bindings = append(bindings, c.createAcl("_confluent-ksql-" + cluster.OutputTopicPrefix, kafkav1.PatternTypes_PREFIXED, op, kafkav1.ResourceTypes_GROUP, serviceAccountId))
	}
	for _, op := range[]kafkav1.ACLOperations_ACLOperation{
		kafkav1.ACLOperations_DESCRIBE,
		kafkav1.ACLOperations_DESCRIBE_CONFIGS,
	} {
		bindings = append(bindings, c.createAcl("*", kafkav1.PatternTypes_LITERAL, op, kafkav1.ResourceTypes_TOPIC, serviceAccountId))
		bindings = append(bindings, c.createAcl("*", kafkav1.PatternTypes_LITERAL, op, kafkav1.ResourceTypes_GROUP, serviceAccountId))
	}
	for _ ,op := range []kafkav1.ACLOperations_ACLOperation{
		kafkav1.ACLOperations_DESCRIBE,
		kafkav1.ACLOperations_DESCRIBE_CONFIGS,
		kafkav1.ACLOperations_READ,
	} {
		for _, t := range topics {
			bindings = append(bindings, c.createAcl(t, kafkav1.PatternTypes_LITERAL, op, kafkav1.ResourceTypes_TOPIC, serviceAccountId))
		}
	}
	return bindings
}

func (c *clusterCommand) getServiceAccount(cluster *ksqlv1.KSQLCluster) (string, error) {
	users, err := c.userClient.GetServiceAccounts(context.Background())
	if err != nil {
		return "", err
	}
	for _, user := range users {
		if user.ServiceName == fmt.Sprintf("KSQL.%s", cluster.Id) {
			return strconv.Itoa(int(user.Id)), nil
		}
	}
	return "", errors.New(fmt.Sprintf("No service account found for %s", cluster.Id))
}

func (c *clusterCommand) configureACLs(cmd *cobra.Command, args[]string) error {
	ctx := context.Background()

	// Get the Kafka Cluster
	kafkaCluster, err := pcmd.GetKafkaCluster(cmd, c.config)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	// Ensure the Kafka Cluster is an Enterprise cluster
	kafkaCluster, err = c.kafkaClient.Describe(ctx, kafkaCluster)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	// Ensure the KSQL cluster talks to the current Kafka Cluster
	req := &ksqlv1.KSQLCluster{AccountId: c.config.Auth.Account.Id, Id: args[0]}
	cluster, err := c.client.Describe(context.Background(), req)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	if cluster.KafkaClusterId != kafkaCluster.Id {
		pcmd.Printf(cmd, "This KSQL cluster is not backed by the current Kafka cluster.")
	}

	// Ensure the Kafka Cluster is an Enterprise cluster
	if !kafkaCluster.Enterprise {
		pcmd.Printf(cmd, "Cluster is not an enterprise cluster. No ACLS need to be set")
		return nil
	}

	serviceAccountId, err := c.getServiceAccount(cluster)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}

	// Setup ACLs
	bindings := c.buildACLBindings(serviceAccountId, cluster, args[1:])
	if aclsDryRun {
		acl.PrintAcls(bindings, cmd.OutOrStderr())
		return nil
	}
	err = c.kafkaClient.CreateACL(ctx, kafkaCluster, bindings)
	if err != nil {
		return errors.HandleCommon(err, cmd)
	}
	return nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}