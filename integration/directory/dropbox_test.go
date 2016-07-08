package directory

import (
	"github.com/watermint/dcfg/integration/context"
	"testing"
)

func TestDropboxDirectory(t *testing.T) {
	ctx, err := context.NewExecutionContextForTest()
	if err != nil {
		t.Skip()
	}
	ctx.InitDropboxClient()

	dd := DropboxDirectory{
		ExecutionContext: ctx,
	}

	dd.Load()
	accounts := dd.Accounts()
	if len(accounts) < 1 {
		t.Error("No accounts loaded from Dropbox")
	}
}
