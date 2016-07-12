package usersync

import (
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/cli/explorer"
	"github.com/watermint/dcfg/integration/directory"
)

// Deprovision Dropbox account based on Google side status.
// If the account, which identified by email, is not exist on Google Apps,
// this function deletes Dropbox account.
func (d *UserSync) SyncDeprovision() {
	seelog.Trace("Account Sync: Deprovision")

	dropboxMembers := d.DropboxAccounts.Accounts()
	dropboxMembersNotInGoogle := make([]directory.Account, 0)

	for _, x := range dropboxMembers {
		exist, err := d.GoogleEmail.EmailExist(x.Email)
		if err != nil {
			seelog.Errorf("Cannot load emails of Google")
			explorer.FatalShutdown("Please verify Google auth status or network condition")
		}
		if !exist {
			dropboxMembersNotInGoogle = append(dropboxMembersNotInGoogle, x)
		}
	}

	seelog.Tracef("Dropbox [%d] user(s)", len(dropboxMembers))
	seelog.Tracef("Dropbox [%d] user(s) are not in Google", len(dropboxMembersNotInGoogle))
	for _, x := range dropboxMembersNotInGoogle {
		seelog.Tracef("Removing Dropbox User: Email[%s]", x.Email)
		d.DropboxConnector.MembersRemove(x.Email)
	}
}
