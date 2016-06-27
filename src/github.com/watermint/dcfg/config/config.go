package config

import (
	"path"
	"encoding/json"
	"github.com/watermint/dcfg/explorer"
	"os"
	"runtime"
	"github.com/cihub/seelog"
)

type Config struct {
	configPath    string
	google        GoogleConfig
	dropboxToken  DropboxToken
	dropboxClient DropboxClient
}

type GoogleConfig struct {
	Domain string `json:"domain"`
}

type DropboxToken struct {
	TeamManagementToken string `json:"token-team-management"`
}

type DropboxClient struct {
	ClientId     string `json:"app-key"`
	ClientSecret string `json:"app-secret"`
}

var (
	Global Config
)

const (
	googleTokenFileName = "google_token.json"
	googleClientSecret = "google_client_secret.json"
	dropboxTokenFile = "dropbox_token.json"
	systemLog = "dcfg.log"
)

func ReloadConfigForTest() {
	_, file, _, _ := runtime.Caller(1)
	baseDir := path.Dir(path.Dir(path.Dir(file)))
	seelog.Infof("Test Base directory: %s", baseDir)
	testData := path.Join(baseDir, "test_data")
	ReloadConfig(testData)
}

func ReloadConfig(configPath string) {
	seelog.Tracef("Loading configuration: [%s]", configPath)
	loadConfig(configPath)
}

func loadConfig(configPath string) {
	Global = Config{configPath:configPath, google:GoogleConfig{}}
}

// Path to system log
func (d *Config) SystemLogFile() string {
	if d.configPath == "" {
		return systemLog
	} else {
		return path.Join(d.configPath, systemLog)
	}
}

func (d *Config) GoogleTokenFile() string {
	if d.configPath == "" {
		return googleTokenFileName
	} else {
		return path.Join(d.configPath, googleTokenFileName)
	}
}

func (d *Config) GoogleClientFile() string {
	if d.configPath == "" {
		return googleClientSecret
	} else {
		return path.Join(d.configPath, googleClientSecret)
	}
}

func (d *Config) loadConfig(file string, label string, data interface{}) {
	j, err := os.Open(file)
	if err != nil {
		explorer.Fatal("Unable to read file", label, file, err)
	}
	defer j.Close()
	err = json.NewDecoder(j).Decode(data)
	if err != nil {
		explorer.Fatal("Unable to parse file", label, file, err)
	}
}

func (d *Config) loadDropboxToken() {
	c := &DropboxToken{}
	d.loadConfig(d.DropboxTokenFile(), "Dropbox Token", c)
	d.dropboxToken = *c
}

func (d *Config) DropboxToken() DropboxToken {
	if d.dropboxToken.TeamManagementToken == "" {
		d.loadDropboxToken()
	}
	return d.dropboxToken
}

func (d *Config) DropboxTokenFile() string {
	if d.configPath == "" {
		return dropboxTokenFile
	} else {
		return path.Join(d.configPath, dropboxTokenFile)
	}
}
