
package algorithm

import (
	"fmt"
)

type unionPayMac struct{}

func NewUnionPayCalculateMAC() ComputeMACer {
	return &unionPayMac{}
}

func (this *unionPayMac) CalculateMAC(key, data, iv []byte) ([]byte, error) {
	data = ZeroPadding(data, 8)
	var (
		resIV []byte = make([]byte, 8)
		err   error
	)
	des2 := New2DesEnDecrypter()
	copy(resIV, iv)
	//1:循环异或
	lastDataLen := len(data) / 8
	for index := 0; index < lastDataLen; index++ {
		resIV, err = Xor(data[index*8:(index+1)*8], resIV)
		if err != nil {
			return nil, err
		}
	}
	//2:转ASCII码
	resAscii := []byte(fmt.Sprintf("%X", resIV))
	//3:加密左半部分
	leftPart, err := des2.Encrypt(key, resAscii[:8])
	if err != nil {
		return nil, err
	}
	//4:左右异或
	macData, _ := Xor(leftPart, resAscii[8:16])
	//5:加密数据并返回
	resIV, err = des2.Encrypt(key, macData)
	if err != nil {
		return nil, err
	}
	return resIV, nil
}
