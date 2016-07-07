package connector

import (
	"testing"
)

func expectDropboxConnector(DropboxConnector) {
}

func TestConnectorInterface(*testing.T) {
	expectDropboxConnector(&DropboxConnectorImpl{})
	expectDropboxConnector(&DropboxConnectorMock{})
}
