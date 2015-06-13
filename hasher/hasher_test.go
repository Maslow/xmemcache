package hasher

import (
	"encoding/hex"
	"testing"
)

func Test_getHash(t *testing.T) {
	right := "e10adc3949ba59abbe56e057f20f883e"
	value := getHash("123456")
	if right != hex.EncodeToString(value) {
		t.Errorf("MD5 of (kissme) FAILED : %x", value)
	}
}

func Test_GetHashValue(t *testing.T) {
	v := GetHashValue("123456")
	if v == 0 {
		t.Error("Error")
	}
}
