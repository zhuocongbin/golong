
package algorithm

//EnDecrypter support only one operate form
//DES algorithm support ECB  encrypt mode
type EnDecrypter interface {
	Encrypt(key, data []byte) ([]byte, error)
	Decrypt(key, data []byte) ([]byte, error)
}

type Disperser interface {
	OnceDisperse(key, data []byte) ([]byte, error)
	TwoTimesDisperse(key, data []byte) ([]byte, error)
	ThreeTimesDisperse(key, data []byte) ([]byte, error)
}

type ComputeMACer interface {
	//CalculateMAC
	//iv : initialize vector,default 0x0000000000000000
	CalculateMAC(key, data, iv []byte) ([]byte, error)
}
