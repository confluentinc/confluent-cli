package auth

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	chttp "github.com/confluentinc/ccloud-sdk-go"
	sdkMock "github.com/confluentinc/ccloud-sdk-go/mock"
	orgv1 "github.com/confluentinc/ccloudapis/org/v1"
	cliMock "github.com/confluentinc/cli/mock"

	"github.com/confluentinc/cli/command"
	"github.com/confluentinc/cli/log"
	"github.com/confluentinc/cli/shared"
)

func TestCredentialsOverride(t *testing.T) {
	req := require.New(t)
	currentEmail := os.Getenv("XX_CCLOUD_EMAIL")
	currentPassword := os.Getenv("XX_CCLOUD_PASSWORD")

	os.Setenv("XX_CCLOUD_EMAIL", "test-email")
	os.Setenv("XX_CCLOUD_PASSWORD", "test-password")

	prompt := prompt("cody@confluent.io", "iambatman")
	auth := &sdkMock.Auth{
		LoginFunc: func(ctx context.Context, username string, password string) (string, error) {
			return "y0ur.jwt.T0kEn", nil
		},
		UserFunc: func(ctx context.Context) (*orgv1.GetUserReply, error) {
			return &orgv1.GetUserReply{
				User: &orgv1.User{
					Id:        23,
					Email:     "test-email",
					FirstName: "Cody",
				},
			}, nil
		},
	}
	cmds, config := newAuthCommand(prompt, auth, req)

	output, err := command.ExecuteCommand(cmds.Commands[0])
	req.NoError(err)
	req.Contains(output, "Logged in as test-email")

	req.Equal("y0ur.jwt.T0kEn", config.AuthToken)
	req.Equal(&orgv1.User{Id: 23, Email: "test-email", FirstName: "Cody"}, config.Auth.User)

	os.Setenv("XX_CCLOUD_EMAIL", currentEmail)
	os.Setenv("XX_CCLOUD_PASSWORD", currentPassword)
}

func TestLoginSuccess(t *testing.T) {
	req := require.New(t)

	prompt := prompt("cody@confluent.io", "iambatman")
	auth := &sdkMock.Auth{
		LoginFunc: func(ctx context.Context, username string, password string) (string, error) {
			return "y0ur.jwt.T0kEn", nil
		},
		UserFunc: func(ctx context.Context) (*orgv1.GetUserReply, error) {
			return &orgv1.GetUserReply{
				User: &orgv1.User{
					Id:        23,
					Email:     "cody@confluent.io",
					FirstName: "Cody",
				},
			}, nil
		},
	}
	cmds, config := newAuthCommand(prompt, auth, req)

	output, err := command.ExecuteCommand(cmds.Commands[0])
	req.NoError(err)
	req.Contains(output, "Logged in as cody@confluent.io")

	req.Equal("y0ur.jwt.T0kEn", config.AuthToken)
	req.Equal(&orgv1.User{Id: 23, Email: "cody@confluent.io", FirstName: "Cody"}, config.Auth.User)

	config = shared.NewConfig()
	req.NoError(config.Load())
	name := "login-cody@confluent.io-https://confluent.cloud"
	req.Contains(config.Platforms, name)
	req.Equal("https://confluent.cloud", config.Platforms[name].Server)
	req.Contains(config.Credentials, name)
	req.Equal("cody@confluent.io", config.Credentials[name].Username)
	req.Contains(config.Contexts, name)
	req.Equal(name, config.Contexts[name].Platform)
	req.Equal(name, config.Contexts[name].Credential)
}

func TestLoginFail(t *testing.T) {
	req := require.New(t)

	prompt := prompt("cody@confluent.io", "iamrobin")
	auth := &sdkMock.Auth{
		LoginFunc: func(ctx context.Context, username string, password string) (string, error) {
			return "", shared.ErrIncorrectAuth
		},
	}
	cmds, _ := newAuthCommand(prompt, auth, req)

	_, err := command.ExecuteCommand(cmds.Commands[0])
	req.Contains(err.Error(), "You have entered an incorrect username or password.")
}

func TestLogout(t *testing.T) {
	req := require.New(t)

	prompt := prompt("cody@confluent.io", "iamrobin")
	auth := &sdkMock.Auth{}
	cmds, config := newAuthCommand(prompt, auth, req)

	config.AuthToken = "some.token.here"
	config.Auth = &shared.AuthConfig{User: &orgv1.User{Id: 23}}
	req.NoError(config.Save())

	output, err := command.ExecuteCommand(cmds.Commands[1])
	req.NoError(err)
	req.Contains(output, "You are now logged out")

	config = shared.NewConfig()
	req.NoError(config.Load())
	req.Empty(config.AuthToken)
	req.Empty(config.Auth)
}

func Test_credentials_NoSpacesAroundEmail_ShouldSupportSpacesAtBeginOrEnd(t *testing.T) {
	req := require.New(t)

	prompt := prompt(" cody@confluent.io ", " iamrobin ")
	prompt.Out = os.Stdout
	auth := &sdkMock.Auth{}
	cmds, _ := newAuthCommand(prompt, auth, req)

	user, pass, err := cmds.credentials()
	req.NoError(err)
	req.Equal("cody@confluent.io", user)
	req.Equal(" iamrobin ", pass)
}

func prompt(username, password string) *cliMock.Prompt {
	return &cliMock.Prompt{
		Strings:   []string{username},
		Passwords: []string{password},
	}
}

func newAuthCommand(prompt command.Prompt, auth *sdkMock.Auth, req *require.Assertions) (*commands, *shared.Config) {
	var mockAnonHTTPClientFactory = func(baseURL string, logger *log.Logger) *chttp.Client {
		req.Equal("https://confluent.cloud", baseURL)
		return &chttp.Client{Auth: auth}
	}
	var mockJwtHTTPClientFactory = func(ctx context.Context, jwt, baseURL string, logger *log.Logger) *chttp.Client {
		return &chttp.Client{Auth: auth}
	}
	config := shared.NewConfig()
	config.Logger = log.New()
	commands := newCommands(config, prompt, mockAnonHTTPClientFactory, mockJwtHTTPClientFactory)
	for _, c := range commands.Commands {
		c.PersistentFlags().CountP("verbose", "v", "increase output verbosity")
	}
	return commands, config
}