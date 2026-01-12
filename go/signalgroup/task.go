// Source code file, created by Developer@YAN_YING_SONG.

// Async work parallel controller based on system signals.

package signalgroup

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"xie"
	"xie/go/errcause"
)

func exitHistory(message os.Signal) {
	fp := xie.JoinFp("exit.history")
	f, _ := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	t := time.Now().Local().Format("20060102.150405")
	s := fmt.Sprintf("%s %s %s\n", t, signals[message.(syscall.Signal)], xie.ExeFp())
	_, _ = f.WriteString(s)
	_ = f.Close()
}

func thread(deferRun time.Duration, task func() error) {
	defer errcause.Recover()
	defer Quit()

	if deferRun > 0 {
		time.Sleep(deferRun)
	}

	if err := task(); err != nil {
		panic(err) // Last panic.
	}
}

func Quit() {
	// Send process exit signal.

	sig <- syscall.SIGQUIT
}

func Async(deferRun time.Duration, task func() error) {
	// New Goroutine.
	//
	// ! If fail one, And fail all.

	// Task +1.
	atomic.AddInt32(&countWork, 1)

	go thread(deferRun, task)
}

func Wait() {
	// Listen process exit signal.

	if countWork == 0 {
		return
	}

	signal.Notify(sig, listens...)
	message := <-sig
	exitHistory(message)
}

var countWork int32

var sig = make(chan os.Signal)

var listens = []os.Signal{
	syscall.SIGHUP,  // 01
	syscall.SIGINT,  // 02
	syscall.SIGQUIT, // 03
	syscall.SIGILL,  // 04
	syscall.SIGTRAP, // 05
	syscall.SIGABRT, // 06
	syscall.SIGBUS,  // 07
	syscall.SIGFPE,  // 08
	syscall.SIGKILL, // 09
	syscall.SIGSEGV, // 11
	// syscall.SIGPIPE, // 13
	syscall.SIGALRM, // 14
	syscall.SIGTERM, // 15
}

var signals = [...]string{

	// Generic signals
	1:  "SIG01 会话挂断、终端控制台窗口会话被关闭",
	2:  "SIG02 程序中断、来自键盘的操作 Ctrl-C",
	3:  "SIG03 程序退出、或键盘触发 Ctrl-\\",
	4:  "SIG04 非法指令、通常是程序文件本身有错误",
	5:  "SIG05 断点跟踪、通常是工程师正在调试代码",
	6:  "SIG06 异常结束、通常是程序自身代码错误",
	7:  "SIG07 总线错误、非法地址，内存地址未对齐",
	8:  "SIG08 浮点异常、浮点数值运算错误",
	9:  "SIG09 强制终止、KILL -9 PID，无法捕获",
	10: "SIG10 用户定义、用户自定义信号1",
	11: "SIG11 分段冲突、内存地址访问越界",
	12: "SIG12 用户定义、用户自定义信号2",
	13: "SIG13 管道破裂、进程间通信故障",
	14: "SIG14 时钟错误、发生实时定时器时钟错误",
	15: "SIG15 程序终止、进程被杀死 KILL PID",

	// UNIX signals
	16: "SIG16 协栈错误、协处理器栈错误",
	17: "SIG17 子程停止、子进程停止",
	18: "SIG18 程序恢复、恢复被暂停执行的进程",
	19: "SIG19 程序停止、程序被暂停执行，无法捕获",
	20: "SIG20 程序停止、或键盘触发 Ctrl-Z",
	21: "SIG21 程序停止、后台进程请求输入",
	22: "SIG22 程序停止、后台进程请求输出",
	23: "SIG23 紧急协议、使用紧急模式传输重要数据",
	24: "SIG24 处理超时、处理时间超过 CPU 时限",
	25: "SIG25 文件过大、超过文件最大限制",
	26: "SIG26 时钟错误、发生虚拟定制器时钟错误",
	27: "SIG27 时钟错误、发生概况定制器时钟错误",
	28: "SIG28 窗口调整、窗口大小被调整",
	29: "SIG29 文件就绪、文件描述符就绪，可以操作",
	30: "SIG30 主机掉电、主机电源供给失效",
	31: "SIG31 系统错误、通常是错误的系统调用",
}
