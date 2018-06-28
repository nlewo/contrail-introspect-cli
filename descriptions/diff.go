package descriptions

import (
	"fmt"
	"log"
	"os"

	"github.com/nlewo/contrail-introspect-cli/collection"
)
import cli "gopkg.in/urfave/cli.v2"

func RouteDiff() *cli.Command {
	return &cli.Command{
		Name:      "route-diff",
		Usage:     "Show routes diff between controller and vrouter",
		ArgsUsage: "controller-fqdn vrouter-fqdn vrf",
		Action: func(c *cli.Context) error {
			if c.NArg() != 3 {
				fmt.Printf("Wrong argument number\n")
				cli.ShowSubcommandHelp(c)
				os.Exit(1)
			}
			controller := c.Args().Get(0)
			agent := c.Args().Get(1)
			vrf := c.Args().Get(2)

			rdesc := Route()
			crdesc := CtrlRoute()

			col, err := collection.LoadCollection(crdesc, []string{controller, vrf})
			if err != nil {
				log.Fatal(fmt.Sprintf("Route can not be loaded: %s", err))
			}
			elts := col.SearchFuzzy("")
			croutes := []string{}
			for _, e := range elts {
				route, err := e.GetField("prefix")
				if err != nil {
					log.Fatal(err)
				}
				croutes = append(croutes, route)
			}

			col, err = collection.LoadCollection(rdesc, []string{agent, vrf})
			if err != nil {
				log.Fatal(fmt.Sprintf("Route can not be loaded: %s", err))
			}
			elts = col.SearchFuzzy("")
			aroutes := []string{}
			for _, e := range elts {
				route, _ := e.GetField("src_ip")
				if err != nil {
					log.Fatal(err)
				}
				prefix, _ := e.GetField("src_plen")
				if err != nil {
					log.Fatal(err)
				}
				aroutes = append(aroutes, route+"/"+prefix)
			}

			croutes, aroutes = diff(croutes, aroutes)
			output := fmtDiff(croutes, aroutes)
			fmt.Println(output)
			return nil
		},
	}
}

func diff(a, b []string) ([]string, []string) {
	ca := &a
	cb := &b
	if len(a) > len(b) {
		ca = &b
		cb = &a
	}
	ia := 0
	for ia < len(*ca) {
		va := (*ca)[ia]
		ib := 0
		for ib < len(*cb) {
			vb := (*cb)[ib]
			if va == vb {
				*ca = append((*ca)[:ia], (*ca)[ia+1:]...)
				*cb = append((*cb)[:ib], (*cb)[ib+1:]...)
				ia -= 1
				ib -= 1
			}
			ib += 1
		}
		ia += 1
	}
	return a, b
}

func fmtDiff(a, b []string) string {
	res := `diff --route a/controller/routes b/vrouter/routes
--- a/controller/routes
+++ b/vrouter/routes

`
	for _, va := range a {
		res += "-    " + va + "\n"
	}
	if len(a) > 0 {
		res += "\n"
	}
	for _, vb := range b {
		res += "+    " + vb + "\n"
	}
	return res
}
