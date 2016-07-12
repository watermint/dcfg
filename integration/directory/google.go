package directory

import (
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/integration/context"
	"github.com/watermint/dcfg/integration/directory/googleapps"
	"google.golang.org/api/admin/directory/v1"
)

type GoogleDirectory struct {
	ExecutionContext context.ExecutionContext

	googleApps googleapps.GoogleApps

	// All emails
	emailTypes map[string]int

	// Abstract data structure
	accounts map[string]Account
}

const (
	GOOGLE_EMAIL_TYPE_USER = iota
	GOOGLE_EMAIL_TYPE_GROUP
	GOOGLE_EMAIL_TYPE_ALIAS
)

func (g *GoogleDirectory) Group(groupKey string) (Group, bool) {
	seelog.Tracef("Loading Google Group: GroupId[%s]", groupKey)
	group, exist := googleapps.FindGroup(g.googleApps, groupKey)
	if !exist {
		return Group{}, false
	}

	return g.createGroup(group), true
}

func (g *GoogleDirectory) createGroup(rawGroup *admin.Group) Group {
	rawMembers := g.googleApps.GroupMembers(rawGroup.Email)

	members := map[string]Account{}
	for _, x := range rawMembers {
		for _, y := range g.extractMember(x, rawGroup.Email, 0) {
			members[y.Email] = y
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

func (g *GoogleDirectory) extractMember(member *admin.Member, parentGroupKey string, nest int) (members map[string]Account) {
	members = make(map[string]Account)
	switch member.Type {
	case "USER":
		seelog.Tracef("Google Group: Loading user: Nest[%d] Parent[%s] UserEmail[%s]", nest, parentGroupKey, member.Email)

		//TODO: Fetch name for user-provision (for future enhancement like filter by group)
		members[member.Email] = Account{
			Email: member.Email,
		}
	case "GROUP":
		seelog.Tracef("Google Group: Loading Group: Nest[%d] Parent[%s], ChildGroupEmail[%s]", nest, parentGroupKey, member.Email)
		childMembers := g.googleApps.GroupMembers(member.Email)
		for _, x := range childMembers {
			y := g.extractMember(x, member.Email, nest+1)
			for _, z := range y {
				members[z.Email] = z
			}
		}
	case "CUSTOMER":
		seelog.Tracef("Google Group: Loading Customer: Nest[%d] Parent[%s] Customer[%s]", nest, parentGroupKey, member.Id)
		for _, user := range g.googleApps.CustomerUsers(member.Id) {
			_, emails := googleapps.UserEmails(user)

			for _, e := range emails {
				members[e] = Account{
					Email: e,
				}
			}
		}

	default:
		seelog.Warnf("Unknown Google Group Member type: Id[%s] Email[%s] Type[%s]", member.Id, member.Email, member.Type)
	}
	return
}

func (g *GoogleDirectory) preloadEmails() {
	g.emailTypes = make(map[string]int)

	// Group emails
	for _, group := range g.googleApps.Groups() {
		g.emailTypes[group.Email] = GOOGLE_EMAIL_TYPE_GROUP
	}

	for _, user := range g.googleApps.Users() {
		primary, emails := googleapps.UserEmails(user)

		for _, e := range emails {
			g.emailTypes[e] = GOOGLE_EMAIL_TYPE_ALIAS
		}

		// overwrite primary email
		g.emailTypes[primary] = GOOGLE_EMAIL_TYPE_USER
	}
}

func (g *GoogleDirectory) Load() {
	g.googleApps = googleapps.NewGoogleApps(g.ExecutionContext)
	g.preloadEmails()

	g.accounts = g.createAccounts()
}

func (g *GoogleDirectory) createAccounts() (accounts map[string]Account) {
	accounts = make(map[string]Account)
	for _, u := range g.googleApps.Users() {
		accounts[u.PrimaryEmail] = Account{
			Email:     u.PrimaryEmail,
			GivenName: u.Name.GivenName,
			Surname:   u.Name.FamilyName,
		}
	}
	return
}

func (g *GoogleDirectory) Accounts() map[string]Account {
	return g.accounts
}

func (g *GoogleDirectory) EmailExist(email string) (bool, error) {
	_, e := g.emailTypes[email]
	return e, nil
}