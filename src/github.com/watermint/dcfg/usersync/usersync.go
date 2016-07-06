package usersync

import (
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/connector"
	"github.com/watermint/dcfg/context"
	"github.com/watermint/dcfg/directory"
)

type UserSync struct {
	DropboxConnector connector.DropboxConnector
	DropboxAccounts  directory.AccountDirectory
	GoogleAccounts   directory.AccountDirectory
	GoogleGroups     directory.GroupResolver
}

func NewUserSync(context context.ExecutionContext) UserSync {
	gd := directory.GoogleDirectory{ExecutionContext: context}
	dd := directory.DropboxDirectory{ExecutionContext: context}
	dp := connector.CreateConnector(context)
	return UserSync{
		DropboxConnector: dp,
		DropboxAccounts:  &dd,
		GoogleAccounts:   &gd,
		GoogleGroups:     &gd,
	}
}

func (d *UserSync) membersNotInDirectory(member []directory.Account, ad directory.AccountDirectory) (notInDir []directory.Account) {
	for _, x := range member {
		if !directory.ExistInDirectory(ad, x) {
			notInDir = append(notInDir, x)
		}
	}
	return
}

func (d *UserSync) memberInGroups(needle directory.Account, haystack []directory.Group) bool {
	for _, x := range haystack {
		if needle.Email == x.GroupEmail {
			return true
		}
	}
	return false
}

func (d *UserSync) membersNotInGroup(member []directory.Account, gd directory.GroupResolver) (notInDir []directory.Account) {
	for _, x := range member {
		_, exist := gd.Group(x.Email)
		if !exist {
			notInDir = append(notInDir, x)
		}
	}
	return
}

// Deprovision Dropbox account based on Google side status.
// If the account, which identified by email, is not exist on Google Apps,
// this function deletes Dropbox account.
func (d *UserSync) SyncDeprovision() {
	seelog.Trace("Account Sync: Deprovision")

	dropboxMembers := d.DropboxAccounts.Accounts()
	dropboxMembersNotInGoogle := d.membersNotInDirectory(dropboxMembers, d.GoogleAccounts)
	dropboxMembersNotInGroup := d.membersNotInGroup(dropboxMembersNotInGoogle, d.GoogleGroups)

	seelog.Tracef("Dropbox [%d] user(s) are not in Google", len(dropboxMembersNotInGroup))
	for _, x := range dropboxMembersNotInGroup {
		seelog.Tracef("Removing Dropbox User: Email[%s]", x.Email)
		d.DropboxConnector.MembersRemove(x.Email)
	}
}

func (d *UserSync) SyncProvision() {
	seelog.Trace("Account Sync: Provision")

	googleMembers := d.GoogleAccounts.Accounts()
	googleMembersNotInDropbox := d.membersNotInDirectory(googleMembers, d.DropboxAccounts)

	seelog.Tracef("%d users in Google Apps", len(googleMembers))
	seelog.Tracef("Google [%d] user(s) are not in Dropbox", len(googleMembersNotInDropbox))
	for _, x := range googleMembersNotInDropbox {
		seelog.Tracef("Adding Dropbox User: Email[%s]", x)
		d.DropboxConnector.MembersAdd(x.Email, x.GivenName, x.Surname)
	}
}
