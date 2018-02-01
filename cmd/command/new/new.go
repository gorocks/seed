package new

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	path "path/filepath"
	"regexp"
	"strings"

	"github.com/Guazi-inc/seed/cmd/command"
	"github.com/Guazi-inc/seed/cmd/command/version"
	"github.com/Guazi-inc/seed/logger"
	"github.com/Guazi-inc/seed/logger/color"
	"github.com/Guazi-inc/seed/utils"
)

var CmdNew = &commands.Command{
	UsageLine: "new -name=[appname]",
	Short:     "Creates a Grpc Golang app",
	Long: `
Creates a  application for the given app name in the current directory.
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    CreateApp,
}
var (
	appName   string
	groupName string
	template  string
	tempPath  string
)

func init() {
	fs := flag.NewFlagSet("new", flag.ContinueOnError)
	fs.StringVar(&appName, "n", "", "set a name for application")
	fs.StringVar(&groupName, "g", "finance", "this application belong which group")
	fs.StringVar(&template, "tn", "eipis-apply", "template name,use which template")
	fs.StringVar(&tempPath, "tp", "", "template path")
	CmdNew.Flag = *fs
	commands.AvailableCommands = append(commands.AvailableCommands, CmdNew)
}

func CreateApp(cmd *commands.Command, args []string) int {
	if err := cmd.Flag.Parse(args); err != nil {
		logger.Log.Fatalf("Error while parsing flags: %v", err.Error())
	}

	if len(appName) == 0 {
		logger.Log.Fatal("Argument [appname] is missing")
	}

	if len(tempPath) == 0 {
		logger.Log.Fatal("Argument [template path] is missing")
	}
	currpath, _ := os.Getwd()
	appPath := path.Join(currpath, appName)

	if utils.IsExist(appPath) {
		logger.Log.Errorf(colors.Bold("Application '%s' already exists"), appPath)
		logger.Log.Warn(colors.Bold("Do you want to overwrite it? [Yes|No] "))
		if !utils.AskForConfirmation() {
			os.Exit(2)
		}
	}

	logger.Log.Info("Creating application...")
	return CreateFile(cmd, tempPath, appPath)
}

//创建文件
func CreateFile(cmd *commands.Command, templatePath string, appPath string) int {
	files, _ := ioutil.ReadDir(templatePath)
	isTruePath := false
	for _, fi := range files {
		if fi.IsDir() && fi.Name() == template { //找到当前目录下名字为template的文件夹
			isTruePath = true
			//创建总项目目录
			createAllDir(cmd, templatePath)
			//遍历文件夹建立模板文件
			err := path.Walk(path.Join(templatePath, template), func(tempPath string, info os.FileInfo, err error) error {
				if info == nil {
					return err
				}
				if !info.IsDir() {
					data, err := ioutil.ReadFile(tempPath)
					if err != nil {
						return err
					}
					//解析内容识别path部分，如果有path按照path进行文件层级建立，如果没有按照模板目录层级建立
					reg, err := regexp.Compile(`^#path=.+\n`)
					if err != nil {
						return err
					}
					loc := reg.FindIndex(data)
					if len(loc) > 0 { //存在自定义path
						str := string(data)[loc[0]:loc[1]]
						realPath := strings.Split(str, "=")[1]
						careateFile(cmd, appPath, realPath, string(data)[loc[1]:])
					} else {
						realPath := strings.Split(strings.Split(tempPath, template)[1], ".template")[0]
						careateFile(cmd, appPath, realPath, string(data))
					}
				}
				return nil
			})
			if err != nil {
				logger.Log.Error(err.Error())
				return 1
			}
			break
		}
	}
	if !isTruePath {
		logger.Log.Fatalf("the path not find template ,path:%v", templatePath)
	}
	logger.Log.Success("New application successfully created!")
	return 0
}

func careateFile(cmd *commands.Command, templatePath, realPath string, content string) {
	arr := strings.Split(realPath, "/")
	dir := templatePath
	for _, v := range arr[:len(arr)-1] {
		if v == "" {
			continue
		}
		dir = path.Join(dir, v)
	}
	//创建文件需要目录
	createAllDir(cmd, dir)
	//创建文件
	content = strings.Replace(strings.Replace(content, "{{.Appname}}", appName, -1), "{{.GroupName}}", groupName, -1)
	writeFile(cmd, path.Join(dir, strings.Replace(arr[len(arr)-1], "\n", "", -1)), content)
}
func createAllDir(cmd *commands.Command, filePath string) {
	output := cmd.Out()
	if utils.IsExist(filePath) {
		return
	}
	os.MkdirAll(filePath, 0777)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filePath+string(path.Separator), "\x1b[0m")
}
func writeFile(cmd *commands.Command, filePath string, content string) {
	output := cmd.Out()
	utils.WriteToFile(filePath, content)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", filePath+string(path.Separator), "\x1b[0m")
}
