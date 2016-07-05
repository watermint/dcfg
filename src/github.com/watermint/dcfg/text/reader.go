package text

import (
	"bufio"
	"github.com/cihub/seelog"
	"io"
	"os"
	"strings"
)

func ReadLines(filePath string) (lines []string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		seelog.Warnf("Unable to load file: file[%s] err[%s]", filePath, err)
		return []string{}, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		lineRaw, _, err := r.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			seelog.Warnf("Unable to load file: file[%s] err[%s]", filePath, err)
			return []string{}, err
		}
		line := strings.TrimSpace(DecodeUnicodeSequence(lineRaw))
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}
