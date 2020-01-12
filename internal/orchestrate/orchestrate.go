package orchestrate

import (
	"fmt"
	"os"

	"github.com/kassybas/tame/internal/lex"

	"github.com/kassybas/tame/internal/vartable"

	"github.com/sirupsen/logrus"
)

func Make(path, targetName string, targetArgs []string) {

	root, ctx, err := lex.PrepareStep(path, targetName, targetArgs)
	if err != nil {
		logrus.Fatal(err)
	}

	// TODO: put cli args in here
	status := root.RunStep(ctx, vartable.NewVarTable())
	if status.Err != nil {
		logrus.Fatal("error:\n\t", status.Err.Error())
	}
	// pass through the status code
	if status.Stdstatus != 0 {
		fmt.Fprintf(os.Stdout, "tame: *** [%s] Error %d\n", targetName, status.Stdstatus)
	}

	os.Exit(status.Stdstatus)
}
