package googleapps

import (
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/cli/explorer"
	"github.com/watermint/dcfg/integration/auth"
	"github.com/watermint/dcfg/integration/context"
	"google.golang.org/api/admin/directory/v1"
)

type GoogleApps interface {
	Preload()

	Users() []*admin.User
	Groups() []*admin.Group
	GroupMembers(groupEmail string) []*admin.Member
	CustomerUsers(customerId string) []*admin.User
}

func NewGoogleApps(ctx context.ExecutionContext) GoogleApps {
	impl := &GoogleAppsImpl{
		ExecutionContext: ctx,
	}
	impl.Preload()
	cache := &GoogleAppsWithCache{
		Resolver: impl,
	}
	cache.Preload()
	return cache
}

func FindGroup(googleApps GoogleApps, groupKey string) (*admin.Group, bool) {
	for _, x := range googleApps.Groups() {
		if x.Id == groupKey || x.Email == groupKey {
			return x, true
		}
	}
	return nil, false
}

func UserEmails(user *admin.User) (primary string, emails []string) {
	switch user.Emails.(type) {
	case []interface{}:
		userEmails := user.Emails.([]interface{})
		for _, ue1 := range userEmails {
			switch ue1.(type) {
			case map[string]interface{}:
				ue2 := ue1.(map[string]interface{})
				ue3, ue3exist := ue2["address"]
				if ue3exist {
					switch ue3.(type) {
					case string:
						ue4 := ue3.(string)
						emails = append(emails, ue4)
					default:
						seelog.Warnf("Unexpected JSON structure: ue3[%v] type[%T]", ue3, ue3)
					}
				}
			default:
				seelog.Warnf("Unexpected JSON structure: userEmails[%v] type[%T]", ue1, ue1)
			}
		}
		//for _, userEmail := range userEmails {
		//	if email, exist := userEmail["address"]; exist {
		//		emails = append(emails, email)
		//	} else {
		//		seelog.Warnf("Unexpected JSON structure: user.Emails[%v]", user.Emails)
		//	}
		//}

	default:
		seelog.Warnf("Unexpected JSON structure: userEmails[%v] type[%T]", user.Emails, user.Emails)
	}

	return user.PrimaryEmail, emails
}

type GoogleAppsWithCache struct {
	// flags for lazy loading
	lazyUsers         bool
	lazyGroups        bool
	lazyGroupMembers  map[string]bool
	lazyCustomerUsers map[string]bool

	// cache
	cacheUsers         []*admin.User
	cacheGroups        []*admin.Group
	cacheGroupMembers  map[string][]*admin.Member
	cacheCustomerUsers map[string][]*admin.User

	// resolver
	Resolver GoogleApps
}

func (g *GoogleAppsWithCache) Preload() {
	g.lazyGroupMembers = make(map[string]bool)
	g.lazyCustomerUsers = make(map[string]bool)
	g.cacheGroupMembers = make(map[string][]*admin.Member)
	g.cacheCustomerUsers = make(map[string][]*admin.User)
}

func (g *GoogleAppsWithCache) Users() []*admin.User {
	if !g.lazyUsers {
		g.cacheUsers = g.Resolver.Users()
		g.lazyUsers = true
	}
	return g.cacheUsers
}

func (g *GoogleAppsWithCache) Groups() []*admin.Group {
	if !g.lazyGroups {
		g.cacheGroups = g.Resolver.Groups()
		g.lazyGroups = true
	}
	return g.cacheGroups
}

func (g *GoogleAppsWithCache) GroupMembers(groupEmail string) []*admin.Member {
	if v, e := g.lazyGroupMembers[groupEmail]; !v || !e {
		g.cacheGroupMembers[groupEmail] = g.Resolver.GroupMembers(groupEmail)
		g.lazyGroupMembers[groupEmail] = true
	}
	m, e := g.cacheGroupMembers[groupEmail]
	// Should not happen on runtime
	if !e {
		seelog.Errorf("Inconsistent cache state: groupEmail[%s]", groupEmail)
		explorer.FatalShutdown("Please file issue on project page")
	}

	return m
}

func (g *GoogleAppsWithCache) CustomerUsers(customerId string) []*admin.User {
	if v, e := g.lazyCustomerUsers[customerId]; !v || !e {
		g.cacheCustomerUsers[customerId] = g.Resolver.CustomerUsers(customerId)
		g.lazyCustomerUsers[customerId] = true
	}
	u, e := g.cacheCustomerUsers[customerId]
	// Should not happen on runtime
	if !e {
		seelog.Errorf("Inconsistent cache state: customerId[%s]", customerId)
		explorer.FatalShutdown("Please file issue on project page")
	}

	return u
}

const (
	googleLoadChunkSize = 200
)

type GoogleAppsImpl struct {
	ExecutionContext context.ExecutionContext
}

func (g *GoogleAppsImpl) Preload() {
}

