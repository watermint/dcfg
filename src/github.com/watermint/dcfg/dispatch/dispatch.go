package dispatch

import (
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/auth"
	"github.com/watermint/dcfg/cli"
	"github.com/watermint/dcfg/groupsync"
	"github.com/watermint/dcfg/usersync"
)

func DispatchAuth(context cli.ExecutionContext) {
	switch {
	case context.Options.IsModeAuthGoogle():
		seelog.Trace("Start Auth Sequence: Google")
		auth.AuthGoogle(context)
	case context.Options.IsModeAuthDropbox():
		seelog.Trace("Start Auth Sequence: Dropbox")
		auth.AuthDropbox(context)
	}
}

func DispatchSync(context cli.ExecutionContext) {
	context.InitForSync()
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

func Dispatch(context cli.ExecutionContext) {
	switch {
	case context.Options.IsModeAuth():
		DispatchAuth(context)
	case context.Options.IsModeSync():
		DispatchSync(context)
	}
}
