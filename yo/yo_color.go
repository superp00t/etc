package yo

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

var mu = new(sync.Mutex)

var icons = map[LogLevel]string{
	LDebug: "*",
	LWarn:  "!",
	LFatal: "💀",
	LOkay:  "✓",
}

var colors = map[LogLevel]color.Attribute{
	LDebug: color.FgWhite,
	LWarn:  color.FgYellow,
	LFatal: color.FgRed,
	LOkay:  color.FgGreen,
}

type ycolor struct {
	w  io.Writer
	in chan LogData
}

func (c *ycolor) log(l LogData) {
	mu.Lock()
	color.Set(color.FgCyan)
	fmt.Printf("[%s] ", printTime(l.Time))
	color.Set(colors[l.Level])
	fmt.Printf("[%s] ", icons[l.Level])
	color.Set(color.FgWhite)
	fmt.Print(l.Data)
	color.Set(color.FgWhite)
	mu.Unlock()
}

// stolen from golang log
func printTime(t time.Time) string {
	b := []byte{}
	buf := &b
	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '/')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '/')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')
	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)
	return string(b)
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

func (c *ycolor) Log(l LogData) {
	c.log(l)
}

func NewConsoleColor() Logger {
	c := new(ycolor)
	c.in = make(chan LogData)
	c.w = os.Stdout
	go func() {
		for {
			ld := <-c.in
			c.log(ld)
		}
	}()
	return c
}
