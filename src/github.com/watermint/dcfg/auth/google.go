package auth

import (
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/config"
	"github.com/watermint/dcfg/explorer"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/admin/directory/v1"
	"io/ioutil"
	"os"
)

func googleConfig() *oauth2.Config {
	json, err := ioutil.ReadFile(config.Global.GoogleClientFile())
	if err != nil {
		seelog.Errorf("Unable to read Google client file: file[%s] err[%s]", config.Global.GoogleClientFile(), err)
		explorer.FatalShutdown("Ensure file [%s] exist and readable", config.Global.GoogleClientFile())
	}

	c, err := google.ConfigFromJSON(json,
		admin.AdminDirectoryUserReadonlyScope,
		admin.AdminDirectoryGroupReadonlyScope)

	if err != nil {
		seelog.Errorf("Unable to parse Google client file: file[%s] err[%s]", config.Global.GoogleClientFile(), err)
		explorer.FatalShutdown("Ensure file [%s] is appropriate JSON format.", config.Global.GoogleClientFile())
	}

	return c
}

func googleClientByToken(token *oauth2.Token) *admin.Service {
	cfg := googleConfig()
	context := context.Background()
	client := cfg.Client(context, token)
	service, err := admin.New(client)
	if err != nil {
		seelog.Errorf("Unable to create Google client: err[%s]", err)
		explorer.FatalShutdown("Ensure Google API token available to use.")
	}

	return service
}

func GoogleClient() *admin.Service {
	j, err := os.Open(config.Global.GoogleTokenFile())
	if err != nil {
		seelog.Errorf("Unable to read Google token file", config.Global.GoogleTokenFile(), err)
		explorer.FatalShutdown("Ensure file [%s] exist and readable", config.Global.GoogleTokenFile())
	}
	defer j.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(j).Decode(token)
	if err != nil {
		seelog.Errorf("Unable to parse Google token file", config.Global.GoogleTokenFile(), err)
		explorer.FatalShutdown("Ensure file [%s] is appropriate JSON format", config.Global.GoogleTokenFile())
	}
	return googleClientByToken(token)
}

func getGoogleTokenFromWeb() *oauth2.Token {
	seelog.Flush()

	config := googleConfig()
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

func verifyGoogleToken(token *oauth2.Token, domain string) {
	client := googleClientByToken(token)
	_, err := client.Groups.List().Domain(domain).Do()
	if err != nil {
		seelog.Errorf("Authentication failed. domain[%s] err[%s]", domain, err)
		explorer.FatalShutdown("Please re-run `-auth google` sequence")
	}
	_, err = client.Users.List().Domain(domain).Do()
	if err != nil {
		seelog.Errorf("Authentication failed. domain[%s] err[%s]", domain, err)
		explorer.FatalShutdown("Please re-run `-auth google` sequence")
	}
	explorer.ReportSuccess("Verified token for Google Apps")
}

func UpdateGoogleToken(domain string) {
	token := getGoogleTokenFromWeb()

	verifyGoogleToken(token, domain)

	j, err := os.Create(config.Global.GoogleTokenFile())
	if err != nil {
		seelog.Errorf("Unable to open Google token file: file[%s] err[%s]", config.Global.GoogleTokenFile(), err)
		explorer.FatalShutdown("Ensure file [%s] exist and readable", config.Global.GoogleTokenFile())
	}
	defer j.Close()

	err = json.NewEncoder(j).Encode(token)
	if err != nil {
		seelog.Errorf("Unable to write Google token file: file[%s] err[%s]", config.Global.GoogleTokenFile(), err)
		explorer.FatalShutdown("Ensure file [%s] is appropriate JSON format", config.Global.GoogleTokenFile())
	}
	explorer.ReportSuccess("Google Token file updated: [%s]", config.Global.GoogleTokenFile())
}

func AuthGoogle(domain string) {
	seelog.Info("Start authentication sequence for Google Apps")
	UpdateGoogleToken(domain)
}
