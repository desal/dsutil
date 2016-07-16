package dsutil

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/desal/cmd"
	"github.com/desal/richtext"
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

func PosixPath(path string) string {
	if len(path) < 2 || runtime.GOOS != "windows" {
		return path
	}

	if path[1] == ':' {
		return "/" + strings.ToLower(string(path[0])) + filepath.ToSlash(path[2:])
	}

	return filepath.ToSlash(path)
}

var cygMounts = map[string]string{}
var cygErr error

func init() {
	cygErr = cmd.Check()
	if cygErr != nil {
		return
	}

	var mountStr string
	mountStr, _, cygErr = cmd.New("", richtext.Silenced()).Execf("mount")
	if cygErr != nil {
		return
	}

	mountLines := SplitLines(mountStr, true)

	re := regexp.MustCompile(`^(.*) on (.*) type .* \(.*\)$`)
	for _, line := range mountLines {
		splits := re.FindStringSubmatch(line)
		if len(splits) != 3 {
			cygErr = fmt.Errorf("invalid mount output: %s", line)
			return
		}
		win32Path := splits[1]
		posixPath := splits[2]
		//On windows this is completely case insensitive
		cygMounts[strings.ToLower(posixPath)] = filepath.FromSlash(win32Path)
	}
}

func CygError() error {
	return cygErr
}

//On windows, passing in a non-trivial absolute posix path (one which starts
//with /, but is not a under drive letter like /c/...) will panic if 'mount.exe'
// is not available; i.e. when running with cygwin/msys/git bash/etc not in the
//path. To gracefully handle, check in advance with dsutil.CygError() != nil
func NativePath(path string) string {
	if runtime.GOOS != "windows" || len(path) == 0 || path[0] != '/' {
		return path
	}

	if cygErr != nil && path[2] == '/' && ((path[1] >= 'a' && path[1] <= 'z') || (path[1] >= 'A' && path[1] <= 'Z')) {
		return string(path[1]) + ":" + filepath.FromSlash(path[2:])
	}

	if cygErr != nil {
		panic(cygErr)
	}

	var posixMount string
	lowerPath := strings.ToLower(path)

	for posixPath, _ := range cygMounts {
		if strings.HasPrefix(lowerPath, posixPath) && len(posixPath) > len(posixMount) {
			posixMount = posixPath
		}
	}

	//Don't want to use trim prefix as i want to presever the case after the prefix
	return cygMounts[posixMount] + filepath.FromSlash(path[len(posixMount):])
}
