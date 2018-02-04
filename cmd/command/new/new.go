package new

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	path "path/filepath"
	"strings"

	"github.com/Guazi-inc/seed/cmd/command"
	"github.com/Guazi-inc/seed/cmd/command/version"
	"github.com/Guazi-inc/seed/logger"
	"github.com/Guazi-inc/seed/logger/color"
	"github.com/Guazi-inc/seed/utils"
)

var CmdNew = &commands.Command{
	UsageLine: "new -n=[appname] -tp=[template path]",
	Short:     "Creates a  app for template",
	Long: `
Creates a  application for the given app name and template in the current directory.
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
		logger.Fatalf("Error while parsing flags: %v", err.Error())
	}

	if len(appName) == 0 {
		logger.Fatal("Argument [appname] is missing")
	}

	if len(tempPath) == 0 {
		logger.Fatal("Argument [template path] is missing")
	}
	currpath, _ := os.Getwd()
	appPath := path.Join(currpath, appName)

	if utils.IsExist(appPath) {
		logger.Errorf(colors.Bold("Application '%s' already exists"), appPath)
		logger.Warn(colors.Bold("Do you want to overwrite it? [Yes|No] "))
		if !utils.AskForConfirmation() {
			os.Exit(2)
		}
	}

	logger.Info("Creating application...")
	return CreateFile(tempPath, appPath)
}

//创建文件
func CreateFile(templatePath string, appPath string) int {
	files, _ := ioutil.ReadDir(templatePath)
	isTruePath := false
	for _, fi := range files {
		if fi.IsDir() && fi.Name() == template { //找到当前目录下名字为template的文件夹
			isTruePath = true
			//创建总项目目录
			createAllDir(templatePath)
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
					realPath := strings.Split(strings.Split(tempPath, template)[1], ".tmpl")[0]
					careateFile(appPath, realPath, string(data))
				}
				return nil
			})
			if err != nil {
				logger.Error(err.Error())
				return 1
			}
			break
		}
	}
	if !isTruePath {
		logger.Fatalf("the path not find %s template ,path:%v", template, templatePath)
	}
	logger.Success("New application successfully created!")
	return 0
}

func careateFile(templatePath, realPath string, content string) {
	arr := strings.Split(realPath, "/")
	dir := templatePath
	for _, v := range arr[:len(arr)-1] {
		if v == "" {
			continue
		}
		dir = path.Join(dir, v)
	}
	//创建文件需要目录
	createAllDir(dir)
	//创建文件
	content = strings.Replace(strings.Replace(content, "{{.Appname}}", appName, -1), "{{.GroupName}}", groupName, -1)
	writeFile(path.Join(dir, strings.Replace(arr[len(arr)-1], "\n", "", -1)), content)
}

//create dir from path
func createAllDir(filePath string) {
	if utils.IsExist(filePath) {
		return
	}
	err := os.MkdirAll(filePath, 0755)
	if err != nil {
		logger.Fatalf("fail create dir:%s ,err:%v", filePath, err)
	}
	logger.Success(fmt.Sprintf("create dir:%v", filePath+string(path.Separator)))
}

//create file
func writeFile(filePath string, content string) {
	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		logger.Fatalf("fail create file %v,err:%v", filePath, err)
	}
	_, err = f.WriteString(content)
	if err != nil {
		logger.Fatalf("fail create file  %v,err:%v", filePath, err)
	}
	//go fmt
	utils.FormatSourceCode(filePath)
	logger.Success(fmt.Sprintf("create file:%v", filePath))
}
