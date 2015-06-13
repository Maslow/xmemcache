package hasher

import (
	"crypto/md5"
)

func getHash(data string) []byte {
	h := md5.New()
	h.Write([]byte(data))
	return h.Sum(nil)
}

func GetHashValue(data string) (value uint32) {
	digest := getHash(data)
	for i := 0; i < 4; i++ {
		var t uint32 = 0
		for j := 0; j < 4; j++ {
			t = t | ((uint32)(digest[i*4+j]&0xFF))<<uint8(j*8)
		}
		value += t
	}
	return
}
