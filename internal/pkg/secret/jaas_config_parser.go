package secret

import (
	"fmt"
	"regexp"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/confluentinc/properties"
)

type JAASParserInterface interface {
	ConvertPropertiesToJAAS(props *properties.Properties, op string) (*properties.Properties, error)
	SetOriginalConfigKeys(props *properties.Properties)
}

// JAASParser represents a jaas value that is returned from Parse().
type JAASParser struct {
	JaasOriginalConfigKeys *properties.Properties
	JaasProps              *properties.Properties
	Path                   string
	WhitespaceKey          string
	tokenizer              scanner.Scanner
}

func NewJAASParser() *JAASParser {
	return &JAASParser{
		JaasOriginalConfigKeys: properties.NewProperties(),
		JaasProps:              properties.NewProperties(),
		WhitespaceKey:          "",
	}
}

func (j *JAASParser) updateJAASConfig(op string, key string, value string, config string) (string, error) {
	switch op {
	case DELETE:
		keyValuePattern := key + JAAS_VALUE_PATTERN
		pattern := regexp.MustCompile(keyValuePattern)
		delete := ""
		// check if value is in JAAS format
		if pattern.MatchString(config) {
			matched := pattern.FindString(config)
			if matched == "" {
				return "", fmt.Errorf("The configuration " + config + " not present in JAAS configuration.")
			}
			config = pattern.ReplaceAllString(config, delete)
			if strings.HasSuffix(matched, ";") {
				config = config + ";"
			}
		} else {
			keyValuePattern := key + PASSWORD_PATTERN // check if value is in Secrets format
			pattern := regexp.MustCompile(keyValuePattern)
			matched := pattern.FindString(config)
			if matched == "" {
				return "", fmt.Errorf("The configuration " + key + " not present in JAAS configuration.")
			}
			config = pattern.ReplaceAllString(config, delete)
		}
		break
	case UPDATE:
		keyValuePattern := key + JAAS_VALUE_PATTERN
		pattern := regexp.MustCompile(keyValuePattern)
		if pattern.MatchString(config) {
			replaceVal := key + j.WhitespaceKey + "=" + j.WhitespaceKey + value
			matched := pattern.FindString(config)
			config = pattern.ReplaceAllString(config, replaceVal)
			if strings.HasSuffix(matched, ";") {
				config = config + ";"
			}
		} else {
			add := SPACE + key + j.WhitespaceKey + "=" + j.WhitespaceKey + value
			config = strings.TrimSuffix(config, ";") + add + ";"
		}
		break
	default:
		return "", fmt.Errorf("The operation " + op + " is not supported.")
	}

	return config, nil
}

func (j *JAASParser) parseConfig(specialChar rune) (string, int, error) {
	configName := ""
	offset := -1
	if unicode.IsSpace(j.tokenizer.Peek()) {
		j.tokenizer.Scan()
		configName = j.tokenizer.TokenText()
		offset = j.tokenizer.Pos().Offset
	}

	for j.tokenizer.Peek() != scanner.EOF && !unicode.IsSpace(j.tokenizer.Peek()) && j.tokenizer.Peek() != specialChar {
		j.tokenizer.Scan()
		configName = configName + j.tokenizer.TokenText()
		if offset == -1 {
			offset = j.tokenizer.Pos().Offset
		}
	}
	err := validateConfig(configName)
	if err != nil {
		return "", offset, err
	}
	return configName, offset, nil
}

func validateConfig(config string) error {
	if config == "}" || config == "{" || config == ";" || config == "=" || config == "};" || config == "" || config == " " {
		return fmt.Errorf("Invalid JAAS configuration: expected a configuration name but received " + config)
	}

	return nil
}

func (j *JAASParser) ignoreBackslash() {
	tokenizer := j.tokenizer
	tokenizer.Scan()
	if tokenizer.TokenText() == "\\" {
		j.tokenizer.Scan()
	}
}

func (j *JAASParser) isClosingBracket() bool {
	// If it's whitespace move ahead
	tokenizer := j.tokenizer
	if unicode.IsSpace(tokenizer.Peek()) {
		tokenizer.Scan()
		if tokenizer.TokenText() == "}" {
			j.tokenizer.Scan()
			return true
		}
	} else if tokenizer.Peek() == '}' {
		j.tokenizer.Scan()
		return true
	}

	return false
}

