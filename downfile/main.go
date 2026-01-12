package main

import (
	"fmt"
	"math"
	"net"
	"net/http"
	"runtime"
	"time"
	"xie/go/xjs"
)

func eatCpu() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
		loop:
			r := math.Pi*2 + 1
			_ = r
			goto loop
		}()
	}
	time.Sleep(30 * time.Second)
}

func main() {
	eatCpu()
	return

	port := ":5555"
	Is, err := NetInterfaces()
	if err != nil {
		fmt.Println(err)
		// return
	}
	fmt.Printf("== Downfile Server - Port %s ==\n%s\n", port, xjs.ToJsonBytes(Is, true))

	err = http.ListenAndServe(port, http.FileServer(http.Dir("/Users/tmp1/Downloads")))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func NetInterfaces() (map[string][]string, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	result := make(map[string][]string, 10)

	for i := 0; i < len(ifs); i++ {
		if len(ifs[i].HardwareAddr) == 0 {
			continue
		}
		list, _ := ifs[i].Addrs()
		ls := make([]string, 0, len(list))
		for j := 0; j < len(list); j++ {
			ls = append(ls, list[j].String())
		}
		if len(ls) > 0 {
			result[ifs[i].HardwareAddr.String()] = ls
		}
	}

	return result, nil
}

/*
Default: 33%
CPU Max: 43% - 10w
GPU Max: 63% - 30w
Media engine - ?w

*/
