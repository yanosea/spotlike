package util

import (
	"io"
	"os"
	"testing"
)

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
		writer  io.Writer
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive testing",
			args: args{
				message: "test message stdout",
				writer:  os.Stdout,
			},
			want: "test message\n",
		}, {
			name: "positive testing",
			args: args{
				message: "test message stderr",
				writer:  os.Stderr,
			},
			want: "test message\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintWithWriterWithBlankLineBelow(tt.args.writer, tt.args.message)
		})
	}
}

func TestPrintlnWithBlankLineAbove(t *testing.T) {
	type args struct {
		message string
		writer  io.Writer
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive testing",
			args: args{
				message: "test message stdout",
				writer:  os.Stdout,
			},
			want: "\ntest message",
		}, {
			name: "positive testing",
			args: args{
				message: "test message stderr",
				writer:  os.Stderr,
			},
			want: "\ntest message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintWithWriterWithBlankLineAbove(tt.args.writer, tt.args.message)
		})
	}
}

func TestPrintBetweenBlankLine(t *testing.T) {
	type args struct {
		message string
		writer  io.Writer
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive testing",
			args: args{
				message: "test message stdout",
				writer:  os.Stdout,
			},
			want: "\ntest message\n",
		}, {
			name: "positive testing",
			args: args{
				message: "test message stderr",
				writer:  os.Stderr,
			},
			want: "\ntest message\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintWithWriterBetweenBlankLine(tt.args.writer, tt.args.message)
		})
	}
}
