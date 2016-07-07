package explorer

import (
	"fmt"
	"github.com/cihub/seelog"
)

var (
	reportSuccess []string
	reportFailure []string
)

func init() {
	reportSuccess = []string{}
	reportFailure = []string{}
}

func ReportSuccess(format string, values ...interface{}) {
	reportSuccess = append(reportSuccess, fmt.Sprintf(format, values...))
}

func ReportFailure(format string, values ...interface{}) {
	reportFailure = append(reportFailure, fmt.Sprintf(format, values...))
}

func reportLine(format string, args ...interface{}) {
	seelog.Infof(format, args...)
}

func Report() {
	if len(reportSuccess) == 0 && len(reportFailure) == 0 {
		reportLine("No update.")
	} else {
		if len(reportSuccess) > 0 {
			for i, line := range reportSuccess {
				reportLine("Success: [%d] %s", i+1, line)
			}
		}
		if len(reportFailure) > 0 {
			for i, line := range reportFailure {
				reportLine("Failure: [%d] %s", i+1, line)
			}
		}
	}
	reportLine("Done")
}
