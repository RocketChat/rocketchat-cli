package internal

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type LoggerGenerator struct {
	logwriter io.Writer
}

func NewLoggerGenerator(logfilewriter io.Writer) *LoggerGenerator {
	return &LoggerGenerator{logwriter: logfilewriter}
}

type logfileHook struct {
	formatter logrus.Formatter
	writer    io.Writer
}

var _ logrus.Hook = &logfileHook{}

func (h *logfileHook) Fire(entry *logrus.Entry) error {
	b, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}
	h.writer.Write(b)
	return nil
}

func (*logfileHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
	}
}

type formatterWrapperForDebugShutUp struct {
	formatter logrus.Formatter
	debug     bool
}

func (f formatterWrapperForDebugShutUp) Format(entry *logrus.Entry) ([]byte, error) {
	if !f.debug && entry.Level == logrus.DebugLevel {
		return nil, nil
	}

	// clear line to avoid progress bar copying
	os.Stdout.Write([]byte("\033[2K"))

	return f.formatter.Format(entry)
}

type Logger struct {
	fields         logrus.Fields
	l              *logrus.Logger
	plainformatter *formatterWrapperForDebugShutUp

	context   string
	generator *LoggerGenerator
}

func (lg *LoggerGenerator) NewLogger(context string, debug bool) *Logger {
	l := logrus.New()

	// always set debug level here, but stdout logs are handled separately through formatter
	l.SetLevel(logrus.DebugLevel)

	formatter := &formatterWrapperForDebugShutUp{
		formatter: &logrus.TextFormatter{},
		debug:     debug,
	}

	l.SetFormatter(formatter)

	l.SetOutput(os.Stdout)

	l.AddHook(&logfileHook{
		formatter: &logrus.JSONFormatter{},
		writer:    lg.logwriter,
	})

	return &Logger{
		l:              l,
		fields:         map[string]interface{}{"context": context},
		plainformatter: formatter,
		context:        context,
		generator:      lg,
	}

}

func (l *Logger) Errorf(err error, format string, args ...interface{}) {
	l.l.WithFields(l.fields).WithField("error", err).Errorf(format, args...)
}

func (l *Logger) Fatalf(err error, format string, args ...interface{}) {
	l.l.WithFields(l.fields).WithField("error", err).Fatalf(format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.l.WithFields(l.fields).Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.l.WithFields(l.fields).Warnf(format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.l.WithFields(l.fields).Debugf(format, args...)
}

func (l *Logger) Debug(debug bool) {
	l.plainformatter.debug = debug
}

func (l *Logger) NewContext(context string) *Logger {
	return l.generator.NewLogger(l.context+"::"+context, l.plainformatter.debug)
}
