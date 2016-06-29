package groupsync

import (
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/connector"
	"github.com/watermint/dcfg/directory"
	"github.com/watermint/dcfg/explorer"
)

type GroupSync struct {
	DropboxConnector        connector.DropboxConnector
	DropboxAccountDirectory directory.AccountDirectory
	DropboxGroupDirectory   directory.GroupDirectory
	GoogleDirectory         directory.GroupDirectory
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
