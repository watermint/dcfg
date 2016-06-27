package directory

import (
	"google.golang.org/api/admin/directory/v1"
	"github.com/watermint/dcfg/auth"
	"github.com/watermint/dcfg/explorer"
	"github.com/cihub/seelog"
)

type GoogleDirectory struct {
	// Parameter
	Domain          string

	// API raw data structure
	rawUsers        []*admin.User
	rawGroups       []*admin.Group
	rawGroupMembers map[string][]*admin.Member

	// Abstract data structure
	groups          []Group
	accounts        []Account
}

const (
	googleLoadChunkSize = 100
)

func (g *GoogleDirectory) appendMember(members []Account, member Account) []Account {
	found := false
	for _, x := range members {
		if x.Email == member.Email {
			found = true
			break
		}
	}
	if !found {
		members = append(members, member)
	}
	return members
}

func (g *GoogleDirectory) loadUsers() {
	g.rawUsers = []*admin.User{}
	client := auth.GoogleClient()

	seelog.Tracef("Loading Google Users of domain[%s]", g.Domain)

	users, err := client.Users.List().MaxResults(googleLoadChunkSize).Domain(g.Domain).Do()
	if err != nil {
		explorer.Fatal("Unable to load Google Users", err)
	}
	seelog.Tracef("Google User loaded (chunk): %d user(s)", len(users.Users))
	for _, u := range users.Users {
		g.rawUsers = append(g.rawUsers, u)
	}
	token := users.NextPageToken
	for token != "" {
		seelog.Trace("Loading Google Users (with token)")
		users, err := client.Users.List().MaxResults(googleLoadChunkSize).PageToken(token).Domain(g.Domain).Do()
		if err != nil {
			explorer.Fatal("Unable to load Google Users (with token)", token, err)
		}
		seelog.Tracef("Google User loaded (chunk): %d user(s), token[%s]", len(users.Users), token)
		for _, x := range users.Users {
			g.rawUsers = append(g.rawUsers, x)
		}
		token = users.NextPageToken
	}
}

func (g *GoogleDirectory) loadGroup() {
	g.rawGroups = []*admin.Group{}
	client := auth.GoogleClient()

	seelog.Tracef("Loading Google Groups of domain[%s]", g.Domain)

	groups, err := client.Groups.List().MaxResults(googleLoadChunkSize).Domain(g.Domain).Do()
	if err != nil {
		explorer.Fatal("Unable to load Google Group", err)
	}
	seelog.Tracef("Google Group Loaded (chunk): %d", len(groups.Groups))
	for _, x := range groups.Groups {
		g.rawGroups = append(g.rawGroups, x)
	}
	token := groups.NextPageToken
	for token != "" {
		seelog.Trace("Loading Google Groups (with token)")
		groups, err := client.Groups.List().MaxResults(googleLoadChunkSize).Domain(g.Domain).PageToken(token).Do()
		if err != nil {
			explorer.Fatal("Unable to load Google Group (with token)", token, err)
		}
		seelog.Tracef("Google Group Loaded (chunk): %d", len(groups.Groups))
		for _, x := range groups.Groups {
			g.rawGroups = append(g.rawGroups, x)
		}
		token = groups.NextPageToken
	}
}

func (g *GoogleDirectory) loadMembers(group *admin.Group) (members []*admin.Member) {
	seelog.Tracef("Loading members of Google Group: GroupId[%s] GroupEmail[%s]", group.Id, group.Email)
	client := auth.GoogleClient()

	m, err := client.Members.List(group.Id).MaxResults(googleLoadChunkSize).Do()
	if err != nil {
		explorer.Fatal("Unable to load Google Group Member", err)
	}
	seelog.Tracef("Google Members of Group loaded (chunk): %d member(s)", len(m.Members))
	for _, x := range m.Members {
		members = append(members, x)
	}
	token := m.NextPageToken
	for token != "" {
		seelog.Trace("Loading Google Group Member (with token)")
		m, err := client.Members.List(group.Id).MaxResults(googleLoadChunkSize).PageToken(token).Do()
		if err != nil {
			explorer.Fatal("Unable to load Google Group member (with token)", err)
		}
		seelog.Tracef("Google Members of Group loaded (chunk): %d member(s)", len(m.Members))
		for _, x := range m.Members {
			members = append(members, x)
		}
		token = m.NextPageToken
	}
	return
}

