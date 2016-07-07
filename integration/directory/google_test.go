package directory

import (
	"testing"
	"github.com/watermint/dcfg/integration/context"
)

func TestGoogleDirectory_Load(t *testing.T) {
	ctx, err := context.NewExecutionContextForTest()
	if err != nil {
		t.Skip()
	}
	ctx.InitGoogleClient()

	gd := GoogleDirectory{
		ExecutionContext: ctx,
	}

	gd.Load()
	accounts := gd.Accounts()
	if len(accounts) < 1 {
		t.Error("No accounts loaded from Google")
	}
}
