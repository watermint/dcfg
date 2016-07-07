package groupsync

import (
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/connector"
	"github.com/watermint/dcfg/context"
	"github.com/watermint/dcfg/directory"
	"github.com/watermint/dcfg/explorer"
	"github.com/watermint/dcfg/text"
)

type GroupSync struct {
	DropboxConnector        connector.DropboxConnector
	DropboxAccountDirectory directory.AccountDirectory
	DropboxGroupDirectory   directory.GroupDirectory
	GoogleDirectory         directory.GroupResolver
}

func NewGroupSync(context context.ExecutionContext) GroupSync {
	gd := directory.GoogleDirectory{ExecutionContext: context}
	dd := directory.DropboxDirectory{ExecutionContext: context}
	dp := connector.CreateConnector(context)

	gd.Load()
	dd.Load()

	return GroupSync{
		DropboxConnector:        dp,
		DropboxAccountDirectory: &dd,
		DropboxGroupDirectory:   &dd,
		GoogleDirectory:         &gd,
	}
}

func (g *GroupSync) onDropboxGroupNotFound(googleGroup directory.Group) {
	g.DropboxConnector.GroupsCreate(googleGroup.GroupName, googleGroup.GroupId)
}

func (g *GroupSync) filterGoogleGroupMemberByAccountExistence(googleGroup directory.Group) (member []directory.Account) {
	for _, x := range googleGroup.Members {
		if directory.ExistInDirectory(g.DropboxAccountDirectory, x) {
			member = append(member, x)
		}
	}
	return
}

func (g *GroupSync) membersNotInGroup(member []directory.Account, group directory.Group) (notInGroup []directory.Account) {
	for _, x := range member {
		if !directory.ExistInGroup(group, x) {
			notInGroup = append(notInGroup, x)
		}
	}
	return
}

func (g *GroupSync) addMembersToDropboxGroup(dropboxGroupId string, member []string) {
}

func (g *GroupSync) syncNewGroup(googleGroup directory.Group) {
	newGroup := g.DropboxConnector.GroupsCreate(googleGroup.GroupName, googleGroup.GroupId)
	if newGroup != "" {
		for _, x := range g.filterGoogleGroupMemberByAccountExistence(googleGroup) {
			g.DropboxConnector.GroupsMembersAdd(newGroup, x.Email)
		}
	}
}

func (g *GroupSync) updateExistingGroup(googleGroup directory.Group, dropboxGroup directory.Group) {
	if googleGroup.GroupName != dropboxGroup.GroupName {
		g.DropboxConnector.GroupsUpdate(dropboxGroup.GroupId, googleGroup.GroupName)
	}

	googleMembers := g.filterGoogleGroupMemberByAccountExistence(googleGroup)

	notInDropboxGroup := g.membersNotInGroup(googleMembers, dropboxGroup)
	for _, x := range notInDropboxGroup {
		g.DropboxConnector.GroupsMembersAdd(dropboxGroup.GroupId, x.Email)
	}

	notInGoogleGroup := g.membersNotInGroup(dropboxGroup.Members, googleGroup)
	for _, x := range notInGoogleGroup {
		g.DropboxConnector.GroupsMembersRemove(dropboxGroup.GroupId, x.Email)
	}
}

func (g *GroupSync) Sync(targetGroup string) {
	seelog.Tracef("Group Sync from Google Group: Email[%s]", targetGroup)
	googleGroup, exist := directory.FindByGroupId(g.GoogleDirectory, targetGroup)
	if !exist {
		explorer.ReportFailure("Sync skipped for Google Group: %s (reason: Google Group not found)", targetGroup)
		seelog.Warnf("Google Group not found for sync: Email[%s]", targetGroup)
		return
	}

	dropboxGroup, exist := directory.FindByCorrelationId(g.DropboxGroupDirectory, googleGroup.GroupId)
	if !exist {
		g.syncNewGroup(googleGroup)
	} else {
		g.updateExistingGroup(googleGroup, dropboxGroup)
	}
}

func (g *GroupSync) SyncFromList(context context.ExecutionContext) {
	path := context.Options.GroupWhiteList
	whiteList, err := text.ReadLines(path)
	if err != nil {
		seelog.Errorf("Unable to load Google Group white list: file[%s]", path)
		explorer.FatalShutdown("Ensure file exist and readable: file[%s]", path)
	}
	for _, x := range whiteList {
		g.Sync(x)
	}
}
