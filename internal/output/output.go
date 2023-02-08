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
	Info    *pterm.PrefixPrinter
	Success *pterm.PrefixPrinter
	Warning *pterm.PrefixPrinter
	Error   *pterm.PrefixPrinter

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

func (p *Printer) Raw(s string) {
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
