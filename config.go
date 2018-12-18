// Simple console progress bars
package pb

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func init() {
	// extract the $SHELL environment variable
	// is the standard mechanism POSIX
	// for communicating the shell your program
	// is running within.
	//
	// like everything POSIX it doesn't work
	// 100% of time, as if you launch a
	// different shell from your normal bash
	// session SHELL will still contain bash,
	// not say `zsh` which was launched in a
	// weird `$ zsh my_program` manner.
	//
	// for additional reading:
	// https://stackoverflow.com/questions/3327013/how-to-determine-the-current-shell-im-working-on
	shell := os.Getenv("SHELL")
	switch runtime.GOOS {
	case "windows":
		// we need to check if we're
		// in a windows emulated
		// shell environment
		if shell == "" {
			// this is just a windows shell
			// we are done
			RequireWindowsCalls = true
			ClearLinePrefixString = "\r"
			ClearLineSuffixString = ""
			return
		}
		// since the environment is telling
		// us this is a POSIX like
		// environment we will treat it like
		// such.
		RequireWindowsCalls = false

		// first we need to check if we
		// are really being emulated
		ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*250)
		defer cancel()
		// check if we have access to `stty`
		//
		// cygwin will, but git 4 windows and mtty won't
		// citation:
		// - personal research
		UseSTTYWindows = exec.CommandContext(ctx, "stty", "size").Run() == nil
		fallthrough
	case "solaris", "darwin", "linux", "dragonfly", "freebsd", "netbsd", "openbsd", "plan9":
		if isBASHLike(shell) {
			ClearLinePrefixString = "\r"
			ClearLineSuffixString = ""
			return
		}
		// oh you arent using `bash`, `dash`, or `sh`
		// we'll set an ugly prefix we know that works
		// with `zsh` and just pray it works in `ksh`
		// `csh`, `tcsh`, and `fish`
		//
		// this will likely get more diverse as time goes on
		ClearLinePrefixString = fmt.Sprintf("%c[%dA%c[K\r", 27, 1, 27)
		ClearLineSuffixString = fmt.Sprintf("%c[1i\n", 27)
		return
	default:
		panic(fmt.Sprintf("we not not yet support %s", runtime.GOOS))
	}
}

// isBASHLike checks for a bash like shells which are the most common
func isBASHLike(shell string) bool {
	if os.Getenv("TERMINATOR_UUID") != "" {
		// don't treat terminator like BASH
		return false
	}
	shell = strings.TrimSpace(shell)
	// dash, bash, and sh are mostly identical for our purposes
	return strings.Contains(shell, "bash") ||
		strings.Contains(shell, "dash") ||
		shell == "/bin/sh" ||
		shell == "/sbin/sh" ||
		shell == "/usr/bin/sh" ||
		shell == "/usr/sbin/sh"
}

// ClearLinePrefixString is used to set the clear prefix
var ClearLinePrefixString string

// ClearLineSuffixString is used to set the clear suffix
var ClearLineSuffixString string

// RequireWindowsCalls is done to check
// if while we maybe in windows, we have
// access to
var RequireWindowsCalls bool

// Use STTY can be used on windows in
// order to ensure sizing works like
// a normal boring unix terminal
var UseSTTYWindows bool
