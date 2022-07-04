package main

import (
	"strings"

	"github.org/tiedotguy/scrubeer/internal/args"
)

type Opts struct {
	InputFile  string   `short:"i" long:"input-file" required:"true"`
	OutputFile string   `short:"o" long:"output-file" required:"true"`
	Keep       []string `short:"k" long:"keep"`
	keep       map[int]struct{}
}

var _ = args.Validater((*Opts)(nil))

func (o *Opts) Validate() []string {
	var keeps []string
	for _, k := range o.Keep {
		keeps = append(keeps, strings.Split(k, ",")...)
	}
	o.keep = make(map[int]struct{})

	var errors []string
	for _, keep := range keeps {
		switch keep {
		case "bool", "boolean":
			o.keep[4] = struct{}{}
		case "char":
			o.keep[5] = struct{}{}
		case "float":
			o.keep[6] = struct{}{}
		case "double":
			o.keep[7] = struct{}{}
		case "byte":
			o.keep[8] = struct{}{}
		case "short":
			o.keep[9] = struct{}{}
		case "int":
			o.keep[10] = struct{}{}
		case "long":
			o.keep[11] = struct{}{}
		default:
			errors = append(errors, "unrecognized type: "+keep+", allowed: bool[ean], char, float, double, byte, short, int, long")
		}
	}
	return errors
}
