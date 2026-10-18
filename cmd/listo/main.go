package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/khinshankhan/listo/controller"
	"github.com/khinshankhan/listo/platform"
	"github.com/khinshankhan/listo/platform/logx"
	"github.com/khinshankhan/logstox/fields"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, os.Args[1:]); err != nil {
		logx.Default().Error("exiting", fields.Error(err))
		os.Exit(1)
	}
}

// run dispatches to a subcommand.
func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		usage()
		return fmt.Errorf("no subcommand given")
	}

	switch cmd := args[0]; cmd {
	case "serve":
		return controller.RunHTTPServer(ctx)
	case "version":
		fmt.Fprintln(os.Stdout, platform.Version("listo").String())
		return nil
	case "help", "-h", "--help":
		usage()
		return nil
	default:
		usage()
		return fmt.Errorf("unknown subcommand %q", cmd)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `listo -- the 1 list for all your content.

usage:
  listo serve           run the web server
  listo help            this message
  listo version         display build info

`)
}
