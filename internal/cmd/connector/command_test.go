package connector

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	schedv1 "github.com/confluentinc/cc-structs/kafka/scheduler/v1"
	opv1 "github.com/confluentinc/cc-structs/operator/v1"
	"github.com/confluentinc/ccloud-sdk-go"
	ccsdkmock "github.com/confluentinc/ccloud-sdk-go/mock"

	v3 "github.com/confluentinc/cli/internal/pkg/config/v3"
	cliMock "github.com/confluentinc/cli/mock"
)

const (
	connectorID   = "lcc-123"
	connectorName = "myTestConnector"
)

type ConnectTestSuite struct {
	suite.Suite
	conf               *v3.Config
	kafkaCluster       *schedv1.KafkaCluster
	connector          *schedv1.Connector
	connectorInfo      *opv1.ConnectorInfo
	connectMock        *ccsdkmock.Connect
	kafkaMock          *ccsdkmock.Kafka
	connectorExpansion *opv1.ConnectorExpansion
}

func (suite *ConnectTestSuite) SetupSuite() {
	suite.conf = v3.AuthenticatedCloudConfigMock()
	ctx := suite.conf.Context()
	suite.kafkaCluster = &schedv1.KafkaCluster{
		Id:         ctx.KafkaClusterContext.GetActiveKafkaClusterId(),
		Name:       "KafkaMock",
		AccountId:  "testAccount",
		Enterprise: true,
	}
	suite.connector = &schedv1.Connector{
		Name:           connectorName,
		Id:             connectorID,
		KafkaClusterId: suite.kafkaCluster.Id,
		AccountId:      "testAccount",
		Status:         schedv1.Connector_RUNNING,
		UserConfigs:    map[string]string{},
	}

	suite.connectorInfo = &opv1.ConnectorInfo{
		Name: connectorName,
		Type: "source",
	}

	suite.connectorExpansion = &opv1.ConnectorExpansion{
		Id: &opv1.ConnectorId{Id: connectorID},
		Info: &opv1.ConnectorInfo{
			Name:   connectorName,
			Type:   "Sink",
			Config: map[string]string{},
		},
		Status: &opv1.ConnectorStateInfo{Name: connectorName, Connector: &opv1.ConnectorState{State: "Running"},
			Tasks: []*opv1.TaskState{{Id: 1, State: "Running"}},
		}}

}

func (suite *ConnectTestSuite) SetupTest() {
	suite.kafkaMock = &ccsdkmock.Kafka{
		DescribeFunc: func(ctx context.Context, cluster *schedv1.KafkaCluster) (*schedv1.KafkaCluster, error) {
			return suite.kafkaCluster, nil
		},
	}
	suite.connectMock = &ccsdkmock.Connect{
		CreateFunc: func(arg0 context.Context, arg1 *schedv1.ConnectorConfig) (connector *opv1.ConnectorInfo, e error) {
			return suite.connectorInfo, nil
		},
		UpdateFunc: func(arg0 context.Context, arg1 *schedv1.ConnectorConfig) (info *opv1.ConnectorInfo, e error) {
			return suite.connectorInfo, nil
		},
		PauseFunc: func(arg0 context.Context, arg1 *schedv1.Connector) error {
			return nil
		},
		ResumeFunc: func(arg0 context.Context, arg1 *schedv1.Connector) error {
			return nil
		},
		DeleteFunc: func(arg0 context.Context, arg1 *schedv1.Connector) error {
			return nil
		},
		ListWithExpansionsFunc: func(arg0 context.Context, arg1 *schedv1.Connector, arg2 string) (expansions map[string]*opv1.ConnectorExpansion, e error) {
			return map[string]*opv1.ConnectorExpansion{connectorID: suite.connectorExpansion}, nil
		},
		GetExpansionByIdFunc: func(arg0 context.Context, arg1 *schedv1.Connector) (expansion *opv1.ConnectorExpansion, e error) {
			return suite.connectorExpansion, nil
		},
		GetExpansionByNameFunc: func(ctx context.Context, connector *schedv1.Connector) (expansion *opv1.ConnectorExpansion, e error) {
			return suite.connectorExpansion, nil
		},
		GetFunc: func(arg0 context.Context, arg1 *schedv1.Connector) (connector *opv1.ConnectorInfo, e error) {
			return suite.connectorInfo, nil
		},
	}

}

