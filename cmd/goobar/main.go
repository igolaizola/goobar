package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"

	"github.com/igolaizola/goobar"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

// Build flags
var Version = ""
var Commit = ""
var Date = ""

func main() {
	// Create signal based context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Launch command
	cmd := newCommand()
	if err := cmd.ParseAndRun(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func newCommand() *ffcli.Command {
	fs := flag.NewFlagSet("goobar", flag.ExitOnError)

	return &ffcli.Command{
		ShortUsage: "goobar [flags] <subcommand>",
		FlagSet:    fs,
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
		Subcommands: []*ffcli.Command{
			newVersionCommand(),
			newServeCommand(),
			newRunCommand(),
		},
	}
}

func newVersionCommand() *ffcli.Command {
	return &ffcli.Command{
		Name:       "version",
		ShortUsage: "goobar version",
		ShortHelp:  "print version",
		Exec: func(ctx context.Context, args []string) error {
			v := Version
			if v == "" {
				if buildInfo, ok := debug.ReadBuildInfo(); ok {
					v = buildInfo.Main.Version
				}
			}
			if v == "" {
				v = "dev"
			}
			versionFields := []string{v}
			if Commit != "" {
				versionFields = append(versionFields, Commit)
			}
			if Date != "" {
				versionFields = append(versionFields, Date)
			}
			fmt.Println(strings.Join(versionFields, " "))
			return nil
		},
	}
}

func newServeCommand() *ffcli.Command {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	_ = fs.String("config", "", "config file (optional)")

	port := fs.Int("port", 0, "port number")

	return &ffcli.Command{
		Name:       "serve",
		ShortUsage: "goobar serve [flags] <key> <value data...>",
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ff.PlainParser),
			ff.WithEnvVarPrefix("GOOBAR"),
		},
		ShortHelp: "run goobar server",
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			if *port == 0 {
				return errors.New("missing port")
			}
			return goobar.Serve(ctx, *port)
		},
	}
}

func newRunCommand() *ffcli.Command {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	_ = fs.String("config", "", "config file (optional)")

	return &ffcli.Command{
		Name:       "run",
		ShortUsage: "goobar serve [flags] <key> <value data...>",
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ff.PlainParser),
			ff.WithEnvVarPrefix("GOOBAR"),
		},
		ShortHelp: "run goobar action",
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			return goobar.Run(ctx)
		},
	}
}
