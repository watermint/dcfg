package auth

import (
	"io/ioutil"
	"github.com/watermint/dcfg/config"
	"github.com/watermint/dcfg/explorer"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/admin/directory/v1"
	"golang.org/x/oauth2"
	"golang.org/x/net/context"
	"encoding/json"
	"os"
	"fmt"
	"github.com/cihub/seelog"
)

func googleConfig() *oauth2.Config {
	json, err := ioutil.ReadFile(config.Global.GoogleClientFile())
	if err != nil {
		explorer.Fatal("Unable to read Google client file", config.Global.GoogleClientFile(), err)
	}

	c, err := google.ConfigFromJSON(json,
		admin.AdminDirectoryUserReadonlyScope,
		admin.AdminDirectoryGroupReadonlyScope)

	if err != nil {
		explorer.Fatal("Unable to parse Google client file", config.Global.GoogleClientFile(), err)
	}

	return c
}

func googleClientByToken(token *oauth2.Token) *admin.Service {
	cfg := googleConfig()
	context := context.Background()
	client := cfg.Client(context, token)
	service, err := admin.New(client)
	if err != nil {
		explorer.Fatal("Unable to create Google client", err)
	}

	return service
}

func GoogleClient() *admin.Service {
	j, err := os.Open(config.Global.GoogleTokenFile())
	if err != nil {
		explorer.Fatal("Unable to read Google token file", config.Global.GoogleTokenFile(), err)
	}
	defer j.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(j).Decode(token)
	if err != nil {
		explorer.Fatal("Unable to parse Google token file", config.Global.GoogleTokenFile(), err)
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
		explorer.Fatal("Unable to read authorization code %v", err)
	}

	fmt.Println("")

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		explorer.Fatal("Unable to retrieve token from web %v", err)
	}
	return tok
}

func verifyGoogleToken(token *oauth2.Token, domain string) {
	client := googleClientByToken(token)
	_, err := client.Groups.List().Domain(domain).Do()
	if err != nil {
		explorer.Fatal("Authentication failed. domain[%s] err[%s]", domain, err)
	}
	_, err = client.Users.List().Domain(domain).Do()
	if err != nil {
		explorer.Fatal("Authentication failed. domain[%s] err[%s]", domain, err)
	}
	explorer.ReportSuccess("Verified token for Google Apps")
}

func UpdateGoogleToken(domain string) {
	token := getGoogleTokenFromWeb()

	verifyGoogleToken(token, domain)

	j, err := os.Create(config.Global.GoogleTokenFile())
	if err != nil {
		explorer.Fatal("Unable to open Google token file", config.Global.GoogleTokenFile(), err)
	}
	defer j.Close()

	err = json.NewEncoder(j).Encode(token)
	if err != nil {
		explorer.Fatal("Unable to write Google token file", config.Global.GoogleTokenFile(), err)
	}
	explorer.ReportSuccess("Google Token file updated: [%s]", config.Global.GoogleTokenFile())
}

func AuthGoogle(domain string) {
	seelog.Info("Start authentication sequence for Google Apps")
	UpdateGoogleToken(domain)
}
