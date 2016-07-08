package directory

type AccountDirectory interface {
	Load()
	Accounts() map[string]Account // email -> Account
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
	Members       map[string]Account // email -> Account
	CorrelationId string
}

type GroupDirectory interface {
	Load()
	Groups() map[string]Group // groupId -> Group
}

type GroupResolver interface {
	// Find by group key. groupKey matches both GroupId and GroupEmail.
	Group(groupKey string) (Group, bool)
}

type AccountDirectoryMock struct {
	MockData []Account
}

func (adm *AccountDirectoryMock) Load() {
}

func (adm *AccountDirectoryMock) Accounts() map[string]Account {
	accounts := make(map[string]Account)
	for _, x := range adm.MockData {
		accounts[x.Email] = x
	}
	return accounts
}

type GroupDirectoryMock struct {
	MockData []Group
}

func (gdm *GroupDirectoryMock) Load() {
}

func (gdm *GroupDirectoryMock) Groups() map[string]Group {
	groups := make(map[string]Group)
	for _, x := range gdm.MockData {
		groups[x.GroupId] = x
	}
	return groups
}

func (gdm *GroupDirectoryMock) Group(groupId string) (Group, bool) {
	for _, x := range gdm.MockData {
		if x.GroupId == groupId || x.GroupEmail == groupId {
			return x, true
		}
	}
	return Group{}, false
}
