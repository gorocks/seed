package main

import (
	"flag"
	"github.com/Guazi-inc/seed/cmd"
	"github.com/Guazi-inc/seed/cmd/command"
	"github.com/Guazi-inc/seed/utils"
	"log"
	"os"
)

func main() {
	//currentpath, _ := os.Getwd()

	flag.Usage = cmd.Usage
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()

	if len(args) < 1 {
		cmd.Usage()
		os.Exit(2)
		return
	}

	if args[0] == "help" {
		cmd.Help(args[1:])
		return
	}
	for _, c := range commands.AvailableCommands {
		if c.Name() == args[0] && c.Run != nil {
			c.Flag.Usage = func() { c.Usage() }
			if c.CustomFlags {
				args = args[1:]
			} else {
				c.Flag.Parse(args[1:])
				args = c.Flag.Args()
			}

			if c.PreRun != nil {
				c.PreRun(c, args)
			}
			//// Check if current directory is inside the GOPATH,
			//// if so parse the packages inside it.
			//if utils.IsInGOPATH(currentpath) && cmd.IfGenerateDocs(c.Name(), args) {
			//	swaggergen.ParsePackagesFromDir(currentpath)
			//}
			os.Exit(c.Run(c, args))
			return
		}
	}

	utils.PrintErrorAndExit("Unknown subcommand", cmd.ErrorTemplate)
}
