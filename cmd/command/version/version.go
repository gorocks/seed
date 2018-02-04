package version

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Guazi-inc/seed/cmd/command"
	"github.com/Guazi-inc/seed/logger"
	"github.com/Guazi-inc/seed/logger/color"

	"gopkg.in/yaml.v2"
)

const verboseVersionBanner string = `%s%s seed v{{ .SeedVersion }}%s
%s%s
├── SeedVersion : {{ .SeedVersion }}
├── GoVersion : {{ .GoVersion }}
├── GOOS      : {{ .GOOS }}
├── GOARCH    : {{ .GOARCH }}
├── NumCPU    : {{ .NumCPU }}
├── GOPATH    : {{ .GOPATH }}
├── GOROOT    : {{ .GOROOT }}
├── Compiler  : {{ .Compiler }}
└── Date      : {{ Now "Monday, 2 Jan 2006" }}%s
`

const shortVersionBanner = `seedVersion:{{ .SeedVersion }}
`

var CmdVersion = &commands.Command{
	UsageLine: "version",
	Short:     "Prints the current Seed version",
	Long:      `Prints the current Seed and Go version alongside the platform information.`,
	Run:       versionCmd,
}
var outputFormat string

const version = "v0.0.1"

func init() {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	fs.StringVar(&outputFormat, "o", "", "Set the output format. Either json or yaml.")
	CmdVersion.Flag = *fs
	commands.AvailableCommands = append(commands.AvailableCommands, CmdVersion)
}

func versionCmd(cmd *commands.Command, args []string) int {

	cmd.Flag.Parse(args)
	stdout := cmd.Out()

	if outputFormat != "" {
		runtimeInfo := RuntimeInfo{
			GetGoVersion(),
			runtime.GOOS,
			runtime.GOARCH,
			runtime.NumCPU(),
			os.Getenv("GOPATH"),
			runtime.GOROOT(),
			runtime.Compiler,
			version,
		}
		switch outputFormat {
		case "json":
			{
				b, err := json.MarshalIndent(runtimeInfo, "", "    ")
				if err != nil {
					logger.Error(err.Error())
				}
				fmt.Println(string(b))
				return 0
			}
		case "yaml":
			{
				b, err := yaml.Marshal(&runtimeInfo)
				if err != nil {
					logger.Error(err.Error())
				}
				fmt.Println(string(b))
				return 0
			}
		}
	}

	coloredBanner := fmt.Sprintf(verboseVersionBanner, "\x1b[35m", "\x1b[1m",
		"\x1b[0m", "\x1b[32m", "\x1b[1m", "\x1b[0m")
	InitBanner(stdout, bytes.NewBufferString(coloredBanner))
	return 0
}

// ShowShortVersionBanner prints the short version banner.
func ShowShortVersionBanner() {
	output := colors.NewColorWriter(os.Stdout)
	InitBanner(output, bytes.NewBufferString(colors.MagentaBold(shortVersionBanner)))
}

func GetGoVersion() string {
	var (
		cmdOut []byte
		err    error
	)
	if cmdOut, err = exec.Command("go", "version").Output(); err != nil {
		logger.Fatalf("There was an error running 'go version' command: %s", err)
	}
	return strings.Split(string(cmdOut), " ")[2]
}
