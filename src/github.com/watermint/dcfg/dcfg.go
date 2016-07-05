package main

import (
	"flag"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/auth"
	"github.com/watermint/dcfg/config"
	"github.com/watermint/dcfg/connector"
	"github.com/watermint/dcfg/directory"
	"github.com/watermint/dcfg/explorer"
	"github.com/watermint/dcfg/groupsync"
	"github.com/watermint/dcfg/text"
	"github.com/watermint/dcfg/usersync"
	"log"
	"os"
	"strings"
)

const (
	seeLogXmlTemplate = `
	<seelog type="adaptive" mininterval="200000000" maxinterval="1000000000" critmsgcount="5">
	<formats>
    		<format id="detail" format="date:%%Date(2006-01-02T15:04:05Z07:00)%%tloc:%%File:%%FuncShort:%%Line%%tlevel:%%Level%%tmsg:%%Msg%%n" />
    		<format id="short" format="%%Time [%%LEV] %%Msg%%n" />
	</formats>
	<outputs formatid="detail">
    		<filter levels="trace,info,warn,error,critical">
        		<rollingfile formatid="detail" filename="%s/dcfg.log" type="size" maxsize="52428800" maxrolls="7" />
    		</filter>
		<filter levels="info,warn,error,critical">
        		<console formatid="short" />
    		</filter>
    	</outputs>
	</seelog>
	`
)

var (
	AppVersion string
)

func replaceLogger(basePath string) {
	seeLogXml := fmt.Sprintf(seeLogXmlTemplate, basePath)
	logger, err := seelog.LoggerFromConfigAsString(seeLogXml)
	if err != nil {
		log.Fatalln("Failed to load logger", err.Error())
	}
	seelog.ReplaceLogger(logger)
}

func DirExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func syncUserProvision(googleDirectory *directory.GoogleDirectory, dropboxDirectory *directory.DropboxDirectory, provisioning *connector.DropboxConnector, syncOptions SyncOptions) {
	userSync := usersync.UserSync{
		DropboxConnector: *provisioning,
		DropboxAccounts:  dropboxDirectory,
		GoogleAccounts:   googleDirectory,
		GoogleGroups:     googleDirectory,
	}

	seelog.Infof("Provisioning Users (Google Users -> Dropbox Users)")
	userSync.SyncProvision()
}
func syncUserDeprovision(googleDirectory *directory.GoogleDirectory, dropboxDirectory *directory.DropboxDirectory, provisioning *connector.DropboxConnector, syncOptions SyncOptions) {
	userSync := usersync.UserSync{
		DropboxConnector: *provisioning,
		DropboxAccounts:  dropboxDirectory,
		GoogleAccounts:   googleDirectory,
		GoogleGroups:     googleDirectory,
	}

	seelog.Infof("Deprovisioning Users (Google Users -> Dropbox Users)")
	userSync.SyncDeprovision()
}
func syncGroupProvision(googleDirectory *directory.GoogleDirectory, dropboxDirectory *directory.DropboxDirectory, provisioning *connector.DropboxConnector, syncOptions SyncOptions) {
	if syncOptions.GroupProvisionWhiteList == "" {
		seelog.Errorf("Group white list file required for group provisioning")
		explorer.FatalShutdown("Please specify with option -group-provision-list")
	}

	groupSync := groupsync.GroupSync{
		DropboxConnector:        *provisioning,
		DropboxAccountDirectory: dropboxDirectory,
		DropboxGroupDirectory:   dropboxDirectory,
		GoogleDirectory:         googleDirectory,
	}

	seelog.Infof("Syncing Group (Google Group -> Dropbox Group)")
	whiteList, err := text.ReadLines(syncOptions.GroupProvisionWhiteList)
	if err != nil {
		seelog.Errorf("Unable to load Google Group white list: file[%s]", syncOptions.GroupProvisionWhiteList)
		explorer.FatalShutdown("Ensure file exist and readable: file[%s]", syncOptions.GroupProvisionWhiteList)
	}
	for _, x := range whiteList {
		groupSync.Sync(x)
	}
}

