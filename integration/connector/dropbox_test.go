package connector

import "testing"

func TestDropboxConnectorMock_AssertLogs(t *testing.T) {
	mock := DropboxConnectorMock{}

	logGroupsCreate := mock.CreateOperationLog("GroupsCreate", "TEST-GRP", "test-grp@example.com")
	logMissing := mock.CreateOperationLog("MembersRemove", "test-grp@example.com")

	u, m, s := mock.AssertLogs([]string{})
	if !s || len(m) > 0 || len(u) > 0 {
		t.Errorf("Unexpected state: Unexpected[%v] Missing[%v] Success[%t]", u, m, s)
	}

	mock.GroupsCreate("TEST-GRP", "test-grp@example.com")

	u, m, s = mock.AssertLogs([]string{logGroupsCreate})
	if !s || len(m) > 0 || len(u) > 0 {
		t.Errorf("Unexpected state: Unexpected[%v] Missing[%v] Success[%t]", u, m, s)
	}

	// Missing
	u, m, s = mock.AssertLogs([]string{logMissing})
	if s {
		t.Errorf("Unexpected state: Unexpected[%v] Missing[%v] Success[%t]", u, m, s)
	}
	if len(m) != 1 || m[0] != logMissing {
		t.Errorf("Unexpected state: Unexpected[%v] Missing[%v] Success[%t]", u, m, s)
	}

	// Unexpected
	if len(u) != 1 || u[0] != logGroupsCreate {
		t.Errorf("Unexpected state: Unexpected[%v] Missing[%v] Success[%t]", u, m, s)
	}

	// Clear logs
	mock.ClearOperationHistory()
	u, m, s = mock.AssertLogs([]string{})
	if !s || len(m) > 0 || len(u) > 0 {
		t.Errorf("Unexpected state: Unexpected[%v] Missing[%v] Success[%t]", u, m, s)
	}
}