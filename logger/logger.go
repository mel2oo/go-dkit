/*
 * Copyright 2022 by Mel2oo <https://github.com/saferun/owl>
 *
 * Licensed under the GNU General Public License version 3 (GPLv3)
 *
 * If you distribute GPL-licensed software the license requires
 * that you also distribute the complete, corresponding source
 * code (as defined by GPL) to that GPL-licensed software.
 *
 * You should have received a copy of the GNU General Public License
 * with this program. If not, see <https://www.gnu.org/licenses/>
 */

package logger

import (
	"fmt"
	"runtime"
	"strings"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	App    string `yaml:"app" json:"app,omitempty"`
	Level  string `yaml:"level" json:"level,omitempty"`
	Format bool   `yaml:"format" json:"format,omitempty"`

	// output file, if need.
	Output    string `yaml:"output" json:"output,omitempty"`
	MaxSize   int    `yaml:"max_size" json:"max_size,omitempty"`
	MaxAge    int    `yaml:"max_age" json:"max_age,omitempty"`
	MaxBackup int    `yaml:"max_backup" json:"max_backup,omitempty"`
}

func Init(cfg *Config) error {
	// setup log formatter
	if cfg.Format {
		logrus.SetFormatter(&SimpleFormatter{})
	} else {
		logrus.SetFormatter(&nested.Formatter{TimestampFormat: TimeFormat})
	}

	// output level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	opts := []Option{
		WithAppName(cfg.App),
		WithLevels(logrus.AllLevels[:level+1]...),
	}

	// output file with roll
	if cfg.Output != "" {
		opts = append(opts, WithWriter(&lumberjack.Logger{
			Filename:   cfg.Output,
			MaxSize:    cfg.MaxSize,
			MaxAge:     cfg.MaxAge,
			MaxBackups: cfg.MaxBackup,
		}))
	}

	logrus.AddHook(NewHook(opts...))
	return err
}

func findCaller() string {
	var (
		file string
		line int
	)

	for i := 0; i < 20; i++ {
		file, line = getCaller(i + 5)
		if !skipFile(file) {
			break
		}
	}

	return fmt.Sprintf("%s:%d", file, line)
}

var skipPrefixes = []string{"logrus/", "logrus@", "v4@", "logger/"}

func skipFile(file string) bool {
	for i := range skipPrefixes {
		if strings.HasPrefix(file, skipPrefixes[i]) {
			return true
		}
	}
	return false
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}

	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}

	return file, line
}