func loadDirectories(syncOptions SyncOptions) (*directory.GoogleDirectory, *directory.DropboxDirectory) {
	dropbox := directory.DropboxDirectory{}
	google := directory.GoogleDirectory{}
	dropbox.Load()
	google.Load()

	return &google, &dropbox
}

func sync(target string, connector *connector.DropboxConnector, syncOptions SyncOptions) {
	for _, x := range strings.Split(target, ",") {
		switch x {
		case "group-provision":
			googleDirectory, dropboxDirectory := loadDirectories(syncOptions)
			syncGroupProvision(googleDirectory, dropboxDirectory, connector, syncOptions)
		case "user-provision":
			googleDirectory, dropboxDirectory := loadDirectories(syncOptions)
			syncUserProvision(googleDirectory, dropboxDirectory, connector, syncOptions)
		case "user-deprovision":
			googleDirectory, dropboxDirectory := loadDirectories(syncOptions)
			syncUserDeprovision(googleDirectory, dropboxDirectory, connector, syncOptions)
		default:
			seelog.Errorf("Undefined sync type [%s]", x)
		}
	}
}

func usage() {
	seelog.Flush()
	flag.Usage()
	seelog.Flush()
}

func executeAuth(authMode string) {
	switch authMode {
	case "google":
		defer explorer.Report()
		seelog.Trace("Google: Authorisation sequence")
		auth.AuthGoogle()

	case "dropbox":
		defer explorer.Report()
		seelog.Trace("Dropbox: Authorisation sequence")
		auth.AuthDropbox()

	default:
		fmt.Errorf("Undefined auth mode: [%s]", authMode)
		usage()
	}
}

type SyncOptions struct {
	GroupProvisionWhiteList string
	DryRun                  bool
}

func createDropboxConnector(dryRun bool) connector.DropboxConnector {
	if dryRun {
		seelog.Info("Run sync as dry run")
		return &connector.DropboxConnectorMock{}
	} else {
		return &connector.DropboxConnectorImpl{}
	}
}

func main() {
	authTarget := flag.String("auth", "", "Store API token. Choose target API provider (google/dropbox)")
	syncTarget := flag.String("sync", "", "Sync mode. Separate by comma if want to use multiple mode. (user-provision/user-deprovision/group-provision)")
	basePath := flag.String("path", "", "Path for config/log files.")
	dryRun := flag.Bool("dryrun", true, "Dry run")
	proxy := flag.String("proxy", "", "Proxy hostname:port")
	groupProvisionWhiteList := flag.String("group-provision-list", "", "White list file for group-provision")

	flag.Parse()

	if *basePath == "" {
		fmt.Printf("Path required\n")
		usage()
		os.Exit(1)
	}
	if !DirExists(*basePath) {
		fmt.Printf("Directory not exist: %s\n", *basePath)
		usage()
		os.Exit(1)
	}

	replaceLogger(*basePath)

	seelog.Info("dcfg version: ", AppVersion)

	if *proxy != "" {
		seelog.Tracef("Explicit proxy configuration: HTTP_PROXY [%s]", *proxy)
		seelog.Tracef("Explicit proxy configuration: HTTPS_PROXY [%s]", *proxy)
		os.Setenv("HTTP_PROXY", *proxy)
		os.Setenv("HTTPS_PROXY", *proxy)
	}

	explorer.Start()
	config.ReloadConfig(*basePath)
	defer seelog.Flush()

	syncOptions := SyncOptions{
		GroupProvisionWhiteList: *groupProvisionWhiteList,
		DryRun:                  *dryRun,
	}

	if *authTarget != "" {
		executeAuth(*authTarget)
	} else if *syncTarget != "" {

		defer explorer.Report()

		connector := createDropboxConnector(syncOptions.DryRun)

		sync(*syncTarget, &connector, syncOptions)

	} else {
		usage()
	}
}
