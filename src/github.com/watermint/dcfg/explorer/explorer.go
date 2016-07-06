package explorer

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/watermint/dcfg/cli"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
)

var (
	startupSystemLog bool
)

const (
	seeLogXmlTemplate = `
	<seelog type="adaptive" mininterval="200000000" maxinterval="1000000000" critmsgcount="5">
	<formats>
    		<format id="detail" format="date:%%Date(2006-01-02T15:04:05Z07:00)%%tloc:%%File:%%FuncShort:%%Line%%tlevel:%%Level%%tmsg:%%Msg%%n" />
    		<format id="short" format="%%Time [%%LEV] %%Msg%%n" />
	</formats>
	<outputs formatid="detail">
    		<filter levels="trace,info,warn,error,critical">
        		<rollingfile formatid="detail" filename="%s/dcfg.log" type="size" maxsize="52428800" maxrolls="7" />
    		</filter>
		<filter levels="info,warn,error,critical">
        		<console formatid="short" />
    		</filter>
    	</outputs>
	</seelog>
	`
)

func replaceLogger(basePath string, appVersion string) {
	seeLogXml := fmt.Sprintf(seeLogXmlTemplate, basePath)
	logger, err := seelog.LoggerFromConfigAsString(seeLogXml)
	if err != nil {
		log.Fatalln("Failed to load logger", err.Error())
	}
	seelog.ReplaceLogger(logger)
	seelog.Info("dcfg version: ", appVersion)
}

func FatalShutdown(suggestedWorkaround string, values ...interface{}) {
	seelog.Errorf("Suggested workaround:")
	seelog.Errorf(suggestedWorkaround, values...)
	seelog.Flush()
	os.Exit(1)
}

func verifyNetwork(host string) {
	seelog.Tracef("Verifying network reachability to [%s]", host)
	resp, err := http.Head(host)

	if resp == nil {
		seelog.Tracef("Response: nil")
	} else {
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			seelog.Tracef("Status code: %d", resp.StatusCode)
			seelog.Tracef("Status: %s", resp.Status)
			seelog.Tracef("TLS Cipher Suite: %d", resp.TLS.CipherSuite)
			seelog.Tracef("TLS Server name: %s", resp.TLS.ServerName)
			seelog.Tracef("TLS Version: %d", resp.TLS.Version)
			seelog.Tracef("Proto: %s", resp.Proto)
			seelog.Tracef("Content length: %d", resp.ContentLength)
		} else {
			respLines := string(respDump)
			for i, l := range strings.Split(respLines, "\n") {
				seelog.Tracef("Response[%d]: %s", i, l)
			}
		}
	}
	if err != nil {
		seelog.Errorf("Error during verify connecto to host [%s]: %v", host, err)
		FatalShutdown("Please check you network configuration to host [%s]", host)
	}
	seelog.Tracef("Verification finished without error: host[%s]", host)
	seelog.Infof("Network test: Success: %s", host)
}

func logOs() {
	hostname, _ := os.Hostname()
	seelog.Tracef("OS: Hostname: %s", hostname)

	wd, _ := os.Getwd()
	seelog.Tracef("OS: Working directory: %s", wd)
	seelog.Tracef("OS: Executor UID: %d", os.Geteuid())
}

func logRuntime() {
	seelog.Tracef("Runtime: Operating system: %s", runtime.GOOS)
	seelog.Tracef("Runtime: Architecture: %s", runtime.GOARCH)
	seelog.Tracef("Runtime: Go Version: %s", runtime.Version())
	seelog.Tracef("Runtime: Num CPUs: %d", runtime.NumCPU())
}

func logEnv() {
	for i, a := range os.Args[1:] {
		seelog.Tracef("Arg[%d]: [%s]", i, a)
	}
	for _, e := range os.Environ() {
		seelog.Tracef("Env: %s", e)
	}
}

func logSystem() {
	if !startupSystemLog {
		logOs()
		logRuntime()
		logEnv()
	}
	startupSystemLog = true
}

func logNetwork() {
	verifyNetwork("https://www.googleapis.com")
	verifyNetwork("https://www.dropbox.com")
	verifyNetwork("https://api.dropboxapi.com")
	verifyNetwork("https://content.dropboxapi.com")
	verifyNetwork("https://notify.dropboxapi.com")
}

// Start the explorer
func Start(options cli.Options, appVersion string) {
	replaceLogger(options.BasePath, appVersion)
	logSystem()
	options.UpdateEnv()
	logNetwork()
}