func (j *JAASParser) parseControlFlag() error {
	j.tokenizer.Scan()
	val := j.tokenizer.TokenText()
	switch val {
	case CONTROL_FLAG_REQUIRED, CONTROL_FLAG_REQUISITE, CONTROL_FLAG_OPTIONAL, CONTROL_FLAG_SUFFICIENT:
		j.ignoreBackslash()
		return nil
	default:
		return fmt.Errorf("Invalid JAAS configuration: login module control flag is not specified.")
	}
}

func (j *JAASParser) ParseJAASConfigurationEntry(jaasConfig string, key string) (*properties.Properties, error) {
	j.tokenizer.Init(strings.NewReader(jaasConfig))
	_, _, parsedToken, parentKey, err := j.parseConfigurationEntry(key)
	if err != nil {
		return nil, err
	}
	j.JaasOriginalConfigKeys.DisableExpansion = true
	_, _, err = j.JaasOriginalConfigKeys.Set(key+KEY_SEPARATOR+parentKey, jaasConfig)
	if err != nil {
		return nil, err
	}

	return parsedToken, nil
}

func (j *JAASParser) SetOriginalConfigKeys(props *properties.Properties) {
	j.JaasOriginalConfigKeys.Merge(props)
}

func (j *JAASParser) ConvertPropertiesToJAAS(props *properties.Properties, op string) (*properties.Properties, error) {
	configKey := ""
	result := properties.NewProperties()
	result.DisableExpansion = true
	j.JaasOriginalConfigKeys.DisableExpansion = true
	for key, value := range props.Map() {
		keys := strings.Split(key, KEY_SEPARATOR)
		configKey = keys[CLASS_ID] + KEY_SEPARATOR + keys[PARENT_ID]
		jaas, ok := j.JaasOriginalConfigKeys.Get(configKey)
		if !ok {
			return nil, fmt.Errorf("Failed to convert the properties to a JAAS configuration.")
		}
		jaas, err := j.updateJAASConfig(op, keys[KEY_ID], value, jaas)
		if err != nil {
			return nil, err
		}
		_, _, err = j.JaasOriginalConfigKeys.Set(configKey, jaas)
		if err != nil {
			return nil, err
		}
		_, _, err = result.Set(keys[CLASS_ID], jaas)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (j *JAASParser) parseConfigurationEntry(prefixKey string) (int, int, *properties.Properties, string, error) {
	// Parse Parent Key
	parsedConfigs := properties.NewProperties()

	parentKey, startIndex, err := j.parseConfig('=')
	if err != nil {
		return 0, 0, nil, "", err
	}

	// Parse Control Flag
	err = j.parseControlFlag()
	if err != nil {
		return 0, 0, nil, "", err
	}

	key := ""
	for j.tokenizer.Peek() != scanner.EOF && j.tokenizer.Peek() != ';' {
		// Parse Key
		key, _, err = j.parseConfig('=')
		if err != nil {
			return 0, 0, nil, "", err
		}

		if j.tokenizer.Peek() == ' ' {
			j.WhitespaceKey = " "
		}

		// Parse =
		if j.tokenizer.Peek() == scanner.EOF || j.tokenizer.Scan() != '=' || j.tokenizer.TokenText() == "" {
			return 0, 0, nil, "", fmt.Errorf("Invalid JAAS configuration: value is not specified for the key " + key)
		}

		// Parse Value
		value := ""
		value, _, err = j.parseConfig(';')
		if err != nil {
			return 0, 0, nil, "", err
		}
		newKey := prefixKey + KEY_SEPARATOR + parentKey + KEY_SEPARATOR + key
		_, _, err := parsedConfigs.Set(newKey, value)
		if err != nil {
			return 0, 0, nil, "", err
		}
		j.ignoreBackslash()
	}
	if j.tokenizer.Scan() != ';' {
		return 0, 0, nil, "", fmt.Errorf("Invalid JAAS configuration: configuration is not terminated with a ';'")
	}
	endIndex := j.tokenizer.Pos().Offset

	return startIndex, endIndex, parsedConfigs, parentKey, nil
}