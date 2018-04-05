package directory

import (
	"github.com/watermint/dcfg/integration/context"
	"google.golang.org/api/admin/directory/v1"
	"testing"
)

func TestGoogleDirectory_Load(t *testing.T) {
	ctx, err := context.NewExecutionContextForTest()
	if err != nil {
		t.Skip()
	}
	ctx.InitGoogleClient()

	gd := NewGoogleDirectory(ctx)
	accounts := gd.Accounts()
	if len(accounts) < 1 {
		t.Error("No accounts loaded from Google")
	}
}

func TestGoogleDirectory_GroupSimple(t *testing.T) {
	ga := GoogleAppsMock{
		MockGroups: []*admin.Group{
			&admin.Group{
				Id:    "id-g1",
				Name:  "g1",
				Email: "g1@example.com",
			},
		},
	}
	gd := NewGoogleDirectoryForTest(&ga)

	g1, e1 := gd.Group("id-g1")
	if !e1 || g1.GroupId != "g1@example.com" || g1.GroupEmail != "g1@example.com" || g1.GroupName != "g1" {
		t.Errorf("Invalid result: %v", g1)
	}
	g2, e2 := gd.Group("g1@example.com")
	if !e2 || g2.GroupId != "g1@example.com" || g2.GroupEmail != "g1@example.com" || g2.GroupName != "g1" {
		t.Errorf("Invalid result: %v", g2)
	}
	_, e3 := gd.Group("noexistent")
	if e3 {
		t.Errorf("Invalid state")
	}
}

func TestGoogleDirectory_GroupComplex(t *testing.T) {
	gd := CreateGoogleDirectoryForIntegrationTest()

	emailsShouldExist := []string{
		"a@example.com",
		"b@example.com",
		"b2@example.com",
		"b3@example.com",
		"b@example.net",
		"c@example.com",
		"d@example.com",
		"d2@example.com",
		"d@example.org",
		"all@example.com",
		"tokyo@example.com",
		"meguro@example.com",
		"minato@example.com",
	}

	for _, e := range emailsShouldExist {
		if ex, er := gd.EmailExist(e); !ex || er != nil {
			t.Errorf("%s Should exist", e)
		}
	}
}
