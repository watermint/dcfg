package dcfg

import (
	"testing"
	"github.com/watermint/dcfg/integration/connector"
)

func expectDropboxConnector(connector.DropboxConnector) {
}

func TestConnectorInterface(*testing.T) {
	expectDropboxConnector(&connector.DropboxConnectorImpl{})
	expectDropboxConnector(&connector.DropboxConnectorMock{})
}
