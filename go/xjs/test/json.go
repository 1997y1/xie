package main

import (
	"fmt"
	"time"

	"xie/go/xjs"
)

func doBiz(i int, show bool) {
	set := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
		"field3": "<h1>html</h1>",
		"field4": i,
		"field5": map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
			"field3": "<h1>html</h1>",
			"field4": i,
		},
	}

	// 序列化
	// b := xjs.ToJsonBytes(set)
	b := xjs.ToJsonBytesFast(set)

	// 反序列化
	data := make(map[string]interface{})
	// err := json.Unmarshal(b, &data)
	err := xjs.Unmarshal(b, &data)
	if err != nil {
		panic(err)
	}

	if show {
		fmt.Printf("%s - %v\n", b, data)
	}
}

func testSpeed() {
	now := time.Now()
	for i := 0; i < 10*10000; i++ {
		doBiz(i, false)
	}
	fmt.Println(time.Now().Sub(now))
}

func main() {
	now := time.Now()
	for i := 0; i < 10; i++ {
		doBiz(i, true)
	}
	for i := 0; i < 10; i++ {
		testSpeed()
	}
	fmt.Println(time.Now().Sub(now))
}

/*
2.841850458s ========== ========== ========
1.779047792s ========== =======

*/
