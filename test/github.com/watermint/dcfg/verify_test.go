package dcfg

import (
	"github.com/watermint/dcfg/config"
	"testing"
)

func init() {
	config.ReloadConfigForTest()
}

func TestVerifyDropbox(*testing.T) {
	// skip
	//auth.VerifyDropbox()
}
