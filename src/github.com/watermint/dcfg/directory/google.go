package directory

import (
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/auth"
	"github.com/watermint/dcfg/explorer"
	"google.golang.org/api/admin/directory/v1"
)

type GoogleDirectory struct {
	// API raw data structure
	rawUsers        []*admin.User

	// Abstract data structure
	accounts []Account
}

const (
	googleLoadChunkSize = 200
)

func (g *GoogleDirectory) Group(groupId string) (Group, bool) {
	seelog.Tracef("Loading Google Group: GroupId[%s]", groupId)
	rawGroup, exist := g.loadGroup(groupId)
	if !exist {
		seelog.Tracef("Group not found for GroupId[%s]", groupId)
		return Group{}, false
	}
	seelog.Tracef("Loading Google Group Member: GroupId[%s]", groupId)
	rawGroupMembers := g.loadRawGroupMembers(rawGroup.Email, rawGroup.Email)

	return g.createGroupFromRaw(rawGroup, rawGroupMembers), true
}

func (g *GoogleDirectory) createGroupFromRaw(rawGroup *admin.Group, rawMembers []*admin.Member) Group {
	members := []Account{}
	for _, x := range rawMembers {
		for _, y := range g.getFlattenMember(x, rawGroup.Email, 0) {
			members = g.appendMember(members, y)
		}
	}

	group := Group{
		GroupId:    rawGroup.Email,
		GroupEmail: rawGroup.Email,
		GroupName:  rawGroup.Name,
		Members:    members,
	}

	return group
}

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

	seelog.Tracef("Loading Google Users")

	users, err := client.Users.List().MaxResults(googleLoadChunkSize).Customer(auth.GOOGLE_CUSTOMER_ID).Do()
	if err != nil {
		seelog.Errorf("Unable to load Google Users: Err[%v]", err)
		explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth google`")
	}
	seelog.Tracef("Google User loaded (chunk): %d user(s)", len(users.Users))
	for _, u := range users.Users {
		g.rawUsers = append(g.rawUsers, u)
	}
	token := users.NextPageToken
	for token != "" {
		seelog.Trace("Loading Google Users (with token)")
		users, err := client.Users.List().MaxResults(googleLoadChunkSize).PageToken(token).Customer(auth.GOOGLE_CUSTOMER_ID).Do()
		if err != nil {
			seelog.Errorf("Unable to load Google Users: Err[%v]", err)
			explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth google`")
		}
		seelog.Tracef("Google User loaded (chunk): %d user(s), token[%s]", len(users.Users), token)
		for _, x := range users.Users {
			g.rawUsers = append(g.rawUsers, x)
		}
		token = users.NextPageToken
	}
}

func (g *GoogleDirectory) loadGroup(groupKey string) (*admin.Group, bool) {
	client := auth.GoogleClient()

	seelog.Tracef("Loading Google Groups: GroupKey[%s]", groupKey)

	rawGroup, err := client.Groups.Get(groupKey).Do()
	if err != nil {
		seelog.Tracef("Unable to load Google Group: err[%v]", err)
		return nil, false
	}
	seelog.Tracef("Google Group Loaded: GroupKey[%s] GroupEmail[%s] GroupEtag[%s]", groupKey, rawGroup.Email, rawGroup.Etag)
	return rawGroup, true
}

func (g *GoogleDirectory) loadRawGroupMembers(groupKey, parentGroupKey string) (members []*admin.Member) {
	seelog.Tracef("Loading members of Google Group: ParentKey[%s] GroupKey[%s]", parentGroupKey, groupKey)
	client := auth.GoogleClient()

	m, err := client.Members.List(groupKey).MaxResults(googleLoadChunkSize).Do()
	if err != nil {
		seelog.Errorf("Unable to load Google Group Member: err[%s]", err)
		explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth google`")
	}
	seelog.Tracef("Google Members of Group loaded: ParentKey[%s] GroupKey[%s]: %d member(s)", parentGroupKey, groupKey, len(m.Members))
	for _, x := range m.Members {
		members = append(members, x)
	}
	token := m.NextPageToken
	for token != "" {
		m, err := client.Members.List(groupKey).MaxResults(googleLoadChunkSize).PageToken(token).Do()
		if err != nil {
			seelog.Errorf("Unable to load Google Group member (with token): Err[%s]", err)
			explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth google`")
		}
		seelog.Tracef("Google Members of Group loaded: ParentKey[%s]/GroupKey[%s]: %d member(s)", parentGroupKey, groupKey, len(m.Members))
		for _, x := range m.Members {
			members = append(members, x)
		}
		token = m.NextPageToken
	}
	return
}

