package main

import (
	"golang.org/x/crypto/ripemd160"
	"fmt"
	"encoding/hex"
)

func main() {
	hash := ripemd160.New()
	hash.Write([]byte("hello world"))
	fmt.Println(hex.EncodeToString(hash.Sum(nil)))//16进制打印
	//98c615784ccb5fe5936fbc0cbe9dfdb408d92f0f
	//sha256 --> 256bit  --> 32字节 --> 64个16进制

}