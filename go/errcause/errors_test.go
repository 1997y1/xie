package errcause

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func doWork() (deferErr error) {
	source, err := os.Open("1.txt")
	if source != nil {
		defer func() { DeferErr(&deferErr, source.Close()) }()
	}
	if err != nil {
		return LinkErr(err, "os.Open")
	}
	_ = source

	return nil
}

func Test_open(t *testing.T) {

	go func() {
		defer Recover()

		if err := doWork(); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)
	fmt.Println("Test end.")
}
