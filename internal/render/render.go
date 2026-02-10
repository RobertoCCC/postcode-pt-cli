package render

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-isatty"

	"github.com/RobertoCCC/postcode-pt-cli/internal/api"
)

type Options struct {
	JSON     bool
	NoColor  bool
	Out      io.Writer
}

func NewOptions(out io.Writer) Options {
	return Options{Out: out}
}

func (o Options) useColor() bool {
	if o.NoColor {
		return false
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	f, ok := o.Out.(*os.File)
	if !ok {
		return false
	}
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}

const (
	ansiReset = "\x1b[0m"
	ansiBold  = "\x1b[1m"
	ansiDim   = "\x1b[2m"
	ansiCyan  = "\x1b[36m"
)

func (o Options) bold(s string) string {
	if !o.useColor() {
		return s
	}
	return ansiBold + s + ansiReset
}

func (o Options) dim(s string) string {
	if !o.useColor() {
		return s
	}
	return ansiDim + s + ansiReset
}

func (o Options) cyan(s string) string {
	if !o.useColor() {
		return s
	}
	return ansiCyan + s + ansiReset
}

func PostalCodes(o Options, entries []api.PostalCodeEntry) error {
	if o.JSON {
		return writeJSON(o.Out, entries)
	}
	if len(entries) == 0 {
		fmt.Fprintln(o.Out, "No entries found.")
		return nil
	}
	for i, e := range entries {
		if i > 0 {
			fmt.Fprintln(o.Out)
		}
		fmt.Fprintf(o.Out, "%s  %s\n", o.bold(o.cyan(e.Code)), o.bold(e.Designation))
		if street := formatStreet(e.Street); street != "" {
			fmt.Fprintf(o.Out, "  %s %s\n", o.dim("street     "), street)
		}
		fmt.Fprintf(o.Out, "  %s %s (%s)\n", o.dim("locality   "), e.Locality.Name, e.Locality.Code)
		fmt.Fprintf(o.Out, "  %s %s (%s)\n", o.dim("municipality"), e.Municipality.Name, e.Municipality.Code)
		fmt.Fprintf(o.Out, "  %s %s (%s)\n", o.dim("district   "), e.District.Name, e.District.Code)
	}
	return nil
}

func Districts(o Options, districts []api.DistrictBrief) error {
	if o.JSON {
		return writeJSON(o.Out, districts)
	}
	if len(districts) == 0 {
		fmt.Fprintln(o.Out, "No districts found.")
		return nil
	}
	for _, d := range districts {
		fmt.Fprintf(o.Out, "%s  %s\n", o.bold(o.cyan(d.Code)), d.Name)
	}
	return nil
}

func Municipalities(o Options, municipalities []api.MunicipalityWithDistrict) error {
	if o.JSON {
		return writeJSON(o.Out, municipalities)
	}
	if len(municipalities) == 0 {
		fmt.Fprintln(o.Out, "No municipalities found.")
		return nil
	}
	for _, m := range municipalities {
		fmt.Fprintf(o.Out, "%s  %s\n", o.bold(o.cyan(m.Code)), m.Name)
	}
	return nil
}

func formatStreet(s api.Street) string {
	parts := []string{}
	if s.Type != nil && *s.Type != "" {
		parts = append(parts, *s.Type)
	}
	if s.Name != nil && *s.Name != "" {
		parts = append(parts, *s.Name)
	}
	return strings.Join(parts, " ")
}

func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
