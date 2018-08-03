package common

import (
	"fmt"

	"github.com/codyaray/go-editor"

	"github.com/confluentinc/cli/shared"
)

// HandleError provides standard error messaging for common errors.
func HandleError(err error) error {
	switch err {
	case shared.ErrNoContext:
		fallthrough
	case shared.ErrUnauthorized:
		fmt.Println("You must login to access Confluent Cloud.")
	case shared.ErrExpiredToken:
		fmt.Println("Your access to Confluent Cloud has expired. Please login again.")
	case shared.ErrIncorrectAuth:
		fmt.Println("You have entered an incorrect username or password. Please try again.")
	case shared.ErrMalformedToken:
		fmt.Println("Your auth token has been corrupted. Please login again.")
	case shared.ErrNotImplemented:
		fmt.Println("Sorry, this functionality is not yet available in the CLI.")
	case shared.ErrNotFound:
		fmt.Println("Kafka cluster not found.")  // TODO: parametrize ErrNotFound for better error messaging
	default:
		switch err.(type) {
		case editor.ErrEditing:
			fmt.Println(err)
		case shared.ErrNotAuthenticated:
			fmt.Println(err)
		case shared.ErrKafka:
			fmt.Println(err)
		default:
			return err
		}
	}
	return nil
}
