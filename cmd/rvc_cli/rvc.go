package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rapidsai/rvc/pkg/rvc"
)

const (
	ProgramName = "rvc"
)

type Args struct {
	UcxPyVersion  string
	RapidsVersion string
}

func initFlags(flagset *flag.FlagSet) *Args {
	args := &Args{}

	flagset.StringVar(&args.RapidsVersion, "rapids", "", "Rapids version")
	flagset.StringVar(&args.UcxPyVersion, "ucx-py", "", "ucx-py version")

	return args
}

func main() {
	flags := flag.NewFlagSet(ProgramName, flag.ExitOnError)
	args := initFlags(flags)

	_ = flags.Parse(os.Args[1:])
	if len(flags.Args()) > 0 {
		fmt.Printf("Unknown command line argument \"%v\"\n", flags.Args()[0])
		flags.Usage()
		os.Exit(1)
	}

	if len(args.RapidsVersion) < 1 && len(args.UcxPyVersion) < 1 {
		flags.Usage()
		os.Exit(1)
	}

	if len(args.RapidsVersion) > 0 && len(args.UcxPyVersion) > 0 {
		fmt.Println("\"-rapids\" and \"-ucx-py\" parameters are mutually exclusives")
		os.Exit(1)
	}

	var version string
	var err error
	if len(args.RapidsVersion) > 1 {
		version, err = rvc.GetUcxPyFromRapids(args.RapidsVersion)
	} else {
		version, err = rvc.GetRapidsFromUcxPy(args.UcxPyVersion)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Print(version)
}
