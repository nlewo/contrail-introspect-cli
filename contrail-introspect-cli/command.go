package main

import "fmt"
import "os"
import "strings"
import "log"

import cli "gopkg.in/urfave/cli.v2"

import "github.com/nlewo/contrail-introspect-cli/collection"
import "github.com/nlewo/contrail-introspect-cli/utils"

func GenCommand(descCol collection.DescCollection, name string, usage string) *cli.Command {
	return &cli.Command{
		Name:      name,
		Usage:     usage,
		ArgsUsage: fmt.Sprintf("%s\n", strings.Join(descCol.PageArgs, " ")),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "long", Aliases: []string{"l"},
				Usage: "Long format",
			},
			&cli.BoolFlag{
				Name: "xml", Aliases: []string{"x"},
				Usage: "XML output format",
			},
			&cli.BoolFlag{
				Name:  "from-file",
				Usage: "Load file instead of URL (for debugging)",
			},
			&cli.BoolFlag{
				Name: "url", Aliases: []string{"u"},
				Usage: "Just show the used URL",
			},
			&cli.StringFlag{
				Name: "search", Aliases: []string{"s"},
				Usage: fmt.Sprintf("Fuzzy search by %s", descCol.PrimaryField),
				Value: "",
			},
			&cli.StringFlag{
				Name: "strict-search", Aliases: []string{"S"},
				Usage: fmt.Sprintf("Strict search by %s", descCol.PrimaryField),
				Value: "",
			},
		},
		Action: func(c *cli.Context) error {
			var page collection.Sourcer
			if c.IsSet("from-file") {
				page = collection.File{Path: c.Args().Get(0)}
			} else {
				if c.NArg() < len(descCol.PageArgs) {
					cli.ShowSubcommandHelp(c)
					os.Exit(1)
				}
				page = descCol.PageBuilder(c.Args().Slice())
			}
			col, e := page.Load(descCol)
			if e != nil {
				log.Fatal(e)
			}

			if c.IsSet("url") {
				fmt.Println(col.Url)
				return nil
			}

			var list collection.Shower

			if c.String("s") != "" {
				list = col.SearchFuzzy(c.String("s"))
			} else if c.String("S") != "" {
				list = col.SearchStrict(c.String("S"))
			} else {
				list = col
			}

			if c.IsSet("xml") {
				list.Xml()
				return nil
			}
			if c.IsSet("long") {
				list.Long()
				return nil
			}
			list.Short()

			return nil
		},
		ShellComplete: func(c *cli.Context) {
			// We only complete the first argument
			if c.NArg() == 0 {
				for _, fqdn := range utils.HostMap {
					fmt.Println(fqdn)
				}
			}
		},
	}
}
