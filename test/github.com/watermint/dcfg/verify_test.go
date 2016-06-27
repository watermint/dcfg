package dcfg

import (
	"testing"
	"github.com/watermint/dcfg/auth"
	"github.com/watermint/dcfg/config"
)

func init() {
	config.ReloadConfigForTest()
}

func TestVerifyDropbox(*testing.T) {
	auth.VerifyDropbox()
}
