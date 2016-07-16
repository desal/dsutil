package dsutil_test

import (
	"testing"

	"github.com/desal/dsutil"
	"github.com/stretchr/testify/assert"
)

//These tests aren't particularly portable.

func TestNativePathCyg(t *testing.T) {
	//Works on msys2 default install
	assert.Equal(t, "C:\\users", dsutil.NativePath("/c/users"))
	assert.NotPanics(t, func() {
		assert.Equal(t, "C:\\msys64\\usr\\bin\\core_perl", dsutil.NativePath("/bin/core_perl"))
	})
}

func TestNativePathWinConsole(t *testing.T) {
	//Command prompt
	assert.Equal(t, "c:\\users", dsutil.NativePath("/c/users"))
	assert.Panics(t, func() {
		dsutil.NativePath("/bin/core_perl")
	})

}

func TestPosixPath(t *testing.T) {
	assert.Equal(t, "/c/users", dsutil.PosixPath("c:\\users"))
	assert.Equal(t, "/c/msys64/usr/bin/core_perl", dsutil.PosixPath("c:\\msys64\\usr\\bin\\core_perl"))
}
