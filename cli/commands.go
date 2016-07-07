package cli

import (
	"errors"
	"flag"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/common/domain"
	"github.com/watermint/dcfg/common/file"
	"os"
	"path"
	"strings"
)

type Options struct {
	ModeAuth       string
	ModeSync       string
	BasePath       string
	DryRun         bool
	Proxy          string
	GroupWhiteList string
}

const (
	MODE_AUTH_DROPBOX = "dropbox"
	MODE_AUTH_GOOGLE  = "google"

	MODE_SYNC_GROUP_PROVISION  = "group-provision"
	MODE_SYNC_USER_PROVISION   = "user-provision"
	MODE_SYNC_USER_DEPROVISION = "user-deprovision"

	optNameModeAuth       = "auth"
	optNameModeSync       = "sync"
	optNameBasePath       = "path"
	optNameDryRun         = "dryrun"
	optNameProxy          = "proxy"
	optNameGroupWhiteList = "group-provision-list"

	FILENAME_GOOGLE_TOKEN         = "google_token.json"
	FILENAME_GOOGLE_CLIENT_SECRET = "google_client_secret.json"
	FILENAME_DROPBOX_TOKEN        = "dropbox_token.json"
)

var (
	modeAuthOpts = []string{MODE_AUTH_GOOGLE, MODE_AUTH_DROPBOX}
	modeSyncOpts = []string{MODE_SYNC_USER_PROVISION, MODE_SYNC_USER_DEPROVISION, MODE_SYNC_GROUP_PROVISION}

	optDescModeAuth       = fmt.Sprintf("Update API token. Choose API provider (%s)", strings.Join(modeAuthOpts, ", "))
	optDescModeSync       = fmt.Sprintf("Sync mode. Separate by comma if you want ot use multiple modes. (%s)", strings.Join(modeSyncOpts, ", "))
	optDescBasePath       = "Path for config/log files."
	optDescDryRun         = "Dry run"
	optDescProxy          = "HTTP(S) proxy (hostname:port)"
	optDescGroupWhiteList = "White list file for group-provision"
)

func (o *Options) IsModeAuth() bool {
	return o.ModeAuth != ""
}
func (o *Options) IsModeSync() bool {
	return o.ModeSync != ""
}

func (o *Options) IsModeAuthGoogle() bool {
	return o.ModeAuth == MODE_AUTH_GOOGLE
}
func (o *Options) IsModeAuthDropbox() bool {
	return o.ModeAuth == MODE_AUTH_DROPBOX
}
func (o *Options) IsModeSyncUserProvision() bool {
	modes := strings.Split(o.ModeSync, ",")
	return domain.ContainsString(modes, MODE_SYNC_USER_PROVISION)
}
func (o *Options) IsModeSyncUserDeprovision() bool {
	modes := strings.Split(o.ModeSync, ",")
	return domain.ContainsString(modes, MODE_SYNC_USER_DEPROVISION)
}
func (o *Options) IsModeGroupProvision() bool {
	modes := strings.Split(o.ModeSync, ",")
	return domain.ContainsString(modes, MODE_SYNC_GROUP_PROVISION)
}
func (o *Options) PathGoogleToken() string {
	return path.Join(o.BasePath, FILENAME_GOOGLE_TOKEN)
}
func (o *Options) PathGoogleClientSecret() string {
	return path.Join(o.BasePath, FILENAME_GOOGLE_CLIENT_SECRET)
}
func (o *Options) PathDropboxToken() string {
	return path.Join(o.BasePath, FILENAME_DROPBOX_TOKEN)
}

func (o *Options) Parse() error {
	modeAuth := flag.String(optNameModeAuth, "", optDescModeAuth)
	modeSync := flag.String(optNameModeSync, "", optDescModeSync)
	basePath := flag.String(optNameBasePath, "", optDescBasePath)
	proxy := flag.String(optNameProxy, "", optDescProxy)
	dryRun := flag.Bool(optNameDryRun, true, optDescDryRun)
	groupWhiteList := flag.String(optNameGroupWhiteList, "", optDescGroupWhiteList)

	flag.Parse()

	o.ModeAuth = *modeAuth
	o.ModeSync = *modeSync
	o.BasePath = *basePath
	o.Proxy = *proxy
	o.DryRun = *dryRun
	o.GroupWhiteList = *groupWhiteList

	return nil
}

func (o *Options) UpdateEnv() {
	if o.Proxy != "" {
		seelog.Tracef("Explicit proxy configuration: HTTP_PROXY [%s]", o.Proxy)
		seelog.Tracef("Explicit proxy configuration: HTTPS_PROXY [%s]", o.Proxy)
		os.Setenv("HTTP_PROXY", o.Proxy)
		os.Setenv("HTTPS_PROXY", o.Proxy)
	}
}

func (o *Options) Validate() error {
	if o.BasePath == "" {
		return errors.New(fmt.Sprintf("`-%s` option required", optNameBasePath))
	}
	if !file.IsDirectory(o.BasePath) {
		return errors.New(fmt.Sprintf("Directory [%s] not exist.", o.BasePath))
	}
	if o.ModeAuth != "" {
		if !domain.ContainsString(modeAuthOpts, o.ModeAuth) {
			return errors.New(fmt.Sprintf("Undefined option for `-%s`: %s", optNameModeAuth, o.ModeAuth))
		}
	}
	if o.ModeSync != "" {
		syncCmds := strings.Split(o.ModeSync, ",")
		for _, x := range syncCmds {
			if !domain.ContainsString(modeSyncOpts, x) {
				return errors.New(fmt.Sprintf("Undefined option for `-%s`: %s", optNameModeSync, x))
			}
			if x == MODE_SYNC_GROUP_PROVISION && o.GroupWhiteList == "" {
				return errors.New(fmt.Sprintf("Mode `%s` requires Google Group white list file", MODE_SYNC_GROUP_PROVISION))
			}
			if x == MODE_SYNC_GROUP_PROVISION && !file.FileExistAndReadable(o.GroupWhiteList) {
				return errors.New(fmt.Sprintf("Google Group white list file [%s] not exist", o.GroupWhiteList))
			}
		}
	}
	return nil
}

func (o *Options) Usage() {
	flag.Usage()
}
