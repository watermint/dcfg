package context

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial"
	"github.com/watermint/dcfg/cli"
	"github.com/watermint/dcfg/common/file"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/admin/directory/v1"
	"io/ioutil"
)

type ExecutionContext struct {
	// Options
	Options cli.Options

	// Dropbox Client
	DropboxClient dropbox.Api
	DropboxToken  DropboxToken

	// Google Client
	GoogleClient       *admin.Service
	GoogleClientConfig *oauth2.Config
	GoogleToken        *oauth2.Token
}

type DropboxToken struct {
	TeamManagementToken string `json:"token-team-management"`
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

func (e *ExecutionContext) CreateDropboxClientByToken(token string) dropbox.Api {
	return dropbox.Client(token, dropbox.Options{})
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
