package lint_cli

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/client9/gospell"
	"github.com/gobuffalo/flect"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	alnum, _ = regexp.Compile("[^a-zA-Z0-9]+")
)

type Rule func(cmd *cobra.Command) error
type FlagRule func(flag *pflag.Flag, cmd *cobra.Command) error

var vocab *gospell.GoSpell

// TODO/HACK: this is to inject a vocab "global" object for use by the rules
func SetVocab(v *gospell.GoSpell) {
	vocab = v
}

// RequireRealWords checks that a field uses delimited-real-words, not smushcasecommands
func RequireRealWords(field string, delimiter rune) Rule {
	return func(cmd *cobra.Command) error {
		fieldValue := getValueByName(cmd, field)
		var issues *multierror.Error
		bareCmd := strings.Split(fieldValue, " ")[0] // TODO should we check all parts?
		for _, w := range strings.Split(bareCmd, string(delimiter)) {
			if ok := vocab.Spell(w); !ok {
				issue := fmt.Errorf("%s should consist of delimited real english words for %s on %s",
					normalizeDesc(field), bareCmd, FullCommand(cmd))
				issues = multierror.Append(issues, issue)
			}
		}
		return issues
	}
}

// RequireEndWithPunctuation checks that a field ends with a period
func RequireEndWithPunctuation(field string, ignoreIfEndsWithCodeBlock bool) Rule {
	return func(cmd *cobra.Command) error {
		fieldValue := getValueByName(cmd, field)
		chomped := strings.TrimRight(fieldValue, "\n")
		lines := strings.Split(fieldValue, "\n")
		if cmd.Long != "" && chomped[len(chomped)-1] != '.' {
			lastLine := len(lines) - 1
			if lines[len(lines)-1] == "" {
				lastLine = len(lines) - 2
			}
			// ignore rule if last line is code block
			if !strings.HasPrefix(lines[lastLine], "  ") || !ignoreIfEndsWithCodeBlock {
				return fmt.Errorf("%s should end with punctuation on %s", normalizeDesc(field), FullCommand(cmd))
			}
		}
		return nil
	}
}

// RequireNotEndWithPunctuation checks that a field does not end with a period
func RequireNotEndWithPunctuation(field string) Rule {
	return func(cmd *cobra.Command) error {
		fieldValue := getValueByName(cmd, field)
		if fieldValue[len(fieldValue)-1] == '.' {
			return fmt.Errorf("%s should not end with punctuation on %s", normalizeDesc(field), FullCommand(cmd))
		}
		return nil
	}
}

// RequireStartWithCapital checks that a field starts with a capital letter
func RequireStartWithCapital(field string) Rule {
	return func(cmd *cobra.Command) error {
		fieldValue := getValueByName(cmd, field)
		if fieldValue != "" && (fieldValue[0] < 'A' || fieldValue[0] > 'Z') {
			return fmt.Errorf("%s should start with a capital on %s", normalizeDesc(field), FullCommand(cmd))
		}
		return nil
	}
}

// RequireCapitalizeProperNouns checks that a field capitalizes proper nouns
func RequireCapitalizeProperNouns(field string, properNouns []string) Rule {
	index := map[string]string{}
	for _, n := range properNouns {
		index[strings.ToLower(n)] = n
	}
	return func(cmd *cobra.Command) error {
		fieldValue := getValueByName(cmd, field)
		var issues *multierror.Error
		for _, word := range strings.Split(fieldValue, " ") {
			if v, found := index[strings.ToLower(word)]; found && word != v {
				issue := fmt.Errorf("%s should capitalize %s on %s", normalizeDesc(field), v, FullCommand(cmd))
				issues = multierror.Append(issues, issue)
			}
		}
		return issues.ErrorOrNil()
	}
}

