package logging

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	if fp, err := os.OpenFile("./.sweetfdtd.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777); err == nil {
		logger.SetOutput(fp)
		logger.SetLevel(logrus.InfoLevel)
	} else {
		msg := fmt.Sprintf("can not write to log file %+v", err)
		panic(msg)
	}
}

func GetLogger() *logrus.Logger {
	return logger
}
