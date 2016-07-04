package auth

import (
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/dropbox/dropbox-sdk-go-unofficial"
	"github.com/watermint/dcfg/config"
	"github.com/watermint/dcfg/explorer"
	"os"
)

func verifyDropboxToken(token string) {
	client := dropboxClientFromToken(token)
	team, err := client.GetInfo()
	if err != nil {
		seelog.Errorf("Authentication failed [%s]", err)
		explorer.FatalShutdown("Please regenerate Dropbox Business API token, then update token file [%s]", config.Global.DropboxTokenFile())
	}
	explorer.ReportSuccess("Verified token for Dropbox Team: TeamId[%s] TeamName[%s] Provisioned[%d] Num Licenses[%d]", team.TeamId, team.Name, team.NumProvisionedUsers, team.NumLicensedUsers)
}

func dropboxClientFromToken(token string) dropbox.Api {
	return dropbox.Client(token, dropbox.Options{})
}

func getDropboxTokenFromConsole() string {
	seelog.Flush()

	fmt.Println("Dropbox Business API (permisson type: Team member management)")
	fmt.Println("")
	fmt.Println("------")
	fmt.Println("Paste generated code here:")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		seelog.Errorf("Unable to read authorization code %v", err)
		explorer.FatalShutdown("Please re-run auth command. Then, paste generated code")
	}

	fmt.Println("")

	return code
}

func updateDropboxToken() {
	token := getDropboxTokenFromConsole()

	verifyDropboxToken(token)

	content := config.DropboxToken{
		TeamManagementToken: token,
	}

	j, err := os.Create(config.Global.DropboxTokenFile())
	if err != nil {
		seelog.Errorf("Unable to open Dropbox token file: file[%s] err[%s]", config.Global.DropboxTokenFile(), err)
		explorer.FatalShutdown("Ensure file [%s] exist and readable", config.Global.DropboxTokenFile())
	}
	defer j.Close()

	err = json.NewEncoder(j).Encode(content)
	if err != nil {
		seelog.Errorf("Unable to write Dropbox token file", config.Global.DropboxTokenFile(), err)
		explorer.FatalShutdown("Ensure file [%s] is appropriate JSON format.", config.Global.DropboxTokenFile())
	}
	explorer.ReportSuccess("Dropbox Token file updated: [%s]", config.Global.DropboxTokenFile())

}

func DropboxClient() dropbox.Api {
	return dropboxClientFromToken(config.Global.DropboxToken().TeamManagementToken)
}

func AuthDropbox() {
	seelog.Info("Start authentication sequence for Dropbox")
	updateDropboxToken()
}
