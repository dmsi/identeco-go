package mylog

import (
	"log/slog"
	"path/filepath"

	"github.com/dmsi/identeco-go/pkg/lib/prettylog"
)

var Lg *slog.Logger

func init() {
	Lg = newLogger()
}

func newLogger() *slog.Logger {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}

	opts := &slog.HandlerOptions{
		Level:       slog.LevelDebug,
		AddSource:   true,
		ReplaceAttr: replace,
	}

	lg := slog.New(prettylog.NewHandler(opts))

	return lg
}
