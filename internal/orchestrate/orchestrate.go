package orchestrate

import (
	"github.com/kassybas/mate/internal/lex"
)

func Make(path, targetName string, targetArgs []string) {
	lex.Analyse(path, targetName, targetArgs)
	// head, err := lex.Analyse(path, targetName, targetArgs)
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// err = head.Exec()
	// if err != nil {
	// 	logrus.Fatal("Execution error:", err.Error())
	// }
}
