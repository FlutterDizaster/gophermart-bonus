package api

import "testing"

func Test_checkLuhn(t *testing.T) {
	type test struct {
		name   string
		number string
		want   bool
	}
	tests := []test{
		{
			name:   "true test",
			number: "4561261212345467",
			want:   true,
		},
		{
			name:   "false test",
			number: "4561261212345464",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkLuhn(tt.number); got != tt.want {
				t.Errorf("checkLuhn() = %v, want %v", got, tt.want)
			}
		})
	}
}
