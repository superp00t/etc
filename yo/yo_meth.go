package yo

import (
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type Lg struct {
	Level  int64
	Logger Logger
}

func (l *Lg) Log(ld LogData) {
	if l.Level >= Int64G("y") {
		l.Logger.Log(ld)
	}
}

func (l *Lg) Println(args ...interface{}) {
	l.Log(LogData{
		Level: LDebug,
		Time:  time.Now(),
		Data:  fmt.Sprintln(args...),
	})
}

func (l *Lg) Printf(f string, args ...interface{}) {
	l.Log(LogData{
		Level: LDebug,
		Time:  time.Now(),
		Data:  fmt.Sprintf(f, args...),
	})
}

func (l *Lg) Ok(args ...interface{}) {
	l.Log(LogData{
		Level: LOkay,
		Time:  time.Now(),
		Data:  fmt.Sprintln(args...),
	})
}

func (l *Lg) Okf(f string, args ...interface{}) {
	l.Log(LogData{
		Level: LOkay,
		Time:  time.Now(),
		Data:  fmt.Sprintf(f, args...),
	})
}

func (l *Lg) Warn(args ...interface{}) {
	l.Log(LogData{
		Level: LWarn,
		Time:  time.Now(),
		Data:  fmt.Sprintln(args...),
	})
}

func (l *Lg) Warnf(f string, args ...interface{}) {
	l.Log(LogData{
		Level: LWarn,
		Time:  time.Now(),
		Data:  fmt.Sprintf(f, args...),
	})
}

func (l *Lg) Spew(v interface{}) {
	l.Log(LogData{
		Level: LDebug,
		Time:  time.Now(),
		Data:  spew.Sdump(v),
	})
}

func L(i int64) *Lg {
	return &Lg{i, stdLogger}
}

func Println(args ...interface{}) {
	L(0).Println(args...)
}

func Printf(f string, args ...interface{}) {
	L(0).Printf(f, args...)
}

func Spew(v interface{}) {
	L(0).Spew(v)
}

func Puke(v interface{}) {
	Fatal(spew.Sdump(v))
}

func Ok(args ...interface{}) {
	L(0).Ok(args...)
}

func Okf(f string, args ...interface{}) {
	L(0).Okf(f, args...)
}

func Warn(args ...interface{}) {
	L(0).Warn(args...)
}

func Warnf(f string, args ...interface{}) {
	L(0).Okf(f, args...)
}

func Fatalf(f string, args ...interface{}) {
	stdLogger.Log(LogData{
		Level: LFatal,
		Time:  time.Now(),
		Data:  fmt.Sprintf(f, args...),
	})
	os.Exit(-1)
}

func Fatal(args ...interface{}) {
	stdLogger.Log(LogData{
		Level: LFatal,
		Time:  time.Now(),
		Data:  fmt.Sprintln(args...),
	})
	os.Exit(-1)
}
