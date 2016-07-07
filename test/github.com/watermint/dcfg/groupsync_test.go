package dcfg

import (
	"testing"
	"github.com/watermint/dcfg/integration/connector"
	"github.com/watermint/dcfg/integration/directory"
	"github.com/watermint/dcfg/sync/groupsync"
)

func TestGroupSync1(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{
			directory.Group{
				GroupId:   "g1@example.com",
				GroupName: "G1",
				Members: []directory.Account{
					directory.Account{
						Email: "a@example.com",
					},
				},
			},
		},
	}
	dropboxGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{
			directory.Group{
				GroupId:   "g1",
				GroupName: "G1",
				Members: []directory.Account{
					directory.Account{
						Email: "a@example.com",
					},
				},
				CorrelationId: "g1@example.com",
			},
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

	groupSync := groupsync.GroupSync{
		DropboxConnector:        &provision,
		DropboxAccountDirectory: &dropboxAccounts,
		DropboxGroupDirectory:   &dropboxGroups,
		GoogleDirectory:         &googleGroups,
	}

	groupSync.Sync("g1@example.com")
	unexpected, missing, success := provision.AssertLogs([]string{})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestGroupSync2(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{
			directory.Group{
				GroupId:   "g1@example.com",
				GroupName: "G1",
				Members: []directory.Account{
					directory.Account{
						Email: "a@example.com",
					},
				},
			},
		},
	}
	dropboxGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{
			directory.Group{
				GroupId:   "g1",
				GroupName: "G1",
				Members: []directory.Account{
					directory.Account{
						Email: "a@example.com",
					},
					directory.Account{
						Email: "b@example.com",
					},
				},
				CorrelationId: "g1@example.com",
			},
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
	groupSync := groupsync.GroupSync{
		DropboxConnector:        &provision,
		DropboxAccountDirectory: &dropboxAccounts,
		DropboxGroupDirectory:   &dropboxGroups,
		GoogleDirectory:         &googleGroups,
	}

	groupSync.Sync("g1@example.com")
	unexpected, missing, success := provision.AssertLogs([]string{
		provision.CreateOperationLog("GroupsMembersRemove", "g1", "b@example.com"),
	})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestGroupSync3(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{
			directory.Group{
				GroupId:   "g1@example.com",
				GroupName: "G1",
				Members: []directory.Account{
					directory.Account{
						Email: "b@example.com",
					},
				},
			},
		},
	}
	dropboxGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{
			directory.Group{
				GroupId:   "g1",
				GroupName: "G1",
				Members: []directory.Account{
					directory.Account{
						Email: "a@example.com",
					},
					directory.Account{
						Email: "b@example.com",
					},
				},
				CorrelationId: "g1@example.com",
			},
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

	groupSync := groupsync.GroupSync{
		DropboxConnector:        &provision,
		DropboxAccountDirectory: &dropboxAccounts,
		DropboxGroupDirectory:   &dropboxGroups,
		GoogleDirectory:         &googleGroups,
	}

	groupSync.Sync("g1@example.com")

	unexpected, missing, success := provision.AssertLogs([]string{
		provision.CreateOperationLog("GroupsMembersRemove", "g1", "a@example.com"),
	})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestGroupSync4(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{
			directory.Group{
				GroupId:   "g1@example.com",
				GroupName: "G1",
				Members: []directory.Account{
					directory.Account{
						Email: "c@example.com",
					},
				},
			},
		},
	}
	dropboxGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{
			directory.Group{
				GroupId:   "g1",
				GroupName: "G1",
				Members: []directory.Account{
					directory.Account{
						Email: "a@example.com",
					},
					directory.Account{
						Email: "b@example.com",
					},
				},
				CorrelationId: "g1@example.com",
			},
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

	groupSync := groupsync.GroupSync{
		DropboxConnector:        &provision,
		DropboxAccountDirectory: &dropboxAccounts,
		DropboxGroupDirectory:   &dropboxGroups,
		GoogleDirectory:         &googleGroups,
	}

	groupSync.Sync("g1@example.com")

	unexpected, missing, success := provision.AssertLogs([]string{
		provision.CreateOperationLog("GroupsMembersRemove", "g1", "a@example.com"),
		provision.CreateOperationLog("GroupsMembersRemove", "g1", "b@example.com"),
		provision.CreateOperationLog("GroupsMembersAdd", "g1", "c@example.com"),
	})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestGroupSync5(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{
			directory.Group{
				GroupId:   "g1@example.com",
				GroupName: "G1",
				Members: []directory.Account{
					directory.Account{
						Email: "c@example.com",
					},
				},
			},
		},
	}
	dropboxGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{},
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

	groupSync := groupsync.GroupSync{
		DropboxConnector:        &provision,
		DropboxAccountDirectory: &dropboxAccounts,
		DropboxGroupDirectory:   &dropboxGroups,
		GoogleDirectory:         &googleGroups,
	}

	groupSync.Sync("g1@example.com")

	unexpected, missing, success := provision.AssertLogs([]string{
		provision.CreateOperationLog("GroupsCreate", "G1", "g1@example.com"),
		provision.CreateOperationLog("GroupsMembersAdd", "mock-g1@example.com", "c@example.com"),
	})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestGroupSync6(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{
			directory.Group{
				GroupId:   "g1@example.com",
				GroupName: "G1",
				Members: []directory.Account{
					directory.Account{
						Email: "d@example.com",
					},
				},
			},
		},
	}
	dropboxGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{},
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

	groupSync := groupsync.GroupSync{
		DropboxConnector:        &provision,
		DropboxAccountDirectory: &dropboxAccounts,
		DropboxGroupDirectory:   &dropboxGroups,
		GoogleDirectory:         &googleGroups,
	}

	groupSync.Sync("g1@example.com")

	unexpected, missing, success := provision.AssertLogs([]string{
		provision.CreateOperationLog("GroupsCreate", "G1", "g1@example.com"),
	})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestGroupSync7(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{
			directory.Group{
				GroupId:   "g1@example.com",
				GroupName: "G1-RENAME",
				Members: []directory.Account{
					directory.Account{
						Email: "c@example.com",
					},
				},
			},
		},
	}
	dropboxGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{
			directory.Group{
				GroupId:   "g1",
				GroupName: "G1",
				Members: []directory.Account{
					directory.Account{
						Email: "c@example.com",
					},
				},
				CorrelationId: "g1@example.com",
			},
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

	groupSync := groupsync.GroupSync{
		DropboxConnector:        &provision,
		DropboxAccountDirectory: &dropboxAccounts,
		DropboxGroupDirectory:   &dropboxGroups,
		GoogleDirectory:         &googleGroups,
	}

	groupSync.Sync("g1@example.com")

	unexpected, missing, success := provision.AssertLogs([]string{
		provision.CreateOperationLog("GroupsUpdate", "g1", "G1-RENAME"),
	})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}

func TestGroupSync8(t *testing.T) {
	provision := connector.DropboxConnectorMock{}
	googleGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{},
	}
	dropboxGroups := directory.GroupDirectoryMock{
		MockData: []directory.Group{
			directory.Group{
				GroupId:   "g1",
				GroupName: "G1",
				Members: []directory.Account{
					directory.Account{
						Email: "c@example.com",
					},
				},
				CorrelationId: "g1@example.com",
			},
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

	groupSync := groupsync.GroupSync{
		DropboxConnector:        &provision,
		DropboxAccountDirectory: &dropboxAccounts,
		DropboxGroupDirectory:   &dropboxGroups,
		GoogleDirectory:         &googleGroups,
	}

	groupSync.Sync("g1@example.com")

	unexpected, missing, success := provision.AssertLogs([]string{})
	if !success {
		t.Error("Sync failed", unexpected, missing, success)
	}
}
