// Source code file, created by Developer@YAN_YING_SONG.

package console

import (
	"errors"
	"testing"
	"time"

	"github.com/pkg/profile"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 基础测试
func Test(t *testing.T) {
	Logger = &lumberjack.Logger{
		Filename:   "console.log",
		MaxSize:    100,
		MaxBackups: 10,
		LocalTime:  true,
		Compress:   false,
	}

	txt := "error message."

	ERROR(errors.New(txt))
	DEBUG(txt)
	WARN(txt)
	INFO(txt)
	STATE(txt)
	RAW_DEBUG(txt)
}

// 性能测试，每秒可处理的日志量。
/**

  # 日志记录模式结果对比

  - 测试环境A：固态硬盘：2000MB/s

  1. default               1.457s 279879
  2. use channel           1.909s 340256
  3. use channel+bucket    1.501s 1202581

  - 测试环境B：低速硬盘环境：200MB/s

  1. default               1.115s 11339
  2. use channel          12.822s 111137
  3. use channel+bucket    2.522s 173910

*/
func Test2(t *testing.T) {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	Logger = &lumberjack.Logger{
		Filename:   "console.log",
		MaxSize:    100,
		MaxBackups: 3,
		LocalTime:  true,
		Compress:   false,
	}

	txt := "error message."
	go func() {
		for {
			INFO(txt)
		}
	}()

	time.Sleep(time.Second)
}

func Test3(t *testing.T) {
	Logger = &lumberjack.Logger{
		Filename:   "console.log",
		MaxSize:    100,
		MaxBackups: 10,
		LocalTime:  true,
		Compress:   false,
	}

	txt := "error message."

	for {
		INFO(txt)
		time.Sleep(1000 * time.Millisecond)
	}
}
