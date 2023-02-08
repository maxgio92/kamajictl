package output

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"io"
	"k8s.io/kubectl/pkg/cmd/util"
	"os"
	"strings"
)

// Printer represents a helper structure for pretty printing of all the tool's output.
type Printer struct {
	Info     *pterm.PrefixPrinter
	Success  *pterm.PrefixPrinter
	Warning  *pterm.PrefixPrinter
	Error    *pterm.PrefixPrinter
	Action   *pterm.PrefixPrinter
	Generate *pterm.PrefixPrinter
	Waiting  *pterm.PrefixPrinter

	TablePrinter *pterm.TablePrinter

	writer io.Writer
}

type PrintOption func(*Printer)

// NewPrinter returns a new Printer struct.
func NewPrinter(opts ...PrintOption) *Printer {
	generic := &pterm.PrefixPrinter{MessageStyle: pterm.NewStyle(pterm.FgDefault)}

	tablePrinter := pterm.DefaultTable.WithHasHeader().WithBoxed()

	p := &Printer{
		Info: generic.WithPrefix(pterm.Prefix{
			Text:  PrefixInfo,
			Style: pterm.NewStyle(pterm.FgDefault),
		}),

		Success: generic.WithPrefix(pterm.Prefix{
			Text:  PrefixSuccess,
			Style: pterm.NewStyle(pterm.FgLightGreen),
		}),

		Warning: generic.WithPrefix(pterm.Prefix{
			Text:  PrefixWarning,
			Style: pterm.NewStyle(pterm.FgYellow),
		}),

		Error: generic.WithPrefix(pterm.Prefix{
			Text:  PrefixError,
			Style: pterm.NewStyle(pterm.FgRed),
		}),

		Action: generic.WithPrefix(pterm.Prefix{
			Text:  PrefixAction,
			Style: pterm.NewStyle(pterm.FgDefault),
		}),

		Generate: generic.WithPrefix(pterm.Prefix{
			Text:  PrefixGenerate,
			Style: pterm.NewStyle(pterm.FgDefault),
		}),

		Waiting: generic.WithPrefix(pterm.Prefix{
			Text:  PrefixWaiting,
			Style: pterm.NewStyle(pterm.FgDefault),
		}),

		TablePrinter: tablePrinter,
	}

	for _, f := range opts {
		f(p)
	}

	return p
}

func WithWriter(w io.Writer) PrintOption {
	return func(opts *Printer) {
		opts.writer = w
		opts.Info = opts.Info.WithWriter(w)
		opts.Warning = opts.Warning.WithWriter(w)
		opts.Error = opts.Error.WithWriter(w)
		opts.Success = opts.Success.WithWriter(w)
		opts.TablePrinter = opts.TablePrinter.WithWriter(w)
	}
}

func (p *Printer) validate() error {
	if p.writer == nil {
		return errors.New("printer writer cannot be nil")
	}

	return nil
}

// PrintRaw prints the data as string raw to the configured io.Writer.
func (p *Printer) PrintRaw(s string) {
	bufw := bufio.NewWriter(p.writer)
	bufw.WriteString(s)
}

// PrintTable prints a pretty table with the data argument.
func (p *Printer) PrintTable(data [][]string) error {
	return p.TablePrinter.WithData(data).Render()
}

// CheckErr prints a user-friendly error to STDERR and exits with a non-zero exit code.
func (p *Printer) CheckErr(err error) {
	if p != nil {
		util.BehaviorOnFatal(func(msg string, code int) {
			msg = strings.TrimPrefix(msg, "error: ")
			p.Error.Println(strings.TrimRight(msg, "\n"))
		})
	}
	util.CheckErr(err)
}

// ExitOnErr aborts the execution in case of errors, without printing any error message.
func (p *Printer) ExitOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(util.DefaultErrorExitCode)
	}
}

func (p Printer) Infof(format string, a ...interface{}) {
	p.Info.Printfln(format, a...)
}

func (p Printer) Actionf(format string, a ...interface{}) {
	p.Action.Printfln(format, a...)
}

func (p Printer) Generatef(format string, a ...interface{}) {
	p.Generate.Printfln(format, a...)
}

func (p Printer) Waitingf(format string, a ...interface{}) {
	p.Waiting.Printfln(format, a...)
}

func (p Printer) Successf(format string, a ...interface{}) {
	p.Success.Printfln(format, a...)
}

func (p Printer) Warningf(format string, a ...interface{}) {
	p.Warning.Printfln(format, a...)
}

func (p Printer) Failuref(format string, a ...interface{}) {
	p.Error.Printfln(format, a...)
}
