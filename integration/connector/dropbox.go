package connector

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/team"
	"github.com/watermint/dcfg/cli/explorer"
	"github.com/watermint/dcfg/common/util"
	"github.com/watermint/dcfg/integration/context"
)

type DropboxConnector interface {
	GroupsCreate(groupName, groupExternalId string) string
	GroupsUpdate(groupId, newGroupName string)
	GroupsMembersAdd(groupId, accountEmail string)
	GroupsMembersRemove(groupId, accountEmail string)

	MembersRemove(email string)
	MembersAdd(email, givenName, surname string)
}

func CreateConnector(context context.ExecutionContext) DropboxConnector {
	if context.Options.DryRun {
		return &DropboxConnectorMock{}
	} else {
		return &DropboxConnectorImpl{
			ExecutionContext: context,
		}
	}
}

type DropboxConnectorMock struct {
	history []string
}

func (dpm *DropboxConnectorMock) ClearOperationHistory() {
	dpm.history = []string{}
}

func (dpm *DropboxConnectorMock) CreateOperationLog(operationName string, arguments ...string) string {
	return fmt.Sprintf("[%s] %v", operationName, arguments)
}

func (dpm *DropboxConnectorMock) AssertLogs(expected []string) (unexpected []string, missing []string, success bool) {
	for _, x := range expected {
		if !util.ContainsString(dpm.history, x) {
			missing = append(missing, x)
		}
	}
	for _, x := range dpm.history {
		if !util.ContainsString(expected, x) {
			unexpected = append(unexpected, x)
		}
	}
	return unexpected, missing, len(unexpected) == 0 && len(missing) == 0
}

func (dpm *DropboxConnectorMock) enqueueOperationLog(operationName string, arguments ...string) {
	dpm.history = append(dpm.history, dpm.CreateOperationLog(operationName, arguments...))
}

func (dpm *DropboxConnectorMock) GroupsCreate(groupName, groupExternalId string) string {
	dpm.enqueueOperationLog("GroupsCreate", groupName, groupExternalId)
	explorer.ReportSuccess("Dropbox Group should be created: GroupName[%s] ExternalId[%s]", groupName, groupExternalId)
	return fmt.Sprintf("mock-%s", groupExternalId)
}
func (dpm *DropboxConnectorMock) GroupsUpdate(groupId, newGroupName string) {
	dpm.enqueueOperationLog("GroupsUpdate", groupId, newGroupName)
	explorer.ReportSuccess("Dropbox Group should be updated: GroupId[%s] NewGroupName[%s]", groupId, newGroupName)
}
func (dpm *DropboxConnectorMock) GroupsMembersAdd(groupId, accountEmail string) {
	dpm.enqueueOperationLog("GroupsMembersAdd", groupId, accountEmail)
	explorer.ReportSuccess("Member should be added to Dropbox Group: GroupId[%s] Member[%s]", groupId, accountEmail)
}
func (dpm *DropboxConnectorMock) GroupsMembersRemove(groupId, accountEmail string) {
	dpm.enqueueOperationLog("GroupsMembersRemove", groupId, accountEmail)
	explorer.ReportSuccess("Member should be removed from Dropbox Group: GroupId[%s] Member[%s]", groupId, accountEmail)
}
func (dpm *DropboxConnectorMock) MembersRemove(email string) {
	dpm.enqueueOperationLog("MembersRemove", email)
	explorer.ReportSuccess("Member account should be removed from Dropbox: Member[%s]", email)
}
func (dpm *DropboxConnectorMock) MembersAdd(email, givenName, surname string) {
	dpm.enqueueOperationLog("MembersAdd", email, givenName, surname)
	explorer.ReportSuccess("Member account should be added to Dropbox: Email[%s] GivenName[%s] Surname[%s]", email, givenName, surname)
}

type DropboxConnectorImpl struct {
	ExecutionContext context.ExecutionContext
}

func (dps *DropboxConnectorImpl) createGroupSelector(groupId string) (sel *team.GroupSelector) {
	return &team.GroupSelector{
		Tagged:  dropbox.Tagged{Tag: "group_id"},
		GroupId: groupId,
	}
}

