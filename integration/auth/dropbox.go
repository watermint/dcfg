package auth

import (
	"bytes"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/cli/explorer"
	"github.com/watermint/dcfg/common/file"
	"github.com/watermint/dcfg/integration/context"
	"log"
	"strings"
)

func verifyDropboxToken(context context.ExecutionContext, token string) {
	verboseOutput := bytes.NewBufferString("Verbose")
	log.SetOutput(verboseOutput)
	client := context.CreateDropboxClientByToken(token)
	team, err := client.GetInfo()

	verboseLog := verboseOutput.String()
	verboseLines := strings.Split(verboseLog, "\n")
	for i, l := range verboseLines {
		seelog.Tracef("Verbose logs[%d]: %s", i, l)
	}

	if err != nil {
		seelog.Errorf("Authentication failed [%s]", err)
		explorer.FatalShutdown("Please regenerate Dropbox Business API token, then update token file")
	}
	explorer.ReportSuccess("Verified token for Dropbox Team: TeamId[%s] TeamName[%s] Provisioned[%d] Num Licenses[%d]", team.TeamId, team.Name, team.NumProvisionedUsers, team.NumLicensedUsers)
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

func updateDropboxToken(context context.ExecutionContext) {
	token := getDropboxTokenFromConsole()
	path := context.Options.PathDropboxToken()

	verifyDropboxToken(context, token)
	err := file.SaveJSON(path, token)
	if err != nil {
		seelog.Errorf("Unable to write Dropbox token file", path, err)
		explorer.FatalShutdown("Ensure file [%s] is appropriate JSON format.", path)
	}
	explorer.ReportSuccess("Dropbox Token file updated: [%s]", path)
}

func AuthDropbox(context context.ExecutionContext) {
	seelog.Info("Start authentication sequence for Dropbox")
	updateDropboxToken(context)
}
