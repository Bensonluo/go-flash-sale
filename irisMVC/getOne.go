package main

import (
	"log"
	"net/http"
	"sync"
)

var sum int64 = 0

//库存控制
var productNum int64 = 10000

//互斥锁
var mutex sync.Mutex

func GetOneProduct() bool {
	//lock
	mutex.Lock()
	defer mutex.Unlock()
	//check if overload
	if sum < productNum {
		sum += 1
		return true
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