package exec

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-cmd/cmd"
)

type Options struct {
	Silent       bool // no output to stdout and stderr
	IgnoreResult bool // result of stdout and stderr is not returned (enable to save resources + you don't care about the values)
	ShieldEnv    bool // do not expose the current process' environment variables (ignore wath is in os.Environ())

	// ShellExec invocation scheme: "ShellPath ShellExtraFlags ShellCmdFlag Script"
	ShellPath       string   // path of the shell to be executed (defaults to sh)
	ShellCmdFlag    string   // flag that precedes the script string (default '-c')
	ShellExtraFlags []string // list of flags for the invocation of the shell (eg. --posix for bash)
}

const DefaultShellCmdFlag = "-c"
const DefaultShell = "sh"

func ShellExec(script string, envVars []string, opts Options) (resOut, resErr string, resRc int, err error) {
	cmdOptions := cmd.Options{
		Buffered:  !opts.IgnoreResult,
		Streaming: !opts.Silent,
	}

	// Create Cmd with options
	if opts.ShellPath == "" {
		opts.ShellPath = DefaultShell
	}
	scriptCmd := cmd.NewCmdOptions(cmdOptions, opts.ShellPath)
	if opts.ShellExtraFlags != nil {
		scriptCmd.Args = append(scriptCmd.Args, opts.ShellExtraFlags...)
	}

	// Set default command flag of shell if not set
	if opts.ShellCmdFlag == "" {
		opts.ShellCmdFlag = DefaultShellCmdFlag
	}

	scriptCmd.Args = append(scriptCmd.Args, opts.ShellCmdFlag, script)
	if opts.ShieldEnv {
		scriptCmd.Env = []string{}
	} else {
		scriptCmd.Env = os.Environ()
	}
	scriptCmd.Env = append(scriptCmd.Env, envVars...)

	if !opts.Silent {
		// Print STDOUT and STDERR lines streaming from Cmd
		go func() {
			for {
				select {
				case line := <-scriptCmd.Stdout:
					fmt.Println(line)
				case line := <-scriptCmd.Stderr:
					fmt.Fprintln(os.Stderr, line)
				}
			}
		}()
	}

	// Run and wait for Cmd to return
	status := <-scriptCmd.Start()
	if opts.IgnoreResult {
		return
	}

	resOut = strings.Join(status.Stdout, "\n")
	resErr = strings.Join(status.Stderr, "\n")
	resRc = status.Exit
	err = status.Error
	return
}
