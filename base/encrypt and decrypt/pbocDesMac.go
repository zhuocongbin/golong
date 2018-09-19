
package algorithm

import (
	"bytes"
)

type pbocDesMac struct{}

var (
	DefaultInitialVector = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

func NewPbocDesCalculateMAC() ComputeMACer {
	return &pbocDesMac{}
}

func (this *pbocDesMac) CalculateMAC(key, data, iv []byte) ([]byte, error) {
	multiLen := len(data)/8 + 1
	multiByteLen := multiLen * 8
	data = append(data, 0x80)
	tailLen := multiByteLen - len(data)
	data = append(data, bytes.Repeat([]byte{0x00}, tailLen)...)
	//initialize vector
	getPreEightBtyes := data[:8]
	resIV := make([]byte, 8)
	resEncryptIV := make([]byte, 8)
	var err error
	for i := 0; i < 8; i++ {
		resIV[i] = getPreEightBtyes[i] ^ iv[i]
	}
	en := NewDesEnDecrypter()
	resEncryptIV, err = en.Encrypt(key, resIV)
	if err != nil {
		return nil, err
	}
	//encrypt the rest data
	for i := 1; i < multiLen; i++ {
		b2 := data[i*8 : (i+1)*8]
		for byteIndex := 0; byteIndex < 8; byteIndex++ {
			resIV[byteIndex] = b2[byteIndex] ^ resEncryptIV[byteIndex]
		}
		resEncryptIV, err = en.Encrypt(key, resIV)
		if err != nil {
			return nil, err
		}
	}
	return resEncryptIV[:len(key)], nil
}
