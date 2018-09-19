
package algorithm

import (
	"errors"
)

type des2 struct{}

func New2DesEnDecrypter() EnDecrypter {
	return &des2{}
}

func (this *des2) Encrypt(key, data []byte) ([]byte, error) {
	if len(key) != 16 {
		return nil, errors.New("error key")
	}
	//1:取密钥前8个字节数据采用DES加密
	front8Key := key[:8]
	en := NewDesEnDecrypter()
	firstRes, err := en.Encrypt(front8Key, data)
	if err != nil {
		return nil, err
	}
	//2:取密钥后8个字节采用DES解密
	back8Key := key[8:]
	secondRes, err := en.Decrypt(back8Key, firstRes)
	if err != nil {
		return nil, err
	}
	//3:再次用前8个字节的密钥采用DES加密
	last, err := en.Encrypt(front8Key, secondRes)
	if err != nil {
		return nil, err
	}
	return last, nil
}

func (this *des2) Decrypt(key, data []byte) ([]byte, error) {
	if len(key) != 16 {
		return nil, errors.New("error key")
	}
	//1:取密钥前8个字节数据采用DES 解密
	front8Key := key[:8]
	en := NewDesEnDecrypter()
	firstRes, err := en.Decrypt(front8Key, data)
	if err != nil {
		return nil, err
	}
	//2:取密钥后8个字节采用DES加密
	back8Key := key[8:]
	secondRes, err := en.Encrypt(back8Key, firstRes)
	if err != nil {
		return nil, err
	}
	//3:再次用前8个字节的密钥采用DES解密
	last, err := en.Decrypt(front8Key, secondRes)
	if err != nil {
		return nil, err
	}
	return last, nil
}
