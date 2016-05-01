package dsutil

import (
	"os"
	"path/filepath"
	//	"path/filepath"
	"runtime"
	"strings"

	"github.com/desal/cmd"
)

func CheckPath(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func SplitLines(s string, omitEmpty bool) []string {
	split := strings.Split(s, "\n")
	result := []string{}
	for i, s := range split {
		sr := strings.Replace(s, "\r", "", -1)
		if omitEmpty && sr == "" {
			continue
		} else if i == len(split)-1 && sr == "" {
			continue
		}
		result = append(result, sr)
	}
	return result
}

func FirstLine(s string) string {
	return strings.Replace(strings.Split(s, "\n")[0], "\r", "", -1)
}

func SanitisePath(cmdContext *cmd.Context, path string) (string, error) {
	if runtime.GOOS != "windows" {
		return path, nil
	}

	//TODO check if cygpath is available, probably do this once on init
	cygPath, err := cmdContext.Execf("cygpath -w '%s'", path)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(FirstLine(cygPath)), nil
}