// RequireLengthBetween checks that a field is between a certain min and max length
func RequireLengthBetween(field string, minLength, maxLength int) Rule {
	return func(cmd *cobra.Command) error {
		fieldValue := getValueByName(cmd, field)
		var issues *multierror.Error
		if len(fieldValue) < minLength {
			issue := fmt.Errorf("%s is too short on %s - %s", normalizeDesc(field), FullCommand(cmd), cmd.Short)
			issues = multierror.Append(issues, issue)
		}
		if len(fieldValue) > maxLength {
			issue := fmt.Errorf("%s is too long on %s", normalizeDesc(field), FullCommand(cmd))
			issues = multierror.Append(issues, issue)
		}
		return issues
	}
}

// RequireSingular checks that a field is singular (not plural)
func RequireSingular(field string) Rule {
	return func(cmd *cobra.Command) error {
		fieldValue := getValueByName(cmd, field)
		if flect.Singularize(fieldValue) != fieldValue {
			return fmt.Errorf("%s should be singular for %s", normalizeDesc(field), FullCommand(cmd))
		}
		return nil
	}
}

// RequireLowerCase checks that a field is lower case
func RequireLowerCase(field string) Rule {
	return func(cmd *cobra.Command) error {
		fieldValue := getValueByName(cmd, field)
		command := strings.Split(fieldValue, " ")[0]
		if strings.ToLower(command) != command {
			return fmt.Errorf("%s should be lower case for %s", normalizeDesc(field), FullCommand(cmd))
		}
		return nil
	}
}

// NamedArgumentConfig lets you specify different argument names in the help/usage
// for create commands vs other commands; e.g., to pass NAME on create and ID elsewhere.
type NamedArgumentConfig struct {
	CreateCommandArg string
	OtherCommandsArg string
}

// RequireNamedArgument checks that a command has a single argument with the appropriate name.
// You can specify different names for create commands vs other commands; e.g., to pass NAME on create and ID elsewhere.
// You can also pass a string of overrides for some commands, identified by their parent, to have a different config.
func RequireNamedArgument(defConfig NamedArgumentConfig, overrides map[string]NamedArgumentConfig) Rule {
	return func(cmd *cobra.Command) error {
		// check whether arg parsing is setup correctly to expect exactly 1 arg (the ID/Name)
		if reflect.ValueOf(cmd.Args).Pointer() != reflect.ValueOf(cobra.ExactArgs(1)).Pointer() {
			return fmt.Errorf("missing expected argument on %s", FullCommand(cmd))
		}

		// check whether the usage string is setup correctly
		if o, found := overrides[cmd.Parent().Use]; found {
			if strings.HasPrefix(cmd.Use, "create ") {
				if !strings.HasSuffix(cmd.Use, o.CreateCommandArg) {
					return fmt.Errorf("bad usage string: must have %s in %s",
						o.CreateCommandArg, FullCommand(cmd))
				}
			} else if !strings.HasSuffix(cmd.Use, o.OtherCommandsArg) {
				return fmt.Errorf("bad usage string: must have %s in %s",
					o.OtherCommandsArg, FullCommand(cmd))
			}
		} else {
			// check for "create NAME" and "<verb> ID" elsewhere
			if strings.HasPrefix(cmd.Use, "create ") {
				if !strings.HasSuffix(cmd.Use, defConfig.CreateCommandArg) {
					return fmt.Errorf("bad usage string: must have %s in %s",
						defConfig.CreateCommandArg, FullCommand(cmd))
				}
			} else if !strings.HasSuffix(cmd.Use, defConfig.OtherCommandsArg) {
				return fmt.Errorf("bad usage string: must have %s in %s",
					defConfig.OtherCommandsArg, FullCommand(cmd))
			}
		}

		return nil
	}
}

