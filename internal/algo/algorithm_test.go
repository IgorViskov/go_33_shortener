package algo

import "testing"

func TestDecode(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  uint64
	}{
		{
			name:  "DefaultBase",
			token: "e9a",
			want:  19158,
		},
		{
			name:  "DefaultNulStart",
			token: "0001",
			want:  12596221,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decode(tt.token); got != tt.want {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	type args struct {
		n uint64
	}
	tests := []struct {
		name string
		num  uint64
		want string
	}{
		{
			name: "DefaultBase",
			num:  0,
			want: "a",
		},
		{
			name: "DefaultNulStart",
			num:  125,
			want: "cb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Encode(tt.num); got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
