package kafka

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	kafkav1 "github.com/confluentinc/ccloudapis/kafka/v1"
)

// ACLConfiguration wrapper used for flag parsing and validation
type ACLConfiguration struct {
	*kafkav1.ACLBinding
	errors error
}

// aclConfigFlags returns a flag set which can be parsed to create an ACLConfiguration object.
func aclConfigFlags() *pflag.FlagSet {
	flgSet := aclEntryFlags()
	flgSet.SortFlags = false
	flgSet.AddFlagSet(resourceFlags())
	return flgSet
}

// aclEntryFlags returns a flag set which can be parsed to create an AccessControlEntry object.
func aclEntryFlags() *pflag.FlagSet {
	flgSet := pflag.NewFlagSet("acl-entry", pflag.ExitOnError)
	//flgSet.String("cluster", "", "Confluent Cloud cluster ID.")
	flgSet.Bool("allow", false, "Set the ACL to grant access.")
	flgSet.Bool("deny", false, "Set the ACL to restrict access to resource.")
	//flgSet.String( "host", "*", "Set Kafka principal host. Note: Not supported on CCLOUD.")
	flgSet.Int("service-account-id", 0, "The service account ID.")
	flgSet.String("operation", "", fmt.Sprintf("Set ACL Operation to: [%s].",
		listEnum(kafkav1.ACLOperations_ACLOperation_name, []string{"ANY", "UNKNOWN"})))
	// An error is only returned if the flag name is not present.
	// We know the flag name is present so its safe to ignore this.
	_ = cobra.MarkFlagRequired(flgSet, "service-account-id")
	_ = cobra.MarkFlagRequired(flgSet, "operation")
	return flgSet
}

// resourceFlags returns a flag set which can be parsed to create a ResourcePattern object.
func resourceFlags() *pflag.FlagSet {
	flgSet := pflag.NewFlagSet("acl-resource", pflag.ExitOnError)
	//flgSet.String("cluster", "", "The Confluent Cloud cluster ID.")
	flgSet.Bool("cluster-scope", false, `Set the cluster resource. With this option the ACL grants
access to the provided operations on the Kafka cluster itself.`)
	flgSet.String("topic", "", `Set the topic resource. With this option the ACL grants the provided
operations on the topics that start with that prefix, depending on whether
the --prefix option was also passed.`)
	flgSet.String("consumer-group", "", "Set the Consumer Group resource.")
	flgSet.String("transactional-id", "", "Set the TransactionalID resource.")
	flgSet.Bool("prefix", false, "Set to match all resource names prefixed with this value.")

	return flgSet
}

// parse returns ACLConfiguration from the contents of cmd
func parse(cmd *cobra.Command) *ACLConfiguration {
	aclBinding := &ACLConfiguration{
		ACLBinding: &kafkav1.ACLBinding{
			Entry: &kafkav1.AccessControlEntryConfig{
				Host: "*",
			},
			Pattern: &kafkav1.ResourcePatternConfig{},
		},
	}
	cmd.Flags().Visit(fromArgs(aclBinding))
	return aclBinding
}

// fromArgs maps command flag values to the appropriate ACLConfiguration field
func fromArgs(conf *ACLConfiguration) func(*pflag.Flag) {
	return func(flag *pflag.Flag) {
		v := flag.Value.String()
		switch n := flag.Name; n {
		case "consumer-group":
			setResourcePattern(conf, "GROUP", v)
		case "cluster-scope":
			// The only valid name for a cluster is kafka-cluster
			// https://github.com/confluentinc/cc-kafka/blob/88823c6016ea2e306340938994d9e122abf3c6c0/core/src/main/scala/kafka/security/auth/Resource.scala#L24
			setResourcePattern(conf, "cluster", "kafka-cluster")
		case "topic":
			fallthrough
		case "delegation-token":
			fallthrough
		case "transactional-id":
			setResourcePattern(conf, n, v)
		case "allow":
			conf.Entry.PermissionType = kafkav1.ACLPermissionTypes_ALLOW
		case "deny":
			conf.Entry.PermissionType = kafkav1.ACLPermissionTypes_DENY
		case "prefix":
			conf.Pattern.PatternType = kafkav1.PatternTypes_PREFIXED
		case "service-account-id":
			if v == "0" {
				conf.Entry.Principal = "User:*"
				break
			}
			conf.Entry.Principal = "User:" + v
		case "operation":
			v = strings.ToUpper(v)
			v = strings.Replace(v, "-", "_", -1)
			if op, ok := kafkav1.ACLOperations_ACLOperation_value[v]; ok {
				conf.Entry.Operation = kafkav1.ACLOperations_ACLOperation(op)
				break
			}
			conf.errors = multierror.Append(conf.errors, fmt.Errorf("Invalid operation value: "+v))
		}
	}
}

func setResourcePattern(conf *ACLConfiguration, n, v string) {
	/* Normalize the resource pattern name */
	if conf.Pattern.ResourceType != kafkav1.ResourceTypes_UNKNOWN {
		conf.errors = multierror.Append(conf.errors, fmt.Errorf("exactly one of %v must be set",
			listEnum(kafkav1.ResourceTypes_ResourceType_name, []string{"ANY", "UNKNOWN"})))
		return
	}

	n = strings.ToUpper(n)
	n = strings.Replace(n, "-", "_", -1)

	conf.Pattern.ResourceType = kafkav1.ResourceTypes_ResourceType(kafkav1.ResourceTypes_ResourceType_value[n])

	if conf.Pattern.ResourceType == kafkav1.ResourceTypes_CLUSTER {
		conf.Pattern.PatternType = kafkav1.PatternTypes_LITERAL
	}
	conf.Pattern.Name = v
}

func listEnum(enum map[int32]string, exclude []string) string {
	var ops []string

OUTER:
	for _, v := range enum {
		for _, exclusion := range exclude {
			if v == exclusion {
				continue OUTER
			}
		}
		if v == "GROUP" {
			v = "consumer-group"
		}
		if v == "CLUSTER" {
			v = "cluster-scope"
		}
		v = strings.Replace(v, "_", "-", -1)
		ops = append(ops, strings.ToLower(v))
	}

	sort.Strings(ops)
	return strings.Join(ops, ", ")
}