// RequireNotTitleCase checks that a field is Not Title Casing Everything.
// You may pass a list of proper nouns that should always be capitalized, however.
func RequireNotTitleCase(field string, properNouns []string) Rule {
	// TODO: this is an awful IsTitleCase heuristic
	// index properNouns by the first word for easy search; value is the list of all individual words in properNoun
	// since multiple phrases can start with same word, we have an []string for each phrase (hence [][]string)
	index := map[string][][]string{}
	for _, noun := range properNouns {
		parts := strings.Split(noun, " ")
		// we can index multiple phrases starting with the same noun
		if _, ok := index[parts[0]]; ok {
			index[parts[0]] = append(index[parts[0]], parts)
		} else {
			index[parts[0]] = [][]string{parts}
		}
	}
	return func(cmd *cobra.Command) error {
		fieldValue := getValueByName(cmd, field)
		var issues *multierror.Error
		words := strings.Split(fieldValue, " ")
		for i := 0; i < len(words); i++ {
			word := alnum.ReplaceAllString(words[i], "") // Remove any punctuation before comparison
			if word[0] >= 'A' && word[0] <= 'Z' {
				// We have to start our check/loop at i=0 in case the command starts with a multi-word proper noun
				// But we don't consider capitalizing the first word of the sentence (i=0) to be title case
				var isTitleCase bool
				if i == 0 {
					isTitleCase = false
				} else {
					isTitleCase = true
				}
				if _, ok := index[word]; ok {
					// Check if index is a single-word proper noun; if so, allow it.
					if len(index[word]) == 1 {
						isTitleCase = false
					} else {
						// Check if index is a multi-word proper noun; if all words in the phrase match, allow it.
						parts := strings.Split(word, " ")
						if phrases, ok := index[parts[0]]; ok {
							// For each set of proper nouns starting with the current word...
							for _, wp := range phrases {
								// ...check if the next words in the command match the proper noun phrase
								allMatch := true
								for j := 0; j < len(wp); j++ {
									if words[i+j] != wp[j] {
										allMatch = false
										break
									}
								}
								// Hurray! All words in the word phrase matched the next words in the command
								if allMatch {
									isTitleCase = false
									i += len(wp) // skip any remaining words in the phrase; we've already OK'd them
									break        // we don't need to check anymore phrases
								}
							}
						}
					}
				}
				if isTitleCase {
					issue := fmt.Errorf("don't title case %s on %s - %s",
						normalizeDesc(field), FullCommand(cmd), cmd.Short)
					issues = multierror.Append(issues, issue)
				}
			}
		}
		return issues
	}
}

// RequireFlag checks that a flag is defined and whether it should be optional or required
func RequireFlag(flag string, optional bool) Rule {
	return func(cmd *cobra.Command) error {
		f := cmd.Flag(flag)
		if f == nil {
			return fmt.Errorf("missing --%s flag on %s", flag, FullCommand(cmd))
		} else {
			if optional && f.Annotations[cobra.BashCompOneRequiredFlag] != nil &&
				f.Annotations[cobra.BashCompOneRequiredFlag][0] == "true" {
				return fmt.Errorf("required --%s flag should be optional on %s", flag, FullCommand(cmd))
			}
		}
		return nil
	}
}

// RequireFlagType checks that a flag has the specified type, if it exists.
// Please use RequireFlag to check that it exists first.
func RequireFlagType(flag, typeName string) Rule {
	return func(cmd *cobra.Command) error {
		f := cmd.Flag(flag)
		if f != nil {
			// check that --flag has the right type (so its not a different meaning)
			if typeName != "" && f.Value.Type() != typeName {
				return fmt.Errorf("standard --%s flag has the wrong type on %s", flag, FullCommand(cmd))
			}
		}
		return nil
	}
}

// RequireFlagDescription checks that a flag has the specified usage string, if it exists.
// Please use RequireFlag to check that it exists first.
func RequireFlagDescription(flag, description string) Rule {
	return func(cmd *cobra.Command) error {
		f := cmd.Flag(flag)
		if f != nil {
			// check that --flag has the standard description (so its not a different meaning)
			if description != "" && f.Usage != description {
				return fmt.Errorf("bad usage string: expected standard description for --%s on %s",
					flag, FullCommand(cmd))
			}
		}
		return nil
	}
}

