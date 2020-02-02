package orchestrate

import (
	"os"

	"github.com/kassybas/tame/internal/build/compile"

	"github.com/kassybas/tame/internal/vartable"

	"github.com/sirupsen/logrus"
)

func Make(path, targetName string, targetArgs map[string]interface{}) {

	root, ctx, err := compile.PrepareStep(path, targetName, targetArgs)
	if err != nil {
		logrus.Errorf(err.Error())
		os.Exit(1)
	}

	// TODO: put cli args in here
	status := root.RunStep(ctx, vartable.NewVarTable())
	if status.Err != nil {
		logrus.Errorf(status.Err.Error())
		if status.Stdstatus == 0 {
			status.Stdstatus = 1
		}
	}
	// pass through the status code
	if status.Stdstatus != 0 {
		logrus.Errorf("tame: *** [%s] Error %d\n", targetName, status.Stdstatus)
	}

	os.Exit(status.Stdstatus)
}
