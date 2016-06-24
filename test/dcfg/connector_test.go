package dcfg

import (
	"testing"
	"dcfg/connector"
)

func expectDropboxConnector(connector.DropboxConnector) {
}

func TestConnectorInterface(*testing.T) {
	expectDropboxConnector(&connector.DropboxConnectorImpl{})
	expectDropboxConnector(&connector.DropboxConnectorMock{})
}
