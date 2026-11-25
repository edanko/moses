package scan

import "testing"

func TestReadNextLine(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []string
	}{
		{
			name: "preserves case and spaces and strips CR",
			text: " Hello World \r\nSecond\n",
			want: []string{" Hello World ", "Second", eof},
		},
		{
			name: "blank lines stay empty",
			text: "\n\nlast",
			want: []string{"", "", "last", eof},
		},
		{
			name: "missing trailing newline",
			text: "no-newline",
			want: []string{"no-newline", eof},
		},
		{
			name: "empty input",
			text: "",
			want: []string{eof},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.text)
			for i, want := range tt.want {
				if got := s.ReadNextLine(); got != want {
					t.Fatalf("line %d = %q, want %q", i, got, want)
				}
			}
		})
	}
}

func TestReadLine(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []string
	}{
		{
			name: "lowercases and maps negative zero",
			text: " -0\r\nABC\n",
			want: []string{"0", "abc", ""},
		},
		{
			name: "blank line and eof",
			text: "\n\n",
			want: []string{"", "", ""},
		},
		{
			name: "missing trailing newline",
			text: "Some Line",
			want: []string{"some line", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.text)
			for i, want := range tt.want {
				if got := s.ReadLine(); got != want {
					t.Fatalf("line %d = %q, want %q", i, got, want)
				}
			}
		})
	}
}
