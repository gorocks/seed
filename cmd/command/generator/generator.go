package generator

import (
	"os"
	"path/filepath"
	"strings"

	"flag"

	"github.com/Guazi-inc/seed/cmd/command"
	"github.com/Guazi-inc/seed/cmd/command/version"
	"github.com/Guazi-inc/seed/logger"
	"github.com/Guazi-inc/seed/utils"
)

var CmdGen = &commands.Command{
	UsageLine: "gen [Command]",
	Short:     "seed generator proto avro and db model",
	Long: `
can do generator groto avro db-model to go file.
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    Gen,
}

var (
	fPath   string
	outPath string
)

func init() {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.StringVar(&fPath, "p", "", "proto file path ")
	fs.StringVar(&outPath, "o", "", "proto out file path ")
	CmdGen.Flag = *fs
	commands.AvailableCommands = append(commands.AvailableCommands, CmdGen)
}

func Gen(cmd *commands.Command, args []string) int {
	if err := cmd.Flag.Parse(args[1:]); err != nil {
		logger.Fatalf("Error while parsing flags: %v", err.Error())
	}
	logger.Info(fPath)
	if len(fPath) == 0 {
		logger.Fatalf("Argument [p] is missing")
	}
	mcmd := args[0]
	switch mcmd {
	case "proto":
		logger.Info("generator grpc rpc code from .proto")
		genProto()
	default:
		logger.Fatal("Command is missing")
	}
	return 0
}

func genProto() {
	if !utils.CheckProtoc() {
		logger.Fatal("protoc is missing,please install protoc")
	}
	gps := utils.GetGOPATHs()
	if len(gps) == 0 {
		logger.Fatalf("GOPATH environment variable is not set or empty")
	}

	gopath := gps[0]
	logger.Infof("GOPATH: %s", gopath)
	if outPath == "" {
		outPath = gopath + "/src"
	}
	logger.Infof("out_put_file_path:%s", outPath)
	filepath.Walk(fPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".proto") {
			logger.Warnf("%s not proto file ", path)
			return nil
		}
		arr := strings.Split(path, "/")
		iPath := ""
		for k, v := range arr {
			if v == "proto" {
				iPath = strings.Join(arr[:k], "/")
				break
			}
		}
		utils.ProtocGenGo(path, outPath, iPath)
		return nil
	})
}
