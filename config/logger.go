package config

import (
	"log/slog"
	"os"
)

func InitLogger(logLevel string) error {
	var lv slog.Level
	err := lv.UnmarshalText([]byte(logLevel))
	if err != nil {
		return err
	}
	opts := &slog.HandlerOptions{
		Level: lv,
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))
	return nil
}
