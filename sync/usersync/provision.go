package usersync

import "github.com/cihub/seelog"

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

