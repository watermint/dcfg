package auth

import (
	"github.com/watermint/dcfg/explorer"
	"github.com/cihub/seelog"
)

// Verify Dropbox connection. Exit the program if connection failure.
func VerifyDropbox() {
	teamClient := DropboxClient()

	info, err := teamClient.GetInfo()
	if err != nil {
		seelog.Errorf("Failed to load Dropbox Team info: err[%v]", err)
		explorer.FatalShutdown("Please re-run `-auth dropbox`")
	}

	seelog.Infof("Dropbox Team Name: %s", info.Name)
	seelog.Infof("Dropbox Team ID: %s", info.TeamId)
	seelog.Infof("Dropbox Provisioned Number: %d", info.NumProvisionedUsers)
}

// Verify Google connection. Exit the program if connection failure.
func VerifyGoogle(domain string) {
	client := GoogleClient()

	groups, err := client.Groups.List().Domain(domain).Do()
	if err != nil {
		seelog.Errorf("Failed to load Google Group: domain[%s] err[%v]", domain, err)
		explorer.FatalShutdown("Please re-run `-auth google`")
	}
	users, err := client.Users.List().Domain(domain).Do()
	if err != nil {
		seelog.Errorf("Failed to load Google Users: err[%s]", err)
		explorer.FatalShutdown("Please re-run `-auth google`")
	}

	seelog.Infof("Google Groups: %d (partial)", len(groups.Groups))
	seelog.Infof("Google Members: %d (partial)", len(users.Users))
}