// RequireFlagSort checks whether flags should be auto sorted
func RequireFlagSort(sort bool) Rule {
	return func(cmd *cobra.Command) error {
		if cmd.Flags().HasFlags() && cmd.Flags().SortFlags != sort {
			if sort {
				return fmt.Errorf("flags not sorted on %s", FullCommand(cmd))
			}
			return fmt.Errorf("flags unexpectedly sorted on %s", FullCommand(cmd))
		}
		return nil
	}
}

// RequireFlagRealWords checks that a flag uses delimited-real-words, not --smushcaseflags
func RequireFlagRealWords(delim rune) FlagRule {
	return func(flag *pflag.Flag, cmd *cobra.Command) error {
		for _, w := range strings.Split(flag.Name, string(delim)) {
			if ok := vocab.Spell(w); !ok {
				return fmt.Errorf("flag name should consist of delimited real english words for %s on %s",
					flag.Name, FullCommand(cmd))
			}
		}
		return nil
	}
}

// RequireFlagDelimiter checks that a flag uses a specified delimiter at most maxCount times
func RequireFlagDelimiter(delim rune, maxCount int) FlagRule {
	return func(flag *pflag.Flag, cmd *cobra.Command) error {
		countDelim := 0
		for _, l := range flag.Name {
			if l == delim {
				countDelim++
				if countDelim > maxCount {
					return fmt.Errorf("flag name must only have %d delimiter (\"%c\") for %s on %s",
						maxCount, delim, flag.Name, FullCommand(cmd))
				}
			}
		}
		return nil
	}
}

// RequireFlagCharacters checks that a flag consists only of letters and a delimiter
func RequireFlagCharacters(delim rune) FlagRule {
	return func(flag *pflag.Flag, cmd *cobra.Command) error {
		for _, l := range flag.Name {
			if !unicode.IsLetter(l) && l != delim {
				return fmt.Errorf("flag name must be letters and delim (\"%c\") for %s on %s",
					delim, flag.Name, FullCommand(cmd))
			}
		}
		return nil
	}
}

// RequireFlagNotEndWithPunctuation checks that a flag description doesn't end with a period
func RequireFlagNotEndWithPunctuation(flag *pflag.Flag, cmd *cobra.Command) error {
	if flag.Usage[len(flag.Usage)-1] == '.' {
		return fmt.Errorf("flag usage ends with punctuation for %s on %s", flag.Name, FullCommand(cmd))
	}
	return nil
}

// RequireFlagStartWithCapital checks that a flag description starts with a capital letter
func RequireFlagStartWithCapital(flag *pflag.Flag, cmd *cobra.Command) error {
	if flag.Usage[0] < 'A' || flag.Usage[0] > 'Z' {
		return fmt.Errorf("flag usage should start with a capital for %s on %s", flag.Name, FullCommand(cmd))
	}
	return nil
}

// RequireFlagNameLength checks that a flag is between a certain min and max length
func RequireFlagNameLength(minLength, maxLength int) FlagRule {
	return func(flag *pflag.Flag, cmd *cobra.Command) error {
		var issues *multierror.Error
		if len(flag.Name) < minLength {
			issue := fmt.Errorf("flag name is too short for %s on %s", flag.Name, FullCommand(cmd))
			issues = multierror.Append(issues, issue)
		}
		if len(flag.Name) > maxLength {
			issue := fmt.Errorf("flag name is too long for %s on %s", flag.Name, FullCommand(cmd))
			issues = multierror.Append(issues, issue)
		}
		return issues
	}
}

func getValueByName(obj interface{}, name string) string {
	return reflect.Indirect(reflect.ValueOf(obj)).FieldByName(name).String()
}

func normalizeDesc(field string) string {
	switch field {
	case "Use":
		return "command"
	case "Long":
		return "long description"
	case "Short":
		return "short description"
	default:
		return strings.ToLower(field)
	}
}