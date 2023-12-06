package exit

import (
	"testing"
)

func TestExitCode_GetNumber(t *testing.T) {
	type args struct {
		exitCode ExitCode
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "positive testing (ok)",
			args: args{
				exitCode: CodeOk,
			},
			want: 0,
		},
		{
			name: "positive testing (error init)",
			args: args{
				exitCode: CodeErrInit,
			},
			want: 1,
		},
		{
			name: "positive testing (error command)",
			args: args{
				exitCode: CodeErrCmd,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.exitCode.GetNumber(); got != tt.want {
				t.Errorf("ExitCode.GetNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExit(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "positive testing",
			want: 0,
		},
	}
	oldExit := osExit
	defer func() { osExit = oldExit }()
	var status int
	exit := func(code int) {
		status = code
	}
	osExit = exit
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Exit()
			if status != tt.want {
				t.Errorf("Exit code: %v, want %v", status, tt.want)
			}
		})
	}
}

func TestErrorExit(t *testing.T) {
	type args struct {
		exitCode ExitCode
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "positive testing (error init)",
			args: args{
				exitCode: CodeErrInit,
			},
			want: 1,
		},
		{
			name: "positive testing (error command)",
			args: args{
				exitCode: CodeErrCmd,
			},
			want: 2,
		},
	}
	oldExit := osExit
	defer func() { osExit = oldExit }()
	var status int
	exit := func(code int) {
		status = code
	}
	osExit = exit
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ErrorExit(tt.args.exitCode)
			if status != tt.want {
				t.Errorf("Exit code: %v, want %v", status, tt.want)
			}
		})
	}
}
