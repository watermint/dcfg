package usersync

import (
	"github.com/watermint/dcfg/integration/connector"
	"github.com/watermint/dcfg/integration/directory"
	"testing"
)

func TestUserSyncRemoveUser(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleEmail := directory.EmailResolverMock{
		MockData: []string{
			"a@example.com",
			"b@example.com",
		},
	}
	dropboxAccounts := directory.AccountDirectoryMock{
		MockData: []directory.Account{
			directory.Account{
				Email: "a@example.com",
			},
			directory.Account{
				Email: "b@example.com",
			},
			directory.Account{
				Email: "c@example.com",
			},
		},
	}
	userSync := UserSync{
		DropboxConnector: &provision,
		DropboxAccounts:  &dropboxAccounts,
		GoogleAccounts:   &directory.AccountDirectoryMock{},
		GoogleGroups:     &directory.GroupDirectoryMock{},
		GoogleEmail:      &googleEmail,
	}
	userSync.SyncDeprovision()

	unexpected, missing, success := provision.AssertLogs([]string{
		provision.CreateOperationLog("MembersRemove", "c@example.com"),
	})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestUserSyncRemoveUser2(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleEmail := directory.EmailResolverMock{
		MockData: []string{
			"a@example.com",
			"b@example.com",
		},
	}
	dropboxAccounts := directory.AccountDirectoryMock{
		MockData: []directory.Account{
			directory.Account{
				Email: "a@example.com",
			},
			directory.Account{
				Email: "c@example.com",
			},
			directory.Account{
				Email: "d@example.com",
			},
		},
	}
	userSync := UserSync{
		DropboxConnector: &provision,
		DropboxAccounts:  &dropboxAccounts,
		GoogleAccounts:   &directory.AccountDirectoryMock{},
		GoogleGroups:     &directory.GroupDirectoryMock{},
		GoogleEmail:      &googleEmail,
	}
	userSync.SyncDeprovision()

	unexpected, missing, success := provision.AssertLogs([]string{
		provision.CreateOperationLog("MembersRemove", "c@example.com"),
		provision.CreateOperationLog("MembersRemove", "d@example.com"),
	})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestUserSyncEqual(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleEmail := directory.EmailResolverMock{
		MockData: []string{
			"a@example.com",
			"b@example.com",
		},
	}
	dropboxAccounts := directory.AccountDirectoryMock{
		MockData: []directory.Account{
			directory.Account{
				Email: "a@example.com",
			},
			directory.Account{
				Email: "b@example.com",
			},
		},
	}
	userSync := UserSync{
		DropboxConnector: &provision,
		DropboxAccounts:  &dropboxAccounts,
		GoogleAccounts:   &directory.AccountDirectoryMock{},
		GoogleGroups:     &directory.GroupDirectoryMock{},
		GoogleEmail:      &googleEmail,
	}
	userSync.SyncDeprovision()

	unexpected, missing, success := provision.AssertLogs([]string{})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestUserSyncGoogleHasMore(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleEmail := directory.EmailResolverMock{
		MockData: []string{
			"a@example.com",
			"b@example.com",
			"c@example.com",
		},
	}
	dropboxAccounts := directory.AccountDirectoryMock{
		MockData: []directory.Account{
			directory.Account{
				Email: "a@example.com",
			},
			directory.Account{
				Email: "b@example.com",
			},
		},
	}
	userSync := UserSync{
		DropboxConnector: &provision,
		DropboxAccounts:  &dropboxAccounts,
		GoogleAccounts:   &directory.AccountDirectoryMock{},
		GoogleGroups:     &directory.GroupDirectoryMock{},
		GoogleEmail:      &googleEmail,
	}
	userSync.SyncDeprovision()

	unexpected, missing, success := provision.AssertLogs([]string{})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestUserSync_SyncDeprovision(t *testing.T) {
	emailsShouldKeep := []string{
		"a@example.com",
		"b@example.com",
		"b2@example.com",
		"b@example.net",
		"c@example.com",
		"d@example.com",
		"d@example.org",
		"tokyo@example.com",
		"minato@example.com",
		"meguro@example.com",
		"all@example.com",
	}
	emailsShouldRemove := []string{
		"x@example.com",
		"y@example.net",
		"z@example.org",
	}
	dropboxAccounts := make([]directory.Account, 0, len(emailsShouldKeep)+len(emailsShouldRemove))
	for _, x := range emailsShouldKeep {
		dropboxAccounts = append(dropboxAccounts, directory.Account{
			Email: x,
		})
	}
	for _, x := range emailsShouldRemove {
		dropboxAccounts = append(dropboxAccounts, directory.Account{
			Email: x,
		})
	}

	dp := &connector.DropboxConnectorMock{}
	dd := &directory.AccountDirectoryMock{
		MockData: dropboxAccounts,
	}
	gd := directory.CreateGoogleDirectoryForIntegrationTest()

	userSync := UserSync{
		DropboxConnector: dp,
		DropboxAccounts:  dd,
		GoogleAccounts:   gd,
		GoogleGroups:     gd,
		GoogleEmail:      gd,
	}
	userSync.SyncDeprovision()

	expectedOperations := make([]string, 0, len(emailsShouldRemove))
	for _, x := range emailsShouldRemove {
		expectedOperations = append(expectedOperations, dp.CreateOperationLog("MembersRemove", x))
	}
	unexpected, missing, success := dp.AssertLogs(expectedOperations)
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}
