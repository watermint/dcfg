package dcfg

import (
	"testing"
	"github.com/watermint/dcfg/config"
)

func init() {
	config.ReloadConfigForTest()
}

func TestVerifyDropbox(*testing.T) {
	// skip
	//auth.VerifyDropbox()
}
