// This is a set of "linters" defining the CLI specification
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/client9/gospell"
	"github.com/hashicorp/go-multierror"

	"github.com/confluentinc/cli/internal/cmd"
	"github.com/confluentinc/cli/internal/pkg/config"
	linter "github.com/confluentinc/cli/internal/pkg/lint-cli"
	"github.com/confluentinc/cli/internal/pkg/log"
	"github.com/confluentinc/cli/internal/pkg/version"
)

var (
	debug   = flag.Bool("debug", false, "print debug output")
	affFile = flag.String("aff-file", "", "hunspell .aff file")
	dicFile = flag.String("dic-file", "", "hunspell .dic file")

	vocab *gospell.GoSpell

	properNouns = []string{
		"Apache", "Kafka", "CLI", "API", "ACL", "ACLs", "Confluent Cloud", "Confluent Platform",
	}
	vocabWords = []string{
		"ccloud", "kafka", "api", "acl", "url", "config", "multizone", "transactional", "decrypt",
	}
	utilityCommands = []string{
		"login", "logout", "version", "completion <shell>", "update",
	}
	nonClusterScopedCommands = []linter.RuleFilter{
		linter.OnlyLeafCommands, linter.ExcludeCommand(utilityCommands...),
		linter.ExcludeUse("local"), linter.ExcludeParentUse("environment", "service-account"),
		// these all require explicit cluster as id/name args
		linter.ExcludeCommandContains("kafka cluster"),
		// this doesn't need a --cluster override since you provide the api key itself to identify it
		linter.ExcludeCommandContains("api-key update", "api-key delete"),
		// this doesn't need a --cluster
		linter.ExcludeCommandContains("secret"),
	}
)

var rules = []linter.Rule{
	linter.Filter(
		linter.RequireNamedArgument(
			linter.NamedArgumentConfig{CreateCommandArg: "<name>", OtherCommandsArg: "<id>"},
			map[string]linter.NamedArgumentConfig{
				"environment": {CreateCommandArg: "<name>", OtherCommandsArg: "<environment-id>"},
				"topic":       {CreateCommandArg: "<topic>", OtherCommandsArg: "<topic>"},
				"api-key":     {CreateCommandArg: "N/A", OtherCommandsArg: "<apikey>"},
			},
		),
		linter.OnlyLeafCommands, linter.ExcludeCommand(utilityCommands...),
		// skip resource container commands
		linter.ExcludeUse("list", "auth"),
		// skip ACLs which don't have an identity (value objects rather than entities)
		linter.ExcludeCommandContains("kafka acl"),
		// skip api-key create since you don't get to choose a name for API keys
		linter.ExcludeCommandContains("api-key create"),
		// skip local which delegates to bash commands
		linter.ExcludeCommandContains("local"),
		// skip for api-key store command since KEY is not last argument
		linter.ExcludeCommand("api-key store <apikey> <secret>"),
		// skip secret commands
		linter.ExcludeCommandContains("secret"),
	),
	// TODO: ensuring --cluster is optional DOES NOT actually ensure that the cluster context is used
	linter.Filter(linter.RequireFlag("cluster", true), nonClusterScopedCommands...),
	linter.Filter(linter.RequireFlagType("cluster", "string"), nonClusterScopedCommands...),
	linter.Filter(linter.RequireFlagDescription("cluster", "Kafka cluster ID."),
		append(nonClusterScopedCommands, linter.ExcludeParentUse("api-key"))...),
	linter.RequireFlagSort(false),
	linter.RequireLowerCase("Use"),
	linter.RequireSingular("Use"),
	linter.Filter(
		linter.RequireSuffix("Short", "This is only available for Confluent Cloud Enterprise users."),
		// only include ACLs as they have a really long suffix/disclaimer that they're CCE only
		linter.IncludeCommandContains("kafka acl"),
		// only include service-accounts as they have a really long suffix/disclaimer that they're CCE only
		linter.IncludeCommandContains("service-account"),
	),
	linter.Filter(
		linter.RequireLengthBetween("Short", 13, 60),
		linter.ExcludeCommandContains("secret"),
		// skip ACLs as they have a really long suffix/disclaimer that they're CCE only
		linter.ExcludeCommandContains("kafka acl"),
		// skip service-accounts as they have a really long suffix/disclaimer that they're CCE only
		linter.ExcludeCommandContains("service-account"),
	),
	linter.RequireStartWithCapital("Short"),
	linter.RequireEndWithPunctuation("Short", false),
	linter.RequireCapitalizeProperNouns("Short", properNouns),
	linter.RequireStartWithCapital("Long"),
	linter.RequireEndWithPunctuation("Long", true),
	linter.RequireCapitalizeProperNouns("Long", properNouns),
	linter.Filter(linter.RequireNotTitleCase("Short", properNouns),
		linter.ExcludeCommandContains("secret")),
	linter.RequireRealWords("Use", '-'),
}

var flagRules = []linter.FlagRule{
	linter.FlagFilter(linter.RequireFlagNameLength(2, 16),
		linter.ExcludeFlag("service-account-id", "replication-factor", "local-secrets-file", "remote-secrets-file")),
	linter.RequireFlagStartWithCapital,
	linter.RequireFlagEndWithPunctuation,
	linter.RequireFlagCharacters('-'),
	linter.FlagFilter(linter.RequireFlagDelimiter('-', 1),
		linter.ExcludeFlag("service-account-id", "local-secrets-file", "remote-secrets-file")),
	linter.RequireFlagRealWords('-'),
}

func main() {
	flag.Parse()

	var err error
	vocab, err = gospell.NewGoSpell(*affFile, *dicFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, w := range vocabWords {
		vocab.AddWordRaw(w)
	}
	linter.SetVocab(vocab)

	l := linter.Linter{
		Rules:     rules,
		FlagRules: flagRules,
		Vocab:     vocab,
		Debug:     *debug,
	}

	var issues *multierror.Error
	for _, cliName := range []string{"confluent", "ccloud"} {
		cli, err := cmd.NewConfluentCommand(cliName, &config.Config{CLIName: cliName}, &version.Version{Binary: cliName}, log.New())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = l.Lint(cli)
		if err != nil {
			issues = multierror.Append(issues, err)
		}
	}
	if issues.ErrorOrNil() != nil {
		fmt.Println(issues)
		os.Exit(1)
	}
}