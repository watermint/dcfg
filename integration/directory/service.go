package directory

type AccountDirectory interface {
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
	Groups() map[string]Group // groupId -> Group
}

type GroupResolver interface {
	// Find by group key. groupKey matches both GroupId and GroupEmail.
	Group(groupKey string) (Group, bool)
}

type EmailResolver interface {
	// Ensure email exist in the directory.
	EmailExist(email string) (bool, error)
}

type AccountDirectoryMock struct {
	MockData []Account
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

type EmailResolverMock struct {
	MockData []string
}

func (erm *EmailResolverMock) EmailExist(email string) (bool, error) {
	for _, x := range erm.MockData {
		if x == email {
			return true, nil
		}
	}
	return false, nil
}
