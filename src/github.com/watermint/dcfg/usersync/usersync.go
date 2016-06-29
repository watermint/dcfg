package usersync

import (
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/connector"
	"github.com/watermint/dcfg/directory"
)

type UserSync struct {
	DropboxConnector connector.DropboxConnector
	DropboxAccounts  directory.AccountDirectory
	GoogleAccounts   directory.AccountDirectory
	GoogleGroups     directory.AccountDirectory
}

func (d *UserSync) membersNotInDirectory(member []directory.Account, ad directory.AccountDirectory) (notInDir []directory.Account) {
	for _, x := range member {
		if !directory.ExistInDirectory(ad, x) {
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

	seelog.Tracef("Dropbox [%d] user(s) are not in Google", len(dropboxMembersNotInGoogle))
	for _, x := range dropboxMembersNotInGoogle {
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
