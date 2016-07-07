package usersync

import "github.com/cihub/seelog"


// Deprovision Dropbox account based on Google side status.
// If the account, which identified by email, is not exist on Google Apps,
// this function deletes Dropbox account.
func (d *UserSync) SyncDeprovision() {
	seelog.Trace("Account Sync: Deprovision")

	dropboxMembers := d.DropboxAccounts.Accounts()
	dropboxMembersNotInGoogle := d.membersNotInDirectory(dropboxMembers, d.GoogleAccounts)
	dropboxMembersNotInGroup := d.membersNotInGroup(dropboxMembersNotInGoogle, d.GoogleGroups)

	seelog.Tracef("Dropbox [%d] user(s)", len(dropboxMembers))
	seelog.Tracef("Dropbox [%d] user(s) are not in Google", len(dropboxMembersNotInGoogle))
	seelog.Tracef("Dropbox [%d] user(s) are not in Group", len(dropboxMembersNotInGroup))
	for _, x := range dropboxMembersNotInGroup {
		seelog.Tracef("Removing Dropbox User: Email[%s]", x.Email)
		d.DropboxConnector.MembersRemove(x.Email)
	}
}

