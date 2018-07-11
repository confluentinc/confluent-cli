package shared

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
	jsonParser "github.com/hashicorp/hcl/json/parser"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"

	"github.com/confluentinc/cli/log"
	"path"
)

const (
	defaultConfigFile = "~/.confluent/config.hcl"
)

// ErrNoConfig means that no configuration exists.
var ErrNoConfig = fmt.Errorf("no config file exists")

// Config represents the CLI configuration.
type Config struct {
	MetricSink     MetricSink             `json:"-" hcl:"-"`
	Logger         *log.Logger            `json:"-" hcl:"-"`
	Filename       string                 `json:"-" hcl:"-"`
	AuthURL        string                 `json:"auth_url" hcl:"auth_url"`
	AuthToken      string                 `json:"auth_token" hcl:"auth_token"`
	Auth           *AuthConfig            `json:"auth" hcl:"auth"`
	Platforms      map[string]*Platform   `json:"platforms" hcl:"platforms"`
	Credentials    map[string]*Credential `json:"credentials" hcl:"credentials"`
	Contexts       map[string]*Context    `json:"contexts" hcl:"contexts"`
	CurrentContext string                 `json:"current_context" hcl:"current_context"`
}

// Load reads the CLI config from disk.
func (c *Config) Load() error {
	c.Platforms = map[string]*Platform{}
	c.Credentials = map[string]*Credential{}
	c.Contexts = map[string]*Context{}
	filename, err := c.getFilename()
	if err != nil {
		return err
	}
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNoConfig
		}
		return errors.Wrapf(err, "unable to read config file: %s", filename)
	}
	err = hcl.Unmarshal(input, c)
	if err != nil {
		return errors.Wrapf(err, "unable to parse config file: %s", filename)
	}
	return nil
}

// Save writes the CLI config to disk.
func (c *Config) Save() error {
	cfg, err := json.Marshal(c)
	if err != nil {
		return errors.Wrapf(err, "unable to marshal config")
	}
	ast, err := jsonParser.Parse(cfg)
	if err != nil {
		return errors.Wrapf(err, "unable to parse config")
	}
	filename, err := c.getFilename()
	if err != nil {
		return err
	}
	err = os.MkdirAll(path.Dir(filename), 0700)
	if err != nil {
		return errors.Wrapf(err, "unable to create config directory: %s", filename)
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrapf(err, "unable to create config file: %s", filename)
	}
	defer func() { _ = f.Close() }()
	err = printer.Fprint(f, ast)
	if err != nil {
		return errors.Wrapf(err, "unable to write config to file: %s", filename)
	}
	return nil
}

// Context returns the current Context object.
func (c *Config) Context() (*Context, error) {
	if c.CurrentContext == "" {
		return nil, ErrUnauthorized
	}
	return c.Contexts[c.CurrentContext], nil
}

// CheckLogin returns an error if the user is not logged in.
func (c *Config) CheckLogin() error {
	if c.Auth == nil || c.Auth.Account == nil || c.Auth.Account.Id == "" {
		return ErrUnauthorized
	}
	return nil
}

func (c *Config) getFilename() (string, error) {
	if c.Filename == "" {
		c.Filename = defaultConfigFile
	}
	filename, err := homedir.Expand(c.Filename)
	if err != nil {
		return "", err
	}
	return filename, nil
}
