package dcfg

import (
	"testing"
	"github.com/watermint/dcfg/connector"
)

func expectDropboxConnector(connector.DropboxConnector) {
}

func TestConnectorInterface(*testing.T) {
	expectDropboxConnector(&connector.DropboxConnectorImpl{})
	expectDropboxConnector(&connector.DropboxConnectorMock{})
}
