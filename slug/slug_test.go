package slug

import "testing"

func TestMake(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{""},
			want: "",
		},
		{
			name: "simple",
			args: args{"hello world"},
			want: "hello-world",
		},
		{
			name: "a lot of spaces",
			args: args{"hello    world"},
			want: "hello-world",
		},
		{
			name: "a lot of dashes",
			args: args{"hello----world"},
			want: "hello-world",
		},
		{
			name: "a lot of backslashes",
			args: args{"hello\\\\\\\\\\\\\\\\world"},
			want: "hello-world",
		},
		{
			name: "a lot of slashes",
			args: args{"hello////////world"},
			want: "hello-world",
		},
		{
			name: "special chars",
			args: args{"hello world!"},
			want: "hello-world",
		},
		{
			name: "control chars",
			args: args{"hello\nworld"},
			want: "hello-world",
		},
		{
			name: "unicode",
			args: args{"hello 世界"},
			want: "hello-世界",
		},
		{
			name: "unicode emoji",
			args: args{"hello 🌍"},
			want: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Make(tt.args.s); got != tt.want {
				t.Errorf("Make() = %v, want %v", got, tt.want)
			}
		})
	}
}
