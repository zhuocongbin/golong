
package algorithm

import (
	"bytes"
)

type ansix99mac struct{}

func NewANSIx99CalculateMAC() ComputeMACer {
	return &ansix99mac{}
}

//CalculateMAC 采用ANSI X9.9标准算法
func (this *ansix99mac) CalculateMAC(key, data, iv []byte) ([]byte, error) {
	dataLen := len(data)
	bn := dataLen % 8
	if bn != 0 {
		data = append(data, bytes.Repeat([]byte{0x00}, bn)...)
	}
	var (
		resIV []byte = make([]byte, 8)
		err   error
	)
	des1 := NewDesEnDecrypter()
	copy(resIV, iv)
	lastDataLen := len(data) / 8
	for index := 0; index < lastDataLen; index++ {
		resIV, err = Xor(data[index*8:(index+1)*8], resIV)
		if err != nil {
			return nil, err
		}
		resIV, err = des1.Encrypt(key, resIV)
		if err != nil {
			return nil, err
		}
	}
	return resIV, nil
}
