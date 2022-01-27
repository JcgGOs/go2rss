package model

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"config_test", "./config_test.json", "socks5://127.0.0.1:5000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.path)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}

			if got.Proxy != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Parse() = %v, want %v", got, tt.want)
			// }
		})
	}
}