func (suite *ConnectTestSuite) newCMD() *cobra.Command {
	prerunner := cliMock.NewPreRunnerMock(&ccloud.Client{Connect: suite.connectMock, Kafka: suite.kafkaMock}, nil)
	cmd := New(prerunner, suite.conf)
	return cmd
}

func (suite *ConnectTestSuite) TestPauseConnector() {
	cmd := suite.newCMD()
	cmd.SetArgs(append([]string{"pause", connectorID}))
	err := cmd.Execute()
	req := require.New(suite.T())
	req.Nil(err)
	req.True(suite.connectMock.PauseCalled())
	retVal := suite.connectMock.PauseCalls()[0]
	req.Equal(retVal.Arg1.KafkaClusterId, suite.kafkaCluster.Id)
}

func (suite *ConnectTestSuite) TestResumeConnector() {
	cmd := suite.newCMD()
	cmd.SetArgs(append([]string{"resume", connectorID}))
	err := cmd.Execute()
	req := require.New(suite.T())
	req.Nil(err)
	req.True(suite.connectMock.ResumeCalled())
	retVal := suite.connectMock.ResumeCalls()[0]
	req.Equal(retVal.Arg1.KafkaClusterId, suite.kafkaCluster.Id)
}

func (suite *ConnectTestSuite) TestDeleteConnector() {
	cmd := suite.newCMD()
	cmd.SetArgs(append([]string{"delete", connectorID}))
	err := cmd.Execute()
	req := require.New(suite.T())
	req.Nil(err)
	retVal := suite.connectMock.DeleteCalls()[0]
	req.Equal(retVal.Arg1.KafkaClusterId, suite.kafkaCluster.Id)
}

func (suite *ConnectTestSuite) TestListConnectors() {
	cmd := suite.newCMD()
	cmd.SetArgs(append([]string{"list"}))
	err := cmd.Execute()
	req := require.New(suite.T())
	req.Nil(err)
	req.True(suite.connectMock.ListWithExpansionsCalled())
	retVal := suite.connectMock.ListWithExpansionsCalls()[0]
	req.Equal(retVal.Arg1.KafkaClusterId, suite.kafkaCluster.Id)
}

func (suite *ConnectTestSuite) TestDescribeConnector() {
	cmd := suite.newCMD()
	cmd.SetArgs(append([]string{"describe", connectorID}))
	err := cmd.Execute()
	req := require.New(suite.T())
	req.Nil(err)
	req.True(suite.connectMock.GetExpansionByIdCalled())
	retVal := suite.connectMock.GetExpansionByIdCalls()[0]
	req.Equal(retVal.Arg1.KafkaClusterId, suite.kafkaCluster.Id)
}

func (suite *ConnectTestSuite) TestCreateConnector() {
	cmd := suite.newCMD()
	cmd.SetArgs(append([]string{"create", "--config", "../../../test/fixtures/input/connector-config.yaml"}))
	err := cmd.Execute()
	req := require.New(suite.T())
	req.Nil(err)
	req.True(suite.connectMock.CreateCalled())
	retVal := suite.connectMock.CreateCalls()[0]
	req.Equal(retVal.Arg1.KafkaClusterId, suite.kafkaCluster.Id)
}

func (suite *ConnectTestSuite) TestUpdateConnector() {
	cmd := suite.newCMD()
	cmd.SetArgs(append([]string{"update", connectorID, "--config", "../../../test/fixtures/input/connector-config.yaml"}))
	err := cmd.Execute()
	req := require.New(suite.T())
	req.Nil(err)
	req.True(suite.connectMock.UpdateCalled())
	retVal := suite.connectMock.UpdateCalls()[0]
	req.Equal(retVal.Arg1.KafkaClusterId, suite.kafkaCluster.Id)
}

func TestConnectTestSuite(t *testing.T) {
	suite.Run(t, new(ConnectTestSuite))
}