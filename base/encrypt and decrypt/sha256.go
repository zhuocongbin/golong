package main

import (
	"crypto/sha256"
	"fmt"
	"encoding/hex"
	"os"
	"log"
	"io"
)

func main() {
	//3种用法

	//一
	data1 := []byte("hello world")
	sum256 := sha256.Sum256(data1)
	fmt.Printf("%x\n",sum256)//b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9

	//二
	data2 := []byte("hello world123")
	hash := sha256.New()
	hash.Write(data2)
	sum := hash.Sum(nil)
	fmt.Println(hex.EncodeToString(sum))//e6ec8096bfa71305f98061e024e9b2f67b0df12bbd56e9423e91203968ae2e9d

	//三 加密一个文件中的内容
	file, err := os.Open("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	h := sha256.New()
	io.Copy(h,file)
	bytes := h.Sum(nil)
	fmt.Println(hex.EncodeToString(bytes))//85303176ac15c747a2b49a7fb4844e8e8740b4ea8c81d961a31b76f70af4de48
	//sha256 --> 256bit  --> 32字节 --> 64个16进制

}