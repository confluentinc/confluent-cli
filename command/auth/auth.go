package auth

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/confluentinc/cli/command/common"
	chttp "github.com/confluentinc/cli/http"
	"github.com/confluentinc/cli/shared"
)

type commands struct {
	Commands []*cobra.Command
	config   *shared.Config
}

// New returns a list of auth-related Cobra commands.
func New(config *shared.Config) []*cobra.Command {
	cmd := &commands{config: config}
	cmd.init()
	return cmd.Commands
}

func (a *commands) init() {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login to a Confluent Control Plane.",
		RunE:  a.login,
	}
	loginCmd.Flags().String("url", "https://confluent.cloud", "Confluent Control Plane URL")
	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout of a Confluent Control Plane.",
		RunE:  a.logout,
	}
	a.Commands = []*cobra.Command{loginCmd, logoutCmd}
}

func (a *commands) login(cmd *cobra.Command, args []string) error {
	if cmd.Flags().Changed("url") {
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			return err
		}
		a.config.AuthURL = url
	}
	email, password, err := credentials()
	if err != nil {
		return err
	}

	client := chttp.NewClient(chttp.BaseClient, a.config.AuthURL, a.config.Logger)
	token, err := client.Auth.Login(email, password)
	if err != nil {
		err = shared.ConvertAPIError(err)
		if err == shared.ErrUnauthorized { // special case for login failure
			err = shared.ErrIncorrectAuth
		}
		return common.HandleError(err)
	}
	a.config.AuthToken = token

	client = chttp.NewClientWithJWT(context.Background(), a.config.AuthToken, a.config.AuthURL, a.config.Logger)
	user, err := client.Auth.User()
	if err != nil {
		return common.HandleError(shared.ConvertAPIError(err))
	}
	a.config.Auth = user

	a.createOrUpdateContext(user)

	err = a.config.Save()
	if err != nil {
		return errors.Wrap(err, "unable to save user auth")
	}
	fmt.Println("Logged in as", email)
	return nil
}

func (a *commands) logout(cmd *cobra.Command, args []string) error {
	a.config.AuthToken = ""
	a.config.Auth = nil
	err := a.config.Save()
	if err != nil {
		return errors.Wrap(err, "unable to delete user auth")
	}
	fmt.Println("You are now logged out")
	return nil
}

func credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your Confluent Cloud credentials:")

	fmt.Print("Email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Password: ")
	bytePassword, err := terminal.ReadPassword(0)
	fmt.Println()
	if err != nil {
		return "", "", err
	}
	password := string(bytePassword)

	return strings.TrimSpace(email), strings.TrimSpace(password), nil
}

func (a *commands) createOrUpdateContext(user *shared.AuthConfig) {
	name := fmt.Sprintf("login-%s-%s", user.User.Email, a.config.AuthURL)
	if _, ok := a.config.Platforms[name]; !ok {
		a.config.Platforms[name] = &shared.Platform{
			Server: a.config.AuthURL,
		}
	}
	if _, ok := a.config.Credentials[name]; !ok {
		a.config.Credentials[name] = &shared.Credential{
			Username: user.User.Email,
			// don't save password if they entered it interactively
		}
	}
	if _, ok := a.config.Contexts[name]; !ok {
		a.config.Contexts[name] = &shared.Context{
			Platform:   name,
			Credential: name,
		}
	}
	a.config.CurrentContext = name
}
