package directory

import (
	"github.com/cihub/seelog"
	"github.com/dropbox/dropbox-sdk-go-unofficial/team"
	"github.com/dropbox/dropbox-sdk-go-unofficial/team_common"
	"github.com/watermint/dcfg/cli/explorer"
	"github.com/watermint/dcfg/integration/context"
)

type DropboxDirectory struct {
	executionContext context.ExecutionContext

	// API raw data structure
	rawMembers        []*team.TeamMemberInfo
	rawGroupSummaries []*team_common.GroupSummary
	rawGroupFullInfo  map[string]*team.GroupFullInfo

	// Abstract data structure
	groups   map[string]Group   // GroupId -> Group
	accounts map[string]Account // Email -> Account
}

const (
	dropboxLoadChunkSize = 100
)

func NewDropboxDirectory(ctx context.ExecutionContext) *DropboxDirectory {
	dd := DropboxDirectory{
		executionContext: ctx,
	}
	dd.load()
	return &dd
}

func (d *DropboxDirectory) loadMembers() {
	d.rawMembers = []*team.TeamMemberInfo{}
	client := d.executionContext.DropboxClient

	seelog.Trace("Loading Dropbox Team Member Info")

	sel := team.MembersListArg{}
	sel.Limit = dropboxLoadChunkSize
	ms, err := client.MembersList(&sel)
	if err != nil {
		seelog.Errorf("Unable to load Dropbox Team Member: err[%s]", err)
		explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth dropbox`")
	}
	for _, m := range ms.Members {
		d.rawMembers = append(d.rawMembers, m)
	}
	seelog.Tracef("Dropbox Team Member Chunk loaded: %d member(s)", len(ms.Members))
	if !ms.HasMore {
		return
	}

	cursor := ms.Cursor
	for {
		sel := team.MembersListContinueArg{}
		sel.Cursor = cursor
		seelog.Trace("Loading Dropbox Team Member Info (Continue)")

		ms, err := client.MembersListContinue(&sel)
		if err != nil {
			seelog.Errorf("Unable to load Dropbox Team Member: err[%s]", err)
			explorer.FatalShutdown("Please re-run `-auth dropbox`")
		}
		for _, m := range ms.Members {
			d.rawMembers = append(d.rawMembers, m)
		}
		seelog.Tracef("Dropbox Team Member Chunk (Continue) loaded: %d member(s), has more: %t", len(ms.Members), ms.HasMore)
		if !ms.HasMore {
			break
		}
		cursor = ms.Cursor
	}
}

func (d *DropboxDirectory) loadGroupSummaries() {
	client := d.executionContext.DropboxClient

	seelog.Trace("Loading Dropbox Group Summaries")

	sel := team.GroupsListArg{}
	sel.Limit = dropboxLoadChunkSize
	gs, err := client.GroupsList(&sel)
	if err != nil {
		seelog.Errorf("Unable to load Dropbox Group Summary: Err[%s]", err)
		explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth dropbox`")
	}
	for _, g := range gs.Groups {
		d.rawGroupSummaries = append(d.rawGroupSummaries, g)
	}
	seelog.Tracef("Dropbox Group Summary Chunk loaded: %d group(s)", len(gs.Groups))
	if !gs.HasMore {
		return
	}
	cursor := gs.Cursor
	for {
		sel := team.GroupsListContinueArg{}
		sel.Cursor = cursor
		gs, err := client.GroupsListContinue(&sel)
		if err != nil {
			seelog.Errorf("Unable to load Dropbox Group Summary: Err[%s]", err)
			explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth dropbox`")
		}
		seelog.Tracef("Dropbox Group Summary (Continue) Chunk loaded: %d group(s)", len(gs.Groups))
		for _, g := range gs.Groups {
			d.rawGroupSummaries = append(d.rawGroupSummaries, g)
		}
		if !gs.HasMore {
			break
		}
		cursor = gs.Cursor
	}
}

func (d *DropboxDirectory) loadGroups() {
	groups := make(map[string]*team.GroupFullInfo)
	client := d.executionContext.DropboxClient

	for i, gs := range d.rawGroupSummaries {
		sel := team.GroupsSelector{}
		sel.Tag = "group_ids"
		sel.GroupIds = []string{gs.GroupId}
		seelog.Tracef("Loading Dropbox Group Full Info [%d of %d]: Group ID[%s] Group Name[%s]", i, len(d.rawGroupSummaries), gs.GroupId, gs.GroupName)
		results, err := client.GroupsGetInfo(&sel)

		if err != nil {
			seelog.Errorf("Failed to load Dropbox Group: GroupId[%s] GroupName[%s] Err[%v]", gs.GroupId, gs.GroupName, err)
			explorer.FatalShutdown("Please re-run `-sync` if it's network issue. If it looks like auth issue please re-run `-auth dropbox`")
		}

		for _, gr := range results {
			groups[gr.GroupInfo.GroupId] = gr.GroupInfo
		}
	}
	d.rawGroupFullInfo = groups
}

func (d *DropboxDirectory) load() {
	d.loadMembers()
	d.loadGroupSummaries()
	d.loadGroups()

	d.groups = d.createGroups()
	d.accounts = d.createAccounts()
}

func (d *DropboxDirectory) createAccounts() (members map[string]Account) {
	members = make(map[string]Account)
	for _, m := range d.rawMembers {
		members[m.Profile.Email] = Account{
			Email:     m.Profile.Email,
			GivenName: m.Profile.Name.GivenName,
			Surname:   m.Profile.Name.Surname,
		}
	}
	return
}

func (d *DropboxDirectory) createGroups() (groups map[string]Group) {
	groups = make(map[string]Group)
	for gid, g := range d.rawGroupFullInfo {
		members := map[string]Account{}
		for _, m := range g.Members {
			members[m.Profile.Email] = Account{
				Email:     m.Profile.Email,
				GivenName: m.Profile.Name.GivenName,
				Surname:   m.Profile.Name.Surname,
			}
		}
		group := Group{
			GroupId:       gid,
			GroupName:     g.GroupName,
			CorrelationId: g.GroupExternalId,
			Members:       members,
		}
		groups[gid] = group
	}
	return
}

func (d *DropboxDirectory) Groups() map[string]Group {
	return d.groups
}

func (d *DropboxDirectory) Accounts() map[string]Account {
	return d.accounts
}
