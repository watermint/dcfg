package dcfg

import (
	"testing"
	"dcfg/auth"
	"dcfg/config"
)

func init() {
	config.ReloadConfigForTest()
}

func TestVerifyDropbox(*testing.T) {
	auth.VerifyDropbox()
}
