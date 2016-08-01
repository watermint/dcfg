package dispatch

import (
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/cli/explorer"
	"github.com/watermint/dcfg/integration/auth"
	"github.com/watermint/dcfg/integration/context"
	"github.com/watermint/dcfg/sync/groupsync"
	"github.com/watermint/dcfg/sync/usersync"
)

func DispatchAuth(context context.ExecutionContext) {
	switch {
	case context.Options.IsModeAuthGoogle():
		if err := context.InitGoogleAuth(); err != nil {
			seelog.Errorf("Initialisation failure: %v", err)
			explorer.FatalShutdown("Please review file content of: %s", context.Options.PathGoogleClientSecret())
		}
		seelog.Trace("Start Auth Sequence: Google")
		auth.AuthGoogle(context)
	case context.Options.IsModeAuthDropbox():
		seelog.Trace("Start Auth Sequence: Dropbox")
		auth.AuthDropbox(context)
	}
}

func DispatchSync(context context.ExecutionContext) {
	if err := context.InitForSync(); err != nil {
		seelog.Errorf("Initialisation failure: %v", err)
		explorer.FatalShutdown("Please review configuration")
	}
	if context.Options.IsModeSyncUserProvision() {
		seelog.Trace("Start Sync: User Provision")
		seelog.Infof("Provisioning Users (Google Users -> Dropbox Users)")
		userSync := usersync.NewUserSync(context)
		userSync.SyncProvision()
	}
	if context.Options.IsModeSyncUserDeprovision() {
		seelog.Trace("Start Sync: User Deprovision")
		seelog.Infof("Deprovisioning Users (Google Users -> Dropbox Users)")
		userSync := usersync.NewUserSync(context)
		userSync.SyncDeprovision()
	}
	if context.Options.IsModeGroupProvision() {
		seelog.Trace("Start Sync: Group Provision")
		seelog.Infof("Syncing Group (Google Group -> Dropbox Group)")
		groupSync := groupsync.NewGroupSync(context)
		groupSync.SyncFromList(context)
	}
}

func Dispatch(context context.ExecutionContext) {
	defer explorer.Report()

	switch {
	case context.Options.IsModeAuth():
		DispatchAuth(context)
	case context.Options.IsModeSync():
		DispatchSync(context)
	}
}
