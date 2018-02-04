package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/Guazi-inc/seed/logger/color"
)

const (
	levelError = iota
	levelFatal
	levelSuccess
	levelWarn
	levelInfo
)

// Logger logs logging records to the specified io.Writer
type Logger struct {
	mu     sync.Mutex
	output io.Writer
	buf    []byte // for accumulating text to write
}

func New(out io.Writer) *Logger {
	return &Logger{output: out}
}

var log = New(os.Stderr)

// SetOutput sets the logger output destination
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = colors.NewColorWriter(w)
}

func (l *Logger) getColorLevel(level int) string {
	switch level {
	case levelFatal:
		return colors.RedBold("[FATAL]   ")
	case levelInfo:
		return colors.BlueBold("[INFO]    ")
	case levelError:
		return colors.RedBold("[ERROR]   ")
	case levelWarn:
		return colors.YellowBold("[WARN]    ")
	case levelSuccess:
		return colors.GreenBold("[SUCCESS] ")
	default:
		panic("logger: invalid log level")
	}
}

func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
func (l *Logger) formatHeader(buf *[]byte, prefix string, t time.Time) {
	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '/')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '/')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')
	//time
	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)
	*buf = append(*buf, ' ')
	//prefix level
	*buf = append(*buf, prefix...)
	*buf = append(*buf, ": "...)
}

func (l *Logger) Output(calldepth int, s string, level int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, l.getColorLevel(level), time.Now())
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.output.Write(l.buf)
	return err
}

// Info outputs an information log message
func Info(message string) {
	log.Output(2, message, levelInfo)
}

// Infof outputs a formatted information log message
func Infof(message string, vars ...interface{}) {
	log.Output(2, fmt.Sprintf(message, vars...), levelInfo)
}

// Warn outputs a warning log message
func Warn(message string) {
	log.Output(2, message, levelWarn)
}

// Warnf outputs a formatted warning log message
func Warnf(message string, vars ...interface{}) {
	log.Output(2, fmt.Sprintf(message, vars...), levelWarn)
}

// Error outputs an error log message
func Error(message string) {
	log.Output(2, message, levelError)
}

// Errorf outputs a formatted error log message
func Errorf(message string, vars ...interface{}) {
	log.Output(2, fmt.Sprintf(message, vars...), levelError)
}

// Fatal outputs a fatal log message and exists
func Fatal(message string) {
	log.Output(2, message, levelFatal)
	os.Exit(255)
}

// Fatalf outputs a formatted log message and exists
func Fatalf(message string, vars ...interface{}) {
	log.Output(2, fmt.Sprintf(message, vars...), levelFatal)
	os.Exit(255)
}

// Success outputs a success log message
func Success(message string) {
	log.Output(2, message, levelSuccess)

}

// Successf outputs a formatted success log message
func Successf(message string, vars ...interface{}) {
	log.Output(2, fmt.Sprintf(message, vars...), levelSuccess)
}
