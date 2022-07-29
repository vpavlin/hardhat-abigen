package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/vpavlin/hardhat-abigen/internal/extractor"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [FILE|DIRECTORY]\n", os.Args[0])

		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "  --%s\n\t %v (default: '%s')\n", f.Name, f.Usage, f.DefValue) // f.Name, f.Value
		})
	}

	outDir := flag.String("outDir", "output", "Destination directory for extracted ABI json files and generated bindings")
	generateBindings := flag.Bool("bindings", true, "Use Abigen to generate bindings")
	exclude := flag.String("exclude", "", "List of strings to match aginst paths to exclude them from processing")

	flag.Parse()

	inFile := flag.Arg(0)

	info, err := os.Stat(inFile)
	if err != nil {
		logrus.Errorf("No file provided")
		flag.Usage()
		os.Exit(1)
	}

	if info.IsDir() {
		extractor.ProcessAll(inFile, *outDir, *generateBindings, *exclude)
	} else {
		extractor.ProcessSingle(inFile, *outDir, *generateBindings, *exclude)
	}

}
