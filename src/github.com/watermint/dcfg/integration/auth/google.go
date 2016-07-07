package auth

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/cli/explorer"
	"github.com/watermint/dcfg/common/file"
	"github.com/watermint/dcfg/integration/context"
	"golang.org/x/oauth2"
)

const (
	GOOGLE_CUSTOMER_ID = "my_customer"
)

func getGoogleTokenFromWeb(context context.ExecutionContext) *oauth2.Token {
	seelog.Flush()

	config := context.GoogleClientConfig
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Go to the following link in your browser then type the authorization code:")
	fmt.Println("")
	fmt.Println(authURL)
	fmt.Println("")
	fmt.Println("------")
	fmt.Println("Paste code here:")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		seelog.Errorf("Unable to read authorization code %v", err)
		explorer.FatalShutdown("Please re-run, then enter new authorisation code")
	}

	fmt.Println("")

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		seelog.Errorf("Unable to retrieve token from web %v", err)
		explorer.FatalShutdown("Please re-run, then enter new authorisation code")
	}
	return tok
}

func verifyGoogleToken(context context.ExecutionContext, token *oauth2.Token) {
	client, err := context.CreateGoogleClientByToken(token)
	if err != nil {
		seelog.Errorf("Authentication failed. err[%s]", err)
		explorer.FatalShutdown("Please re-run `-auth google` sequence")
	}
	_, err = client.Groups.List().Customer(GOOGLE_CUSTOMER_ID).Do()
	if err != nil {
		seelog.Errorf("Authentication failed. err[%s]", err)
		explorer.FatalShutdown("Please re-run `-auth google` sequence")
	}
	_, err = client.Users.List().Customer(GOOGLE_CUSTOMER_ID).Do()
	if err != nil {
		seelog.Errorf("Authentication failed. err[%s]", err)
		explorer.FatalShutdown("Please re-run `-auth google` sequence")
	}
	explorer.ReportSuccess("Verified token for Google Apps")
}

func UpdateGoogleToken(context context.ExecutionContext) {
	path := context.Options.PathGoogleToken()
	token := getGoogleTokenFromWeb(context)
	verifyGoogleToken(context, token)

	if err := file.SaveJSON(path, token); err != nil {
		seelog.Errorf("Unable to write Google token file: file[%s] err[%s]", path, err)
		explorer.FatalShutdown("Cannot update Google token file: file[%s]", path)
	}
}

func AuthGoogle(context context.ExecutionContext) {
	seelog.Info("Start authentication sequence for Google Apps")
	UpdateGoogleToken(context)
}
