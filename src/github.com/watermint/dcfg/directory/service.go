package directory

type AccountDirectory interface {
	Load()
	Accounts() []Account
}

type Account struct {
	Email     string
	GivenName string
	Surname   string
}

type Group struct {
	GroupId       string
	GroupName     string
	GroupEmail    string
	Members       []Account
	CorrelationId string
}

type GroupDirectory interface {
	Load()
	Groups() []Group
}

func ExistInDirectory(ad AccountDirectory, account Account) bool {
	for _, x := range ad.Accounts() {
		if account.Email == x.Email {
			return true
		}
	}
	return false
}

func ExistInGroup(group Group, account Account) bool {
	for _, x := range group.Members {
		if x.Email == account.Email {
			return true
		}
	}
	return false
}

func FindByCorrelationId(gd GroupDirectory, correlationId string) (Group, bool) {
	for _, x := range gd.Groups() {
		if x.CorrelationId == correlationId {
			return x, true
		}
	}
	return Group{}, false
}

func FindByGroupId(gd GroupDirectory, groupId string) (Group, bool) {
	for _, x := range gd.Groups() {
		if x.GroupId == groupId {
			return x, true
		}
	}
	return Group{}, false
}

type AccountDirectoryMock struct {
	MockData []Account
}

func (adm *AccountDirectoryMock) Load() {
}

func (adm *AccountDirectoryMock) Accounts() []Account {
	return adm.MockData
}

type GroupDirectoryMock struct {
	MockData []Group
}

func (gdm *GroupDirectoryMock) Load() {
}

func (gdm *GroupDirectoryMock) Groups() []Group {
	return gdm.MockData
}
