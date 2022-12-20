package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func LoggerInit(file string) {
	f, err := os.OpenFile(
		file,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter, f)
	Logger = zerolog.New(multi).With().Timestamp().Caller().Logger()
}