func (g *GoogleDirectory) loadCustomerMembers(customerId string) (members []Account) {
	client := auth.GoogleClient()

	seelog.Tracef("Loading Google Customer Members: CustomerId[%s]", customerId)

	r, err := client.Users.List().Customer(customerId).MaxResults(googleLoadChunkSize).Do()
	if err != nil {
		seelog.Errorf("Unable to load Google member in Customer: CustomerId[%s]", customerId)
		explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth google`")
	}
	seelog.Tracef("Google Customer Member loaded (chunk): %d", len(r.Users))
	for _, x := range r.Users {
		members = g.appendMember(members, Account{
			Email:     x.PrimaryEmail,
			GivenName: x.Name.GivenName,
			Surname:   x.Name.FamilyName,
		})
	}
	token := r.NextPageToken

	for token != "" {
		r, err := client.Users.List().Customer(customerId).MaxResults(googleLoadChunkSize).PageToken(token).Do()
		if err != nil {
			seelog.Errorf("Unable to load Google member in Customer: CustomerId[%s]", customerId)
			explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth google`")
		}
		seelog.Tracef("Google Customer Member loaded (chunk): %d", len(r.Users))
		for _, x := range r.Users {
			members = g.appendMember(members, Account{
				Email:     x.PrimaryEmail,
				GivenName: x.Name.GivenName,
				Surname:   x.Name.FamilyName,
			})
		}
		token = r.NextPageToken
	}
	return
}

func (g *GoogleDirectory) getFlattenMember(member *admin.Member, parentGroupKey string, nest int) (members []Account) {
	switch member.Type {
	case "USER":
		seelog.Tracef("Google Group: Loading user: Nest[%d] Parent[%s] UserEmail[%s]", nest, parentGroupKey, member.Email)

		//TODO: Fetch name for user-provision (for future enhancement like filter by group)
		members = g.appendMember(members, Account{
			Email: member.Email,
		})
	case "GROUP":
		seelog.Tracef("Google Group: Loading Group: Nest[%d] Parent[%s], ChildGroupEmail[%s]", nest, parentGroupKey, member.Email)
		childMembers := g.loadRawGroupMembers(member.Email, parentGroupKey)
		for _, x := range childMembers {
			y := g.getFlattenMember(x, member.Email, nest + 1)
			for _, z := range y {
				members = g.appendMember(members, z)
			}
		}
	case "CUSTOMER":
		seelog.Tracef("Google Group: Loading Customer: Nest[%d] Parent[%s] Customer[%s]", nest, parentGroupKey, member.Id)
		for _, y := range g.loadCustomerMembers(member.Id) {
			members = g.appendMember(members, y)
		}

	default:
		seelog.Warnf("Unknown Google Group Member type: Id[%s] Email[%s] Type[%s]", member.Id, member.Email, member.Type)
	}
	return
}

func (g *GoogleDirectory) Load() {
	g.loadUsers()

	g.accounts = g.createAccounts()
}

func (g *GoogleDirectory) createAccounts() (accounts []Account) {
	for _, u := range g.rawUsers {
		accounts = g.appendMember(accounts, Account{
			Email:     u.PrimaryEmail,
			GivenName: u.Name.GivenName,
			Surname:   u.Name.FamilyName,
		})
	}
	return
}

func (g *GoogleDirectory) Accounts() []Account {
	return g.accounts
}
