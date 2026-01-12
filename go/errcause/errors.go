// Source code file, created by Developer@YAN_YING_SONG.

package errcause

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

const DateFormat = "2006-01-02 15:04:05.999 MST"
const deferErrorTag = "deferErr/ "

var st bool

var (
	HiddenPathSDK1 string
	HiddenPathSDK2 string
	HiddenPathHOME string
)

func write(err interface{}) {
	date := time.Now().Local().Format(DateFormat)
	stack := fmt.Sprintf("[  ERROR  ] %s\npanic: %s\n%s\n\n", date, err, debug.Stack())

	if len(HiddenPathSDK1) > 0 {
		stack = strings.ReplaceAll(stack, HiddenPathSDK1, "/")
	}
	if len(HiddenPathSDK2) > 0 {
		stack = strings.ReplaceAll(stack, HiddenPathSDK2, "/")
	}
	if len(HiddenPathHOME) > 0 {
		stack = strings.ReplaceAll(stack, HiddenPathHOME, "~/")
	}

	fmt.Println(stack)
	f, _ := os.OpenFile("panic."+date[:10]+".log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0o666)
	_, _ = f.WriteString(stack)
	_ = f.Close()
}

func Recover() {
	// Discard this thread(goroutine) to keep the main process working.

	if err := recover(); err != nil && !st {

		// Current limit, allow once per second.
		st = true
		go func() {
			time.Sleep(time.Second)
			st = false
		}()

		write(err)
	}
}

func LinkErr(err error, fn string) error {
	if err != nil {
		return fmt.Errorf("%s/ %w", fn, err)
	}
	return nil
}

func DeferErr(deferErr *error, err error) {
	if err = LinkErr(err, deferErrorTag); err != nil {
		*deferErr = err
	}
}

func Defer(err error) bool {
	return strings.HasPrefix(err.Error(), deferErrorTag)
}