func (g *GoogleAppsImpl) Users() []*admin.User {
	rawUsers := make([]*admin.User, 0, googleLoadChunkSize)
	client := g.ExecutionContext.GoogleClient

	seelog.Tracef("Loading Google Users")
	users, err := client.Users.List().MaxResults(googleLoadChunkSize).Customer(auth.GOOGLE_CUSTOMER_ID).Do()
	if err != nil {
		seelog.Errorf("Unable to load Google Users: Err[%v]", err)
		explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth google`")
	}
	seelog.Tracef("Google User loaded (chunk): %d user(s)", len(users.Users))
	rawUsers = append(rawUsers, users.Users...)
	token := users.NextPageToken
	for token != "" {
		seelog.Trace("Loading Google Users (with token)")
		users, err := client.Users.List().MaxResults(googleLoadChunkSize).PageToken(token).Customer(auth.GOOGLE_CUSTOMER_ID).Do()
		if err != nil {
			seelog.Errorf("Unable to load Google Users: Err[%v]", err)
			explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth google`")
		}
		seelog.Tracef("Google User loaded (chunk): %d user(s), token[%s]", len(users.Users), token)
		rawUsers = append(rawUsers, users.Users...)
		token = users.NextPageToken
	}
	return rawUsers
}

func (g *GoogleAppsImpl) Groups() []*admin.Group {
	rawGroups := make([]*admin.Group, 0, googleLoadChunkSize)
	client := g.ExecutionContext.GoogleClient

	seelog.Tracef("Loading Google Groups")
	groups, err := client.Groups.List().MaxResults(googleLoadChunkSize).Customer(auth.GOOGLE_CUSTOMER_ID).Do()
	if err != nil {
		seelog.Errorf("Unable to load Google Groups: Err[%v]", err)
		explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth google`")
	}
	seelog.Tracef("Google Group loaded (chunk): %d group(s)", len(groups.Groups))
	rawGroups = append(rawGroups, groups.Groups...)
	token := groups.NextPageToken
	for token != "" {
		seelog.Trace("Loading Google Groups (with token)")
		groups, err := client.Groups.List().MaxResults(googleLoadChunkSize).PageToken(token).Customer(auth.GOOGLE_CUSTOMER_ID).Do()
		if err != nil {
			seelog.Errorf("Unable to load Google Groups: Err[%v]", err)
			explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth google`")
		}
		seelog.Tracef("Google Groups loaded (chunk): %d groups(s), token[%s]", len(groups.Groups), token)
		rawGroups = append(rawGroups, groups.Groups...)
		token = groups.NextPageToken
	}

	return rawGroups
}

func (g *GoogleAppsImpl) GroupMembers(groupEmail string) []*admin.Member {
	rawMember := make([]*admin.Member, 0, googleLoadChunkSize)
	seelog.Tracef("Loading members of Google Group: GroupKey[%s]", groupEmail)
	client := g.ExecutionContext.GoogleClient

	m, err := client.Members.List(groupEmail).MaxResults(googleLoadChunkSize).Do()
	if err != nil {
		seelog.Errorf("Unable to load Google Group Member: err[%s]", err)
		explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth google`")
	}
	seelog.Tracef("Google Members of Group loaded: GroupKey[%s]: %d member(s)", groupEmail, len(m.Members))
	rawMember = append(rawMember, m.Members...)
	token := m.NextPageToken
	for token != "" {
		m, err := client.Members.List(groupEmail).MaxResults(googleLoadChunkSize).PageToken(token).Do()
		if err != nil {
			seelog.Errorf("Unable to load Google Group member (with token): Err[%s]", err)
			explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth google`")
		}
		seelog.Tracef("Google Members of Group loaded: GroupKey[%s]: %d member(s)", groupEmail, len(m.Members))
		rawMember = append(rawMember, m.Members...)
		token = m.NextPageToken
	}
	return rawMember
}

func (g *GoogleAppsImpl) CustomerUsers(customerId string) []*admin.User {
	rawUsers := make([]*admin.User, 0, googleLoadChunkSize)
	client := g.ExecutionContext.GoogleClient
	seelog.Tracef("Loading Google Customer Members: CustomerId[%s]", customerId)

	r, err := client.Users.List().Customer(customerId).MaxResults(googleLoadChunkSize).Do()
	if err != nil {
		seelog.Errorf("Unable to load Google member in Customer: CustomerId[%s]", customerId)
		explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth google`")
	}
	seelog.Tracef("Google Customer Member loaded (chunk): %d", len(r.Users))
	rawUsers = append(rawUsers, r.Users...)
	token := r.NextPageToken

	for token != "" {
		r, err := client.Users.List().Customer(customerId).MaxResults(googleLoadChunkSize).PageToken(token).Do()
		if err != nil {
			seelog.Errorf("Unable to load Google member in Customer: CustomerId[%s]", customerId)
			explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth google`")
		}
		seelog.Tracef("Google Customer Member loaded (chunk): %d", len(r.Users))
		rawUsers = append(rawUsers, r.Users...)
		token = r.NextPageToken
	}
	return rawUsers
}
