package auth

import (
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/explorer"
	"net/http"
	"net/http/httputil"
	"strings"
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

func verifyNetwork(host string) {
	seelog.Tracef("Verifying network reachability to [%s]", host)
	resp, err := http.Head(host)

	if resp == nil {
		seelog.Tracef("Response: nil")
	} else {
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			seelog.Tracef("Status code: %d", resp.StatusCode)
			seelog.Tracef("Status: %s", resp.Status)
			seelog.Tracef("TLS Cipher Suite: %d", resp.TLS.CipherSuite)
			seelog.Tracef("TLS Server name: %s", resp.TLS.ServerName)
			seelog.Tracef("TLS Version: %d", resp.TLS.Version)
			seelog.Tracef("Proto: %s", resp.Proto)
			seelog.Tracef("Content length: %d", resp.ContentLength)
		} else {
			respLines := string(respDump)
			for i, l := range strings.Split(respLines, "\n") {
				seelog.Tracef("Response[%d]: %s", i, l)
			}
		}
	}
	if err != nil {
		seelog.Errorf("Error during verify connecto to host [%s]: %v", host, err)
		explorer.FatalShutdown("Please check you network configuration to host [%s]", host)
	}
	seelog.Tracef("Verification finished without error: host[%s]", host)
	seelog.Infof("Network test: Success: %s", host)
}
