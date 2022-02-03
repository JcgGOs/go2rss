package util

import (
	"testing"
)

func TestNoExt(t *testing.T) {

	tests := []struct {
		name     string
		fileName string
		want     string
	}{
		{"default", "default", "default"},
		{"default.ext", "default.ext", "default"},
		{"default.js.ext", "default.js.ext", "default.js"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NoExt(tt.fileName); got != tt.want {
				t.Errorf("NoExt() = %v, want %v", got, tt.want)
			}
		})
	}
}