func (g *GoogleDirectory) loadGroupMembers() {
	g.rawGroupMembers = make(map[string][]*admin.Member)
	for _, x := range g.rawGroups {
		g.rawGroupMembers[x.Email] = g.loadMembers(x)
	}
}

func (g *GoogleDirectory) loadCustomerMembers(customerId string) (members []Account) {
	client := auth.GoogleClient()

	seelog.Tracef("Loading Google Customer Members: CustomerId[%s]", customerId)

	r, err := client.Users.List().Customer(customerId).MaxResults(googleLoadChunkSize).Do()
	if err != nil {
		explorer.Fatal("Unable to load Google member in Customer: CustomerId[%s]", customerId)
	}
	seelog.Tracef("Google Customer Member loaded (chunk): %d", len(r.Users))
	for _, x := range r.Users {
		members = g.appendMember(members, Account{
			Email:x.PrimaryEmail,
			GivenName:x.Name.GivenName,
			Surname:x.Name.FamilyName,
		})
	}
	token := r.NextPageToken

	for token != "" {
		r, err := client.Users.List().Customer(customerId).MaxResults(googleLoadChunkSize).PageToken(token).Do()
		if err != nil {
			explorer.Fatal("Unable to load Google member in Customer: CustomerId[%s]", customerId)
		}
		seelog.Tracef("Google Customer Member loaded (chunk): %d", len(r.Users))
		for _, x := range r.Users {
			members = g.appendMember(members, Account{
				Email:x.PrimaryEmail,
				GivenName:x.Name.GivenName,
				Surname:x.Name.FamilyName,
			})
		}
		token = r.NextPageToken
	}
	return
}

func (g *GoogleDirectory) getFlattenMember(member *admin.Member) (members []Account) {
	switch member.Type {
	case "USER":
		seelog.Tracef("Google Group: Loading user: UserEmail[%s]", member.Email)

		//TODO: Fetch name for user-provision (for future enhancement like filter by group)
		members = g.appendMember(members, Account{
			Email: member.Email,
		})
	case "GROUP":
		seelog.Tracef("Google Group: Loading group: ChildGroupEmail[%s]", member.Email)
		for _, y := range g.getFlattenGroupMembers(member.Email) {
			members = g.appendMember(members, y)
		}
	case "CUSTOMER":
		seelog.Tracef("Google Group: Loading Customer: Customer[%s]", member.Id)
		for _, y := range g.loadCustomerMembers(member.Id) {
			members = g.appendMember(members, y)
		}

	default:
		seelog.Warnf("Unknown Google Group Member type: Id[%s] Email[%s] Type[%s]", member.Id, member.Email, member.Type)
	}
	return
}

func (g *GoogleDirectory) getFlattenGroupMembers(groupEmail string) (members []Account) {
	m, exist := g.rawGroupMembers[groupEmail]
	if !exist {
		explorer.Fatal("Google Group member not found", groupEmail)
	}
	seelog.Tracef("Loading flatten group members: GroupEmail[%s]", groupEmail)
	for _, x := range m {
		for _, y := range g.getFlattenMember(x) {
			members = g.appendMember(members, y)
		}
	}
	return
}

func (g *GoogleDirectory) Load() {
	g.loadUsers()
	g.loadGroup()
	g.loadGroupMembers()

	g.accounts = g.createAccounts()
	g.groups = g.createGroups()
}

func (g *GoogleDirectory) createAccounts() (accounts []Account) {
	for _, u := range g.rawUsers {
		accounts = g.appendMember(accounts, Account{
			Email: u.PrimaryEmail,
			GivenName: u.Name.GivenName,
			Surname: u.Name.FamilyName,
		})
	}
	return
}

func (g *GoogleDirectory) createGroups() (groups []Group) {
	for _, x := range g.rawGroups {
		group := Group{
			GroupId: x.Email,
			GroupName: x.Name,
			Members: g.getFlattenGroupMembers(x.Email),
		}
		groups = append(groups, group)
	}
	return
}

func (g *GoogleDirectory) Groups() []Group {
	return g.groups
}

func (g *GoogleDirectory) Accounts() []Account {
	return g.accounts
}
