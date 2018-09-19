
package algorithm

import (
	"fmt"
	"testing"
)

func TestAES(t *testing.T) {
	en := NewAESEnDecrypter()
	d := Test{
		Key:  []byte(TopSecretKey),
		Data: []byte("0123456789ABCDEFFEDCBA98765432100123456789ABCDEF"),
		Want: []byte("0123456789ABCDEFFEDCBA98765432100123456789ABCDEF"),
	}
	res, err := en.Encrypt(d.Key, d.Data)
	if err != nil {
		t.Error("failed:" + err.Error())
	}
	// fmt.Printf("%X\n", res)
	// fmt.Printf("TopKey:%X\n", TopSecretKey)
	out, err := en.Decrypt(d.Key, res)
	if err != nil {
		t.Error("failed:" + err.Error())
	}
	src := fmt.Sprintf("%X", out)
	dst := fmt.Sprintf("%X", d.Want)
	if src != dst {
		t.Errorf("falied \nsrc(%s)\ndst(%s) \n", src, dst)
	}
}