func (dps *DropboxConnectorImpl) createUserSelectArg(accountEmail string) *team.UserSelectorArg {
	return &team.UserSelectorArg{
		Tagged: dropbox.Tagged{Tag: "email"},
		Email:  accountEmail,
	}
}

func (dps *DropboxConnectorImpl) createMemberAccess(accountEmail string) *team.MemberAccess {
	return &team.MemberAccess{
		User: dps.createUserSelectArg(accountEmail),
		AccessType: &team.GroupAccessType{
			Tagged: dropbox.Tagged{Tag: "member"},
		},
	}
}

func (dps *DropboxConnectorImpl) GroupsCreate(groupName, groupExternalId string) string {
	client := dps.ExecutionContext.DropboxClient
	a := team.GroupCreateArg{
		GroupName:       groupName,
		GroupExternalId: groupExternalId,
	}
	g, err := client.GroupsCreate(&a)
	if err != nil {
		seelog.Warnf("Unable to create Dropbox Group: GroupName[%s] ExternalId[%s] Err[%s]", groupName, groupExternalId, err)
		explorer.ReportFailure("Unable to create Dropbox Group: GroupName[%s] ExternalId[%s]", groupName, groupExternalId)
		return ""
	} else {
		seelog.Tracef("Dropbox Group Created: GroupId[%s] GroupName[%s] ExternalId[%s]", g.GroupId, g.GroupName, g.GroupExternalId)
		explorer.ReportSuccess("Dropbox Group Created: GroupId[%s] GroupName[%s] ExternalId[%s]", g.GroupId, g.GroupName, g.GroupExternalId)
		return g.GroupId
	}
}

func (dps *DropboxConnectorImpl) GroupsUpdate(groupId, newGroupName string) {
	client := dps.ExecutionContext.DropboxClient

	a := &team.GroupUpdateArgs{
		Group:        dps.createGroupSelector(groupId),
		NewGroupName: newGroupName,
	}
	g, err := client.GroupsUpdate(a)
	if err != nil {
		seelog.Warnf("Unable to update Dropbox Group: GroupId[%s] NewGroupname[%s] Err[%s]", err)
		explorer.ReportFailure("Unable to update Dropbox Group: GroupId[%s] NewGroupName[%s]", groupId, newGroupName)
	} else {
		seelog.Tracef("Dropbox Group Update: GroupId[%s] GroupName[%s] ExternalId[%s]")
		explorer.ReportSuccess("Dropbox Group Updated: GroupId[%s] GroupName[%s] ExternalId[%s]", g.GroupId, g.GroupName, g.GroupExternalId)
	}
}

func (dps *DropboxConnectorImpl) GroupsMembersAdd(groupId, accountEmail string) {
	client := dps.ExecutionContext.DropboxClient

	m := []*team.MemberAccess{dps.createMemberAccess(accountEmail)}
	a := &team.GroupMembersAddArg{
		Group:   dps.createGroupSelector(groupId),
		Members: m,
	}
	r, err := client.GroupsMembersAdd(a)
	if err != nil {
		seelog.Warnf("Unable to add member to Dropbox Group: GroupId[%s] AccountEmail[%s] Err[%s]", groupId, accountEmail, err)
		explorer.ReportFailure("Unable to add member to Dropbox Group: GroupId[%s] Email[%s]", groupId, accountEmail)
	} else {
		seelog.Tracef("Dropbox Group: Member added (Queued): GroupId[%s] GroupName[%s] AccountEmail[%s] AsyncJobId[%s]", r.GroupInfo.GroupId, r.GroupInfo.GroupName, accountEmail, r.AsyncJobId)
		explorer.ReportSuccess("Dropbox Group: Member added: GroupId[%s] GroupName[%s] AccountEmail[%s]", r.GroupInfo.GroupId, r.GroupInfo.GroupName, accountEmail)
	}
}

