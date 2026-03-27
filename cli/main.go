package main

import (
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, out io.Writer, errOut io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(errOut, "missing command")
		printUsage(errOut)
		return 2
	}

	switch args[0] {
	case "done":
		return cmdDone(args[1:], out, errOut)
	case "help", "-h", "--help":
		printUsage(out)
		return 0
	default:
		fmt.Fprintf(errOut, "unknown command %q\n", args[0])
		printUsage(errOut)
		return 2
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  tm done <id>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  done    Mark task as done by id")
}

func cmdDone(args []string, _ io.Writer, errOut io.Writer) int {
	if len(args) == 0 || args[0] == "" {
		fmt.Fprintln(errOut, "missing id for \"done\"")
		fmt.Fprintln(errOut, "usage: tm done <id>")
		return 2
	}
	id := args[0]
	if _, err := uuid.Parse(id); err != nil {
		fmt.Fprintf(errOut, "invalid id %q: must be a UUID\n", id)
		return 2
	}

	// Network behavior not implemented in this bead; this command focuses on UX for argument validation.
	fmt.Fprintln(errOut, "not implemented: backend request")
	return 1
}
