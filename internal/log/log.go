package log

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

type Color int8

const (
	ColorAlways Color = iota
	ColorNever
	ColorAuto
)

func NewLogger(debug bool, color Color) *slog.Logger {
	return slog.New(newHandler(debug, color))
}

func newHandler(debug bool, color Color) slog.Handler {
	destination := os.Stderr

	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	noColor := false
	switch color {
	case ColorNever:
		noColor = true
	case ColorAuto:
		noColor = !isatty.IsTerminal(destination.Fd())
	}

	options := tint.Options{
		Level:      level,
		NoColor:    noColor,
		TimeFormat: time.TimeOnly,
	}

	return tint.NewHandler(colorable.NewColorable(destination), &options)
}
