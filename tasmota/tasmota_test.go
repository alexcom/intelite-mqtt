package tasmota

import "testing"

func TestEncode(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  string
	}{
		{"test1", []int{8570, 4240, 550, 1580, 550, 510, 565, 1565, 565, 505, 565, 505}, "+8570-4240+550-1580C-510+565-1565F-505FH"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compact(tt.input); got != tt.want {
				t.Errorf("Compact() = %v, want %v", got, tt.want)
			}
		})
	}
}
