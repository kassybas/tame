package orchestrate

import (
	"github.com/kassybas/tame/internal/lex"
	"github.com/sirupsen/logrus"
)




func Make(path, targetName string, targetArgs []string) {
	head, err := lex.Analyse(path, targetName, targetArgs)
	if err != nil {
		logrus.Fatal(err)
	}
	err = head.Exec()
	if err !=nil {
		logrus.Fatal("Execution error:", err.Error())
	}
}
