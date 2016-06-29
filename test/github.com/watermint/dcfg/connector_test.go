package dcfg

import (
	"github.com/watermint/dcfg/connector"
	"testing"
)

func expectDropboxConnector(connector.DropboxConnector) {
}

func TestConnectorInterface(*testing.T) {
	expectDropboxConnector(&connector.DropboxConnectorImpl{})
	expectDropboxConnector(&connector.DropboxConnectorMock{})
}
