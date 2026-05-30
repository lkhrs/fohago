package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/lkhrs/fohago/antispam"
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

// logTurnstileFailure logs the reason for a Turnstile failure at the appropriate log level
func logTurnstileFailure(err error) {
    switch {
    case errors.Is(err, antispam.ErrInvalidInputResponse),
        errors.Is(err, antispam.ErrMissingInputResponse),
        errors.Is(err, antispam.ErrValidationFailed):
        slog.Info("Turnstile rejected submission", slog.Any("error", err))

    case errors.Is(err, antispam.ErrMissingInputSecret):
        slog.Error("Turnstile configuration error", slog.Any("error", err))

    default:
        var statusErr antispam.HTTPStatusError
        if errors.As(err, &statusErr) {
            slog.Warn("Turnstile API returned non-OK status",
                slog.Int("status", statusErr.StatusCode),
                slog.Any("error", err),
            )
            return
        }

        slog.Error("Turnstile check failed", slog.Any("error", err))
    }
}

