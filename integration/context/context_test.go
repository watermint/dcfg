package context

import (
	"testing"
)

func TestExecutionContext_InitDropboxClient(t *testing.T) {
	ctx, err := NewExecutionContextForTest()
	if err != nil {
		t.Skip()
	}

	ctx.InitDropboxClient()
}

func TestExecutionContext_InitGoogleClient(t *testing.T) {
	ctx, err := NewExecutionContextForTest()
	if err != nil {
		t.Skip()
	}

	ctx.loadGoogleClient()
}
