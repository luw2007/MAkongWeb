package main

import (
	"aes"
	"encoding/base64"
	"fmt"
)

var key = []byte("rBwj1MIAivVN222b")

func testAes() {
	key := []byte("rBwj1MIAivVN222b")
	result, err := aes.AesEncrypt([]byte("0"), key)
	if err != nil {
		panic(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(result))
	if base64.StdEncoding.EncodeToString(result) != "NzgOGTK08BvkZN5q8XvG6Q==" {
		panic("不匹配")
	}
	origData, err := aes.AesDecrypt(result, key)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(origData))
}
func main() {
	testAes()
}
