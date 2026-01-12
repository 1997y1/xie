package main

import (
	"fmt"
	"time"

	"xie"
)

const N = 10000 * 10000

var bRaw = []byte("" +
	"Build simple, secure, scalable systems with Go\nAn open-source programming language supported by Google\nEasy to learn and great for teams\nBuilt-in concurrency and a robust standard library\nLarge ecosystem of partners, communities, and tools\n" +
	"使用 Go 构建简单、安全、可扩展的系统\nGoogle 支持的开源编程语言\n易于学习，非常适合团队\n内置并发性和强大的标准库\n庞大的合作伙伴、社区和工具生态系统\n" +
	"")

var sRaw = "" +
	"Build simple, secure, scalable systems with Go\nAn open-source programming language supported by Google\nEasy to learn and great for teams\nBuilt-in concurrency and a robust standard library\nLarge ecosystem of partners, communities, and tools\n" +
	"使用 Go 构建简单、安全、可扩展的系统\nGoogle 支持的开源编程语言\n易于学习，非常适合团队\n内置并发性和强大的标准库\n庞大的合作伙伴、社区和工具生态系统\n" +
	""

func testToString() {
	start := time.Now()
	for i := 0; i < N; i++ {
		_ = string(bRaw)
	}
	fmt.Println("string() =>", time.Since(start))

	start = time.Now()
	for i := 0; i < N; i++ {
		_ = xie.Ts(bRaw)
	}
	fmt.Println("    Ts() =>", time.Since(start))
}

func testToBytes() {
	start := time.Now()
	for i := 0; i < N; i++ {
		_ = []byte(sRaw)
	}
	fmt.Println("[]byte() =>", time.Since(start))
}

func main() {
	testToString()
	testToBytes()
}

/*

string() => 4.165599667s
    Bs() => 31.399709ms
    Ts() => 25.0705ms
[]byte() => 24.956791ms
    Sb() => 25.011625ms

*/
