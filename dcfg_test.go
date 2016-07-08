package main

import (
	"github.com/watermint/dcfg/cli/explorer"
	"github.com/watermint/dcfg/common/file"
	"github.com/watermint/dcfg/integration/context"
	"github.com/watermint/dcfg/sync/groupsync"
	"github.com/watermint/dcfg/sync/usersync"
	"path"
	"testing"
)

func TestIntegration(t *testing.T) {
	ctx, err := context.NewExecutionContextForTest()
	if err != nil {
		t.Skip()
	}
	ctx.InitForSync()
	ctx.Options.DryRun = false

	groupListFile := path.Join(ctx.Options.BasePath, "group_list.csv")
	ctx.Options.GroupWhiteList = groupListFile

	if !file.FileExistAndReadable(groupListFile) {
		t.Skip("Group list file for test not found")
	}

	client := ctx.DropboxClient
	_, err = client.GetInfo()
	if err != nil {
		t.Errorf("Fail: GetInfo: %v", err)
	}

	groupSync := groupsync.NewGroupSync(ctx)
	groupSync.SyncFromList(ctx)

	userSync := usersync.NewUserSync(ctx)
	userSync.SyncProvision()
	userSync.SyncDeprovision()

	explorer.Report()
}
