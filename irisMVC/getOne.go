package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var sum int64 = 0

//库存控制
var productNum int64 = 100000

//互斥锁
var mutex sync.Mutex

//卖出计数
var count int64

func GetOneProduct() bool {
	//lock
	mutex.Lock()
	defer mutex.Unlock()
	count += 1
	//check if over ordering
	//release 1 for every 100 buys
	if count%100 == 0 {
		if sum < productNum {
			sum += 1
			fmt.Println(sum)
			return true
		}
	}

	return false
}

func GetProduct(w http.ResponseWriter, req *http.Request) {
	if GetOneProduct() {
		w.Write([]byte("true"))
		return
	}
	w.Write([]byte("false"))
	return
}

func main() {
	http.HandleFunc("/getOne", GetProduct)
	err := http.ListenAndServe(":8084", nil)
	if err != nil {
		log.Fatal("Err:", err)
	}
}
