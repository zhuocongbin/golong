

package algorithm

import (
	"errors"
)

func Xor(done, dtwo []byte) ([]byte, error) {
	if len(done) != len(dtwo) {
		return nil, errors.New("data length should be equal.")
	}
	dlen := len(done) - 1
	for dlen >= 0 {
		done[dlen] ^= dtwo[dlen]
		dlen = dlen - 1
	}
	return done, nil
}
