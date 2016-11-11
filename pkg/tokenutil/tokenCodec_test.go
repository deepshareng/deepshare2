package tokenutil

import (
	"testing"
)

func TestCodec(t *testing.T) {
	tests := []struct {
		num    int64
		encode string
	}{
		{1, "1"},
		{16, "g"},
		{15, "f"},
		{63, "11"},
	}

	for i, tt := range tests {
		r := Encode(tt.num)
		if r != tt.encode {
			t.Errorf("#%d: encode(%d)=%s, want %s", i, tt.num, r, tt.encode)
		}
		num := Decode(r)
		if num != tt.num {
			t.Errorf("#%d: decode(%s)=%d, want %d", i, r, num, tt.num)
		}
	}
}
