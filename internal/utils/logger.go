package utils

import (
	"os"

	"github.com/rs/zerolog"
)

// Logger is the global logger instance
var Logger zerolog.Logger

// InitLogger initializes the global logger with configuration based on debug and verbose flags
func InitLogger(debug bool) {
	if debug {
		// Debug mode: raw text format, no timestamp, DEBUG level enabled
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		writer := zerolog.ConsoleWriter{
			Out:        os.Stderr,
			NoColor:    false,
			TimeFormat: "", // Empty string disables timestamp
			FormatLevel: func(i interface{}) string {
				if ll, ok := i.(string); ok {
					switch ll {
					case "debug":
						return "[DEBUG]"
					default:
						return "[" + ll + "]"
					}
				}
				return "[???]"
			},
		}
		Logger = zerolog.New(writer).Level(zerolog.DebugLevel).With().
			Logger()
	} else {
		// Silent mode: logger disabled (verbose is no-op without debug)
		zerolog.SetGlobalLevel(zerolog.Disabled)
		Logger = zerolog.New(os.Stderr).Level(zerolog.Disabled).With().
			Logger()
	}
}
