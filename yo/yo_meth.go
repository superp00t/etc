package yo

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"os"
	"time"
)

func Println(args ...interface{}) {
	stdLogger.Log(LogData{
		Level: LDebug,
		Time:  time.Now(),
		Data:  fmt.Sprintln(args...),
	})
}

func Spew(v interface{}) {
	Println(spew.Sdump(v))
}

func Puke(v interface{}) {
	Fatal(spew.Sdump(v))
}

func Ok(args ...interface{}) {
	stdLogger.Log(LogData{
		Level: LOkay,
		Time:  time.Now(),
		Data:  fmt.Sprintln(args...),
	})
}

func Warn(args ...interface{}) {
	stdLogger.Log(LogData{
		Level: LWarn,
		Time:  time.Now(),
		Data:  fmt.Sprintln(args...),
	})
}

func Fatal(args ...interface{}) {
	stdLogger.Log(LogData{
		Level: LFatal,
		Time:  time.Now(),
		Data:  fmt.Sprintln(args...),
	})
	os.Exit(-1)
}

func Printf(f string, args ...interface{}) {
	stdLogger.Log(LogData{
		Level: LDebug,
		Time:  time.Now(),
		Data:  fmt.Sprintf(f, args...),
	})
}

func Okf(f string, args ...interface{}) {
	stdLogger.Log(LogData{
		Level: LOkay,
		Time:  time.Now(),
		Data:  fmt.Sprintf(f, args...),
	})
}
