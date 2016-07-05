package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/auth"
	"github.com/watermint/dcfg/config"
	"github.com/watermint/dcfg/connector"
	"github.com/watermint/dcfg/directory"
	"github.com/watermint/dcfg/explorer"
	"github.com/watermint/dcfg/groupsync"
	"github.com/watermint/dcfg/usersync"
	"io"
	"log"
	"os"
	"strings"
	"bytes"
	"unicode/utf16"
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
	BOM_UTF8 = []byte{0xef, 0xbb, 0xbf}
	BOM_UTF16BE = []byte{0xfe, 0xff}
	BOM_UTF16LE = []byte{0xff, 0xfe}
	BOM_UTF32BE = []byte{0x00, 0x00, 0xfe, 0xff}
	BOM_UTF32LE = []byte{0xff, 0xfe, 0x00, 0x00}
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
	for _, x := range groupSyncGroupList(syncOptions.GroupProvisionWhiteList) {
		groupSync.Sync(x)
	}
}

func trimBom(seq []byte) string {
	if bytes.HasPrefix(seq, BOM_UTF8) {
		return string(bytes.TrimPrefix(seq, BOM_UTF8))
	}
	if bytes.HasPrefix(seq, BOM_UTF16BE) {
		seqWithoutBom := bytes.TrimPrefix(seq, BOM_UTF16BE)
		utf16Seq := make([]uint16, len(seqWithoutBom) / 2)
		for i := range utf16Seq {
			utf16Seq[i] = uint16(seqWithoutBom[2 * i + 1] << 8) | uint16(seqWithoutBom[2 * i])
		}
		return string(utf16.Decode(utf16Seq))
	}
	return string(seq)
}

func groupSyncGroupList(filePath string) (list []string) {
	f, err := os.Open(filePath)
	if err != nil {
		seelog.Errorf("Unable to load Group Sync white list file: file[%s] err[%s]", filePath, err)
		explorer.FatalShutdown("Ensure file [%s] exist and readable", filePath)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		lineRaw, _, err := r.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			seelog.Errorf("Unable to load Group Sync white list file. Error during loading file: file[%s] err[%s]", filePath, err)
			explorer.FatalShutdown("Ensure file [%s] is appropriate format and encoding")
		}
		line := strings.TrimSpace(trimBom(lineRaw))
		if line != "" {
			list = append(list, line)
		}
	}
	return
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
