package migrate

import (
	"flag"

	"github.com/Guazi-inc/seed/cmd/command"
	"github.com/Guazi-inc/seed/cmd/command/version"
	"github.com/Guazi-inc/seed/logger"
)

var CmdMigrate = &commands.Command{
	UsageLine: "migrate [Command]",
	Short:     "Runs database migrations",
	Long: `The command 'migrate' allows you to run database migrations to keep it up-to-date.

  ▶ {{"To run all the migrations:"|bold}}

    $ bee migrate [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"]
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    RunMigration,
}
var (
	mDriver string
	mConn   string
)

func init() {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.StringVar(&mDriver, "driver", "", "Database driver. Either mysql, postgres or sqlite.")
	fs.StringVar(&mConn, "conn", "", "Connection string used by the driver to connect to a database instance.")
	CmdMigrate.Flag = *fs
	commands.AvailableCommands = append(commands.AvailableCommands, CmdMigrate)
}

func RunMigration(cmd *commands.Command, args []string) int {
	if err := cmd.Flag.Parse(args); err != nil {
		logger.Fatalf("Error while parsing flags: %v", err.Error())
	}
	return 0
}
