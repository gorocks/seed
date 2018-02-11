package validate

import (
	"flag"

	"github.com/Guazi-inc/seed/cmd/command"
	"github.com/Guazi-inc/seed/cmd/command/version"
	"github.com/Guazi-inc/seed/utils"
)

var CmdValidate = &commands.Command{
	UsageLine: "validate",
	Short:     "do code validate use gometalinter",
	Long: `
do code validate use gometalinter.
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    doValidate,
}

var (
	jsonPath  string
	isInstall bool
)

func init() {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.StringVar(&jsonPath, "path", "", "gometalinter.json path")
	fs.BoolVar(&isInstall, "i", false, "install and update gometalinter when gometalinter not exists")
	CmdValidate.Flag = *fs
	commands.AvailableCommands = append(commands.AvailableCommands, CmdValidate)
}

func doValidate(cmd *commands.Command, args []string) int {
	if isInstall {
		utils.InstallAndUpdateGometalinter()
	}
	if len(jsonPath) == 0 {
		utils.DoGometalinterCI()
	} else {
		utils.DoGometalinterFromJson(jsonPath)
	}
	return 0
}
