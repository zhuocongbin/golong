package main

import (
	"fmt"
	"crypto/md5"
	"encoding/hex"
)

//md5加密
func main() {
	data := []byte("hello world")
	sum := md5.Sum(data)
	fmt.Printf("%x\n",sum)//16进制打印
	//5eb63bbbe01eeed093cb22bb8f5acdc3

	data1 := []byte("hello world 1111111111111111111")
	hash := md5.New()
	hash.Write(data1)
	bytes := hash.Sum(nil)
	fmt.Println(hex.EncodeToString(bytes))// []byte转为16进制
	//07a1e89a32376deda92c5ed365f5d3d6

	//128bit --> 16字节  -->32个16进制
}