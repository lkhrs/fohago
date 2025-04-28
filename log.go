package main

import (
	"log/slog"
	"os"
)

type ServiceLogger struct {
	stdout *os.File
	file   *os.File
}

func (s *ServiceLogger) Write(event []byte) (n int, err error) {
	n, err = s.stdout.Write(event)
	if err != nil {
		return n, err
	}
	n, err = s.file.Write(event)
	return n, err
}

func ServiceLogHandler() (logHandler slog.Handler) {
	serviceLog, err := os.OpenFile("fohago.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}
	defer serviceLog.Close()
	logWriter := &ServiceLogger{
		stdout: os.Stdout,
		file:   serviceLog,
	}
	logHandler = slog.NewTextHandler(logWriter, nil)
	return
}
func AccessLogHandler() (logHandler slog.Handler) {
	accessLog, err := os.OpenFile("access_log.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		panic(err)
	}
	defer accessLog.Close()
	logHandler = slog.NewJSONHandler(accessLog, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return
}
