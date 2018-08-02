package http

import (
	"net/http"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"

	orgv1 "github.com/confluentinc/cc-structs/kafka/org/v1"
	"github.com/confluentinc/cli/log"
)

// ApiKeyService provides methods for managing API keys on Confluent Control Plane.
type ApiKeyService struct {
	client *http.Client
	sling  *sling.Sling
	logger *log.Logger
}

// NewAPIKeyService returns a new ApiKeyService.
func NewAPIKeyService(client *Client) *ApiKeyService {
	return &ApiKeyService{
		client: client.httpClient,
		logger: client.logger,
		sling:  client.sling,
	}
}

// Create makes a new API Key
func (s *ApiKeyService) Create(key *orgv1.ApiKey) (*orgv1.ApiKey, *http.Response, error) {
	request := &orgv1.CreateApiKeyRequest{ApiKey: key}
	reply := new(orgv1.CreateApiKeyReply)
	resp, err := s.sling.New().Post("/api/api_keys").BodyJSON(request).Receive(reply, reply)
	if err != nil {
		return nil, resp, errors.Wrap(err, "unable to create API key")
	}
	if reply.Error != nil {
		return nil, resp, errors.Wrap(reply.Error, "error creating API key")
	}
	return reply.ApiKey, resp, nil
}
