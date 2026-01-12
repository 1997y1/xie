// Source code file, created by Developer@YAN_YING_SONG.

package console

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

const logKernelTimeLayout = "0102 15:04:05.999"

// Options Log configuration options
type options struct {
	DR, D, S, I, W, E bool // Log level switch.
	Print             bool // Standard output switch.
}

var Control = options{true, true, true, true, true, true, true}

var Logger *lumberjack.Logger

func ERROR(err error) {
	if !Control.E {
		return
	}
	errMsg := err.Error()
	b := logKernel("ERROR  ", []byte(errMsg), Control.Print)
	if Logger != nil {
		_, _ = Logger.Write(b)
	}
}

func DEBUG(txt string, a ...interface{}) {
	if !Control.D {
		return
	}
	if len(a) > 0 {
		txt = fmt.Sprintf(txt, a...)
	}
	b := logKernel("DEBUG  ", []byte(txt), Control.Print)
	if Logger != nil {
		_, _ = Logger.Write(b)
	}
}

func WARN(txt string, a ...interface{}) {
	if !Control.W {
		return
	}
	if len(a) > 0 {
		txt = fmt.Sprintf(txt, a...)
	}
	b := logKernel("WARN   ", []byte(txt), Control.Print)
	if Logger != nil {
		_, _ = Logger.Write(b)
	}
}

func INFO(txt string, a ...interface{}) {
	if !Control.I {
		return
	}
	if len(a) > 0 {
		txt = fmt.Sprintf(txt, a...)
	}
	b := logKernel("INFO   ", []byte(txt), Control.Print)
	if Logger != nil {
		_, _ = Logger.Write(b)
	}
}

func STATE(txt string, a ...interface{}) {
	if !Control.S {
		return
	}
	if len(a) > 0 {
		txt = fmt.Sprintf(txt, a...)
	}
	b := logKernel("", []byte(txt), Control.Print)
	if Logger != nil {
		_, _ = Logger.Write(b)
	}
}

func RAW_DEBUG(txt string, a ...interface{}) {
	if !Control.DR {
		return
	}
	if len(a) > 0 {
		txt = fmt.Sprintf(txt, a...)
	}
	buf := append([]byte(txt), '\n')
	_, _ = os.Stdout.Write(buf)

	if Logger != nil {
		_, _ = Logger.Write(buf)
	}
}

func logKernel(tag string, data []byte, output bool) []byte {

	ts := time.Now().Local().Format(logKernelTimeLayout)
	pc, file, line, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc).Name()

	// Get code file.
	if len(file) > 0 {
		if index := strings.LastIndexByte(file, '/'); index != -1 {
			file = file[index+1:]
		}
	}

	// Get log prefix.
	if len(tag) > 0 {
		file = fmt.Sprintf("%s%s (%s) %s:%d ", tag, ts, fn, file, line)
	} else {
		file = fmt.Sprintf("STATE  %s ", ts)
	}

	// Write buffer.
	buf := make([]byte, 0, 1024)
	if len(file) > 0 {
		buf = append(buf, file...)
	}
	buf = append(buf, data...)
	buf = append(buf, '\n')

	// Print to console.
	if output {
		_, _ = os.Stdout.Write(buf)
	}

	return buf
}
