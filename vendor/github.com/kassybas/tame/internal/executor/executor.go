package executor

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)
import gcmd "github.com/go-cmd/cmd"

func ioReadCloserToStr(closer io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(closer)
	return buf.String()
}

func executeScriptVerboseNew(cmd *gcmd.Cmd) (string, string, string, error) {
	var stdout, stderr string
	var outCount, errCount int
	go func() {
		for {
			select {
			case line := <-cmd.Stdout:
				outCount++
				fmt.Println(line)
			case line := <-cmd.Stderr:
				errCount++
				fmt.Fprintln(os.Stderr, line)
			}
		}
	}()
	<-cmd.Start()

	// Cmd has finished but wait for goroutine to print all lines
	for len(cmd.Stdout) > 0 || len(cmd.Stderr) > 0 {
		time.Sleep(10 * time.Millisecond)
	}
	finalStatus := cmd.Status()

	// handle leftovers (ending output without linebreak)
	if len(finalStatus.Stdout) > outCount{
		fmt.Println(finalStatus.Stdout[len(finalStatus.Stdout)-1])
	}
	if len(finalStatus.Stderr) > errCount{
		fmt.Println(finalStatus.Stderr[len(finalStatus.Stderr)-1])
	}
	stdout = strings.Join(finalStatus.Stdout, "\n")
	stderr = strings.Join(finalStatus.Stderr, "\n")

	return stdout, stderr, strconv.Itoa(finalStatus.Exit), finalStatus.Error
}
func executeScriptSilent(cmd *gcmd.Cmd) (string, string, string, error) {
	status := <-cmd.Start()

	outStr := strings.Join(status.Stdout, "\n")
	errStr := strings.Join(status.Stderr, "\n")
	return outStr, errStr, strconv.Itoa(status.Exit), status.Error
}

func ExecuteScript(name string, script string, varStrings []string, shellPath string, silent bool, shieldEnv bool) (string, string, string, error) {
	logrus.Debug("Begin script execution:", name)
	var cmdOptions gcmd.Options
	if silent {
		cmdOptions = gcmd.Options{
			Buffered:  true,
			Streaming: false,
		}
	} else {
		cmdOptions = gcmd.Options{
			Buffered:  true,
			Streaming: true,
		}
	}

	cmd := gcmd.NewCmdOptions(cmdOptions, shellPath, "-c")
	cmd.Args = append(cmd.Args, script)
	if !shieldEnv {
		cmd.Env = append(cmd.Env, os.Environ()...)
	} else{
		cmd.Env = []string{}
	}
	cmd.Env = append(cmd.Env, varStrings...)

	if silent {
		return executeScriptSilent(cmd)
	}

	return executeScriptVerboseNew(cmd)
}
