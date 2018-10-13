package utils

import "testing"

var privateKey, publicKey = NewKeyPair()

func TestSign(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"a"},
		{"Hello World"},
		{"你好世界"},
	}
	for _, test := range tests {
		signature := Sign(privateKey, []byte(test.input))
		if !Verify(publicKey, signature, []byte(test.input)) {
			t.Errorf("test input %s: signature %x, verify faild.", test.input, signature)
			continue
		}
	}
}