func (dps *DropboxConnectorImpl) GroupsMembersRemove(groupId, accountEmail string) {
	client := dps.ExecutionContext.DropboxClient

	a := &team.GroupMembersRemoveArg{
		Group: dps.createGroupSelector(groupId),
		Users: []*team.UserSelectorArg{dps.createUserSelectArg(accountEmail)},
	}
	r, err := client.GroupsMembersRemove(a)
	if err != nil {
		seelog.Warnf("Unable to remove member form Dropbox Group: GroupId[%s] AccountEmail[%s] Err[%s]", groupId, accountEmail, err)
		explorer.ReportFailure("Unable to remove member from Dropbox Group: GroupId[%s] AccountEmail[%s]", groupId, accountEmail)
	} else {
		seelog.Tracef("Dropbox Group: Member removed (queued): GroupId[%s] GroupName[%s] AccountEmail[%s] AsyncJobId[%s]", r.GroupInfo, r.GroupInfo.GroupName, accountEmail, r.AsyncJobId)
		explorer.ReportSuccess("Dropbox Group: Member removed: GroupId[%s] GroupName[%s] AccountEmail[%s]", r.GroupInfo.GroupId, r.GroupInfo.GroupName, accountEmail)
	}
}

func (dps *DropboxConnectorImpl) MembersRemove(email string) {
	client := dps.ExecutionContext.DropboxClient

	m := team.MembersGetInfoArgs{
		Members: []*team.UserSelectorArg{dps.createUserSelectArg(email)},
	}
	u, err := client.MembersGetInfo(&m)
	if err != nil {
		seelog.Warnf("Unable to load Dropbox member: Email[%s] Err[%s]", email, err)
		explorer.ReportFailure("Unable to remove member Dropbox account: Email[%s] (due to failed to load member info)", email)
		return
	}
	if len(u) != 1 {
		seelog.Warnf("Unable to load Dropbox member: Email[%s] [%v]", email, u)
		explorer.ReportFailure("Unable to remove member Dropbox account: Email[%s] (due to failed to load member info)", email)
		return
	} else {
		if u[0].MemberInfo.Role.Tag == "team_admin" {
			seelog.Warnf("Team Admin should not be removed by script: Email[%s]", email)
			explorer.ReportFailure("Unable to remove Dropbox Team Admin account: Email[%s]", email)
			return
		}
	}

	a := team.MembersRemoveArg{
		MembersDeactivateArg: team.MembersDeactivateArg{
			User:     dps.createUserSelectArg(email),
			WipeData: false,
		},
		KeepAccount: false,
	}
	r, err := client.MembersRemove(&a)
	if err != nil {
		seelog.Warnf("Unable to remove member Dropbox account: Email[%s] Err[%s]", email, err)
		explorer.ReportFailure("Unable to remove member Dropbox account: Email[%s]", email)
	} else {
		seelog.Tracef("Remove Dropbox account: Email[%s] Tag[%s]", email, r.Tag)
		explorer.ReportSuccess("Remove Dropbox account: Email[%s]", email)
	}
}

func (dps *DropboxConnectorImpl) MembersAdd(email, givenName, surname string) {
	client := dps.ExecutionContext.DropboxClient

	a := team.MembersAddArg{
		NewMembers: []*team.MemberAddArg{
			&team.MemberAddArg{
				MemberEmail:     email,
				MemberGivenName: givenName,
				MemberSurname:   surname,
				Role:            &team.AdminTier{Tagged: dropbox.Tagged{Tag: "member_only"}},
			},
		},
	}
	r, err := client.MembersAdd(&a)
	if err != nil {
		seelog.Warnf("Unable to add member Dropbox account: Email[%s] GivenName[%s] Surname[%s] Err[%s]", email, givenName, surname, err)
		explorer.ReportFailure("Unable to add member Dropbox account: Email[%s] GivenName[%s] Surname[%s]", email, givenName, surname)
	} else {
		seelog.Tracef("Add Dropbox account: Email[%s] GivenName[%s] Surname[%s] Tag[%s]", email, givenName, surname, r.Tag)
		explorer.ReportSuccess("Add Dropbox account: Email[%s] GivenName[%s] Surname[%s]", email, givenName, surname)
	}
}
