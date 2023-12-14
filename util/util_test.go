package util

import "testing"

func TestFormatIndent(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive testing",
			args: args{
				message: "test message",
			},
			want: "  test message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatIndent(tt.args.message); got != tt.want {
				t.Errorf("FormatIndent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintlnWithBlankLineBelow(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive testing",
			args: args{
				message: "test message",
			},
			want: "test message\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintlnWithBlankLineBelow(tt.args.message)
		})
	}
}

func TestPrintlnWithBlankLineAbove(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive testing",
			args: args{
				message: "test message",
			},
			want: "\ntest message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintlnWithBlankLineAbove(tt.args.message)
		})
	}
}

func TestPrintBetweenBlankLine(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive testing",
			args: args{
				message: "test message",
			},
			want: "\ntest message\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintBetweenBlankLine(tt.args.message)
		})
	}
}
