package dcfg

import (
	"testing"
	"dcfg/directory"
	"dcfg/usersync"
	"dcfg/connector"
)

func TestUserSyncAddUser(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleAccounts := directory.AccountDirectoryMock{
		MockData:[]directory.Account{
			directory.Account{
				Email: "a@example.com",
			},
			directory.Account{
				Email: "b@example.com",
			},
		},
	}
	dropboxAccounts := directory.AccountDirectoryMock{
		MockData:[]directory.Account{
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
	userSync := usersync.UserSync{
		DropboxConnector: &provision,
		DropboxAccounts: &dropboxAccounts,
		GoogleAccounts: &googleAccounts,
	}
	userSync.SyncProvision()

	unexpected, missing, success := provision.AssertLogs([]string {
	})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestUserSyncAddUser2(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleAccounts := directory.AccountDirectoryMock{
		MockData:[]directory.Account{
			directory.Account{
				Email: "a@example.com",
				GivenName: "Given-A",
				Surname: "Sur-A",
			},
			directory.Account{
				Email: "b@example.com",
				GivenName: "Given-B",
				Surname: "Sur-B",
			},
			directory.Account{
				Email: "c@example.com",
				GivenName: "Given-C",
				Surname: "Sur-C",
			},
		},
	}
	dropboxAccounts := directory.AccountDirectoryMock{
		MockData:[]directory.Account{
			directory.Account{
				Email: "a@example.com",
			},
			directory.Account{
				Email: "b@example.com",
			},
		},
	}
	userSync := usersync.UserSync{
		DropboxConnector: &provision,
		DropboxAccounts: &dropboxAccounts,
		GoogleAccounts: &googleAccounts,
	}
	userSync.SyncProvision()

	unexpected, missing, success := provision.AssertLogs([]string {
		provision.CreateOperationLog("MembersAdd", "c@example.com", "Given-C", "Sur-C"),
	})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestUserSyncAddUser3(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleAccounts := directory.AccountDirectoryMock{
		MockData:[]directory.Account{
			directory.Account{
				Email: "a@example.com",
				GivenName: "Given-A",
				Surname: "Sur-A",
			},
			directory.Account{
				Email: "b@example.com",
				GivenName: "Given-B",
				Surname: "Sur-B",
			},
			directory.Account{
				Email: "c@example.com",
				GivenName: "Given-C",
				Surname: "Sur-C",
			},
		},
	}
	dropboxAccounts := directory.AccountDirectoryMock{
		MockData:[]directory.Account{
			directory.Account{
				Email: "a@example.com",
			},
			directory.Account{
				Email: "d@example.com",
			},
		},
	}
	userSync := usersync.UserSync{
		DropboxConnector: &provision,
		DropboxAccounts: &dropboxAccounts,
		GoogleAccounts: &googleAccounts,
	}
	userSync.SyncProvision()

	unexpected, missing, success := provision.AssertLogs([]string {
		provision.CreateOperationLog("MembersAdd", "b@example.com", "Given-B", "Sur-B"),
		provision.CreateOperationLog("MembersAdd", "c@example.com", "Given-C", "Sur-C"),
	})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}
