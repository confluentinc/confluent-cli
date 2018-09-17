package http

import (
	"net/http"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"

	orgv1 "github.com/confluentinc/cc-structs/kafka/org/v1"
  schedv1 "github.com/confluentinc/cc-structs/kafka/scheduler/v1"
	"github.com/confluentinc/cli/log"
)

// APIKeyService provides methods for managing API keys on Confluent Control Plane.
type APIKeyService struct {
	client *http.Client
	sling  *sling.Sling
	logger *log.Logger
}

// NewAPIKeyService returns a new APIKeyService.
func NewAPIKeyService(client *Client) *APIKeyService {
	return &APIKeyService{
		client: client.httpClient,
		logger: client.logger,
		sling:  client.sling,
	}
}

// Create makes a new API Key
func (s *APIKeyService) Create(key *orgv1.ApiKey) (*orgv1.ApiKey, *http.Response, error) {
	request := &schedv1.CreateApiKeyRequest{ApiKey: &schedv1.ApiKey{ApiKey: key}}
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
