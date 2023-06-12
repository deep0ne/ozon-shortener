package base63

import "testing"

func TestEncodeDecode(t *testing.T) {
	testCases := []struct {
		RandomInteger int64
	}{
		{3141312313},
		{32155477471324},
		{41122},
		{123},
		{1},
		{3215423},
		{222111122},
		{101010101},
		{12348129310234218},
		{0},
	}

	for _, tt := range testCases {
		encoded := Encode(tt.RandomInteger)
		if len(encoded) != 10 {
			t.Errorf("Unexpected length of string %v. Must be 10, got %v", encoded, len(encoded))
		}
		decoded, err := Decode(encoded)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		if decoded != tt.RandomInteger {
			t.Errorf("Error in decoding. Expected %v, got %v", tt.RandomInteger, decoded)
		}
	}
}
