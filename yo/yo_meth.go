package yo

import (
	"fmt"
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
