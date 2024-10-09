package utility

import (
	"github.com/yanosea/spotlike/app/proxy/fmt"
	"github.com/yanosea/spotlike/app/proxy/io"
	"github.com/yanosea/spotlike/app/proxy/os"
)

// UtilityInterface is an interface for Utility.
type UtilityInterface interface {
	FormatIndent(m string) string
	PrintlnWithWriter(w ioproxy.WriterInstanceInterface, a ...any)
	PrintWithWriterWithBlankLineBelow(w ioproxy.WriterInstanceInterface, a ...any)
	PrintWithWriterWithBlankLineAbove(w ioproxy.WriterInstanceInterface, a ...any)
	PrintWithWriterBetweenBlankLine(w ioproxy.WriterInstanceInterface, a ...any)
}

// Utility is a struct that implements UtilityInterface.
type Utility struct {
	FmtProxy fmtproxy.Fmt
	OsProxy  osproxy.Os
}

// New is a constructor for Utility.
func New(
	fmtProxy fmtproxy.Fmt,
	osProxy osproxy.Os,
) *Utility {
	return &Utility{
		FmtProxy: fmtProxy,
		OsProxy:  osProxy,
	}
}

func (u *Utility) FormatIndent(m string) string {
	return "  " + m
}

func (u *Utility) PrintlnWithWriter(w ioproxy.WriterInstanceInterface, a ...any) {
	u.FmtProxy.Fprintf(w, u.FmtProxy.Sprintf("%s", a[0])+"\n")
}

func (u *Utility) PrintWithWriterWithBlankLineBelow(w ioproxy.WriterInstanceInterface, a ...any) {
	u.FmtProxy.Fprintf(w, u.FmtProxy.Sprintf("%s\n", a[0])+"\n")
}

func (u *Utility) PrintWithWriterWithBlankLineAbove(w ioproxy.WriterInstanceInterface, a ...any) {
	u.FmtProxy.Fprintf(w, u.FmtProxy.Sprintf("\n%s", a[0])+"\n")
}

func (u *Utility) PrintWithWriterBetweenBlankLine(w ioproxy.WriterInstanceInterface, a ...any) {
	u.FmtProxy.Fprintf(w, u.FmtProxy.Sprintf("\n%s\n", a[0])+"\n")
}
