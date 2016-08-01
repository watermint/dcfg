package context

import (
	"errors"
	"github.com/cihub/seelog"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/team"
	"github.com/watermint/dcfg/cli"
	"github.com/watermint/dcfg/common/file"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/admin/directory/v1"
	"io/ioutil"
	"path"
	"runtime"
)

type ExecutionContext struct {
	// Options
	Options cli.Options

	// Dropbox Client
	DropboxClient team.Client
	DropboxToken  DropboxToken

	// Google Client
	GoogleClient       *admin.Service
	GoogleClientConfig *oauth2.Config
	GoogleToken        *oauth2.Token
}

type DropboxToken struct {
	TeamManagementToken string `json:"token-team-management"`
}

func NewDropboxToken(mgmtToken string) DropboxToken {
	return DropboxToken{
		TeamManagementToken: mgmtToken,
	}
}

func (e *ExecutionContext) CreateGoogleClientByToken(token *oauth2.Token) (*admin.Service, error) {
	context := context.Background()
	client := e.GoogleClientConfig.Client(context, token)
	service, err := admin.New(client)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (e *ExecutionContext) loadGoogleClient() error {
	client, err := e.CreateGoogleClientByToken(e.GoogleToken)
	if err != nil {
		return err
	}
	e.GoogleClient = client
	return nil
}

func (e *ExecutionContext) loadGoogleClientConfig() error {
	seelog.Tracef("Loading Google Client Config: %s", e.Options.PathGoogleClientSecret())
	json, err := ioutil.ReadFile(e.Options.PathGoogleClientSecret())
	if err != nil {
		return err
	}
	config, err := google.ConfigFromJSON(json,
		admin.AdminDirectoryUserReadonlyScope,
		admin.AdminDirectoryGroupReadonlyScope)
	if err != nil {
		return err
	}
	e.GoogleClientConfig = config
	return nil
}

func (e *ExecutionContext) loadGoogleToken() error {
	path := e.Options.PathGoogleToken()
	token := oauth2.Token{}
	_, err := file.LoadJSON(path, &token)
	if err != nil {
		return err
	}
	e.GoogleToken = &token
	return nil
}

func (e *ExecutionContext) CreateDropboxClientByToken(token string) team.Client {
	config := dropbox.Config{
		Token: token,
	}
	return team.New(config)
}

func (e *ExecutionContext) loadDropboxClient() error {
	e.DropboxClient = e.CreateDropboxClientByToken(e.DropboxToken.TeamManagementToken)
	return nil
}

func (e *ExecutionContext) loadDropboxToken() error {
	path := e.Options.PathDropboxToken()
	token := DropboxToken{}
	_, err := file.LoadJSON(path, &token)
	if err != nil {
		return err
	}
	e.DropboxToken = token
	return nil
}

func (e *ExecutionContext) InitDropboxClient() error {
	if err := e.loadDropboxToken(); err != nil {
		return err
	}
	if err := e.loadDropboxClient(); err != nil {
		return err
	}
	return nil
}

func (e *ExecutionContext) InitGoogleAuth() error {
	if err := e.loadGoogleClientConfig(); err != nil {
		return err
	}
	return nil
}

func (e *ExecutionContext) InitGoogleClient() error {
	if err := e.loadGoogleClientConfig(); err != nil {
		return err
	}
	if err := e.loadGoogleToken(); err != nil {
		return err
	}
	if err := e.loadGoogleClient(); err != nil {
		return err
	}
	return nil
}

func (e *ExecutionContext) InitForSync() error {
	if err := e.InitDropboxClient(); err != nil {
		return err
	}
	if err := e.InitGoogleClient(); err != nil {
		return err
	}
	return nil
}

func getTestBasePath() string {
	_, file, _, _ := runtime.Caller(1)
	projectRoot := path.Dir(path.Dir(path.Dir(file)))
	return path.Join(projectRoot, "test_data")
}

func NewExecutionContextForTest() (ExecutionContext, error) {
	basePath := getTestBasePath()
	if file.IsDirectory(basePath) {
		ctx := ExecutionContext{
			Options: cli.Options{
				BasePath: getTestBasePath(),
			},
		}
		return ctx, nil
	}
	return ExecutionContext{}, errors.New("Test directory not found")
}
