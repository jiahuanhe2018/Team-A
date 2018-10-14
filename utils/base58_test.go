package utils

import "testing"

var tests = []struct {
	raw    string
	encode string
}{
	{"hello world", "StV1DL6CwTryKyV"},
	{"Hello 世界！你好，World！", "5sc47iufhKaZwPBRnLsuXqFUywhevuTPJY7BSCdEABKE"},
}

func TestBase58Encode(t *testing.T) {
	for _, test := range tests {
		gotEncode := Base58Encode([]byte(test.raw))
		if string(gotEncode) != test.encode {
			t.Errorf("test encode %s: got %s, want %s", test.raw, gotEncode, test.encode)
			continue
		}
	}
}

func TestBase58Decode(t *testing.T) {
	for _, test := range tests {
		gotDecode := Base58Decode([]byte(test.encode))
		if string(gotDecode) != test.raw {
			t.Errorf("test decode %s: got %s, want %s", test.encode, gotDecode, test.raw)
			continue
		}
	}
}
