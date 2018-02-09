package new

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	path "path/filepath"
	"strings"
	tmp "text/template"

	"bytes"

	"github.com/Guazi-inc/seed/cmd/command"
	"github.com/Guazi-inc/seed/cmd/command/generator/proto"
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
	protoPath string
	style     string
	template  string
	tempPath  string
)

var serviceTmpl = `package {{.Package}}

import (
	{{range $k,$v:=.Imports}}
	"{{$v}}"
	{{end}}
	"{{.PackPath}}"
	"golang.org/x/net/context"
)

type {{.ServiceName}} struct{}

{{range .Rpc}}

func (s *{{$.ServiceName}}) {{.FunName}}(ctx context.Context, in *{{ tmp .Request $.Package}}) (*{{ tmp .Response $.Package }}, error) {
	return &{{ tmp .Response $.Package }}{}, nil
}


{{end}}

`
var service = `package service

import (
	"google.golang.org/grpc"
)

// RegisterGRPCWebServices RegisterAll grpc web services
func RegisterGRPCWebServices(grpcServer *grpc.Server) {
	
}

// RegisterGRPCServices RegisterAll grpc services
func RegisterGRPCServices(grpcServer *grpc.Server) {
}

`

type serviceTemp struct {
	PackageName string
	PackPath    string
	ServiceName string
	Package     string
	Rpc         []*proto.GFunc
	Imports     []string
}

var isOroverwriteAll = false

func init() {
	fs := flag.NewFlagSet("new", flag.ContinueOnError)
	fs.StringVar(&appName, "n", "", "set a name for application")
	fs.StringVar(&groupName, "g", "finance", "this application belong which group")
	fs.StringVar(&protoPath, "pt", "", "proto path")
	fs.StringVar(&style, "s", "grpcweb", "can choose grpcweb,grpcservice,consumer")
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
		logger.Warn(colors.Bold("Do you want to overwrite all ? [Yes|No] "))
		str := utils.AskForConfirmation()
		if str == "yes" || str == "all" {
			isOroverwriteAll = true
		}
	}

	logger.Info("Creating application...")
	return CreateFile(tempPath, appPath)
}

//创建文件
func CreateFile(templatePath string, appPath string) int {

	files, _ := ioutil.ReadDir(templatePath)
	isTruePath := false
	isNeedGeneratorProto := false
	if len(protoPath) > 0 { //protopath不为空，先生成对应的proto service文件
		isNeedGeneratorProto = true
	}
	for _, fi := range files {
		if fi.IsDir() && fi.Name() == template { //找到当前目录下名字为template的文件夹
			isTruePath = true
			//创建总项目目录
			createAllDir(appPath)
			if isNeedGeneratorProto { //是否是通过proto生成service
				genService(appPath, protoPath)
			}
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
					arr := strings.Split(tempPath, template)
					if len(arr) < 1 {
						logger.Fatalf("the path not find %s template ,path:%v", template, templatePath)
					}
					at := strings.Split(arr[1], "/")
					fileDirPath := appPath
					rfileName := ""
					for k, v := range at {
						//处理path，
						if isNeedGeneratorProto && v == "service" {
							return nil
						}
						if k == (len(at) - 1) {
							v = strings.TrimSuffix(v, ".tmpl")
							rfileName = v
							continue
						}
						fileDirPath = path.Join(fileDirPath, strings.Replace(v, "/n", "", -1))
					}
					careateFile(fileDirPath, rfileName, string(data))
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

func careateFile(fileRPath, fileName string, content string) {
	//创建文件需要目录
	createAllDir(fileRPath)
	//创建文件
	content = strings.Replace(strings.Replace(content, "{{.Appname}}", appName, -1), "{{.GroupName}}", groupName, -1)
	writeFile(path.Join(fileRPath, fileName), content)
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
	if utils.IsExist(filePath) && !isOroverwriteAll {
		logger.Errorf(colors.Bold("Application '%s' already exists"), filePath)
		logger.Warn(colors.Bold("Do you want to overwrite it ? [Yes|No|skip|all] "))
		switch utils.AskForConfirmation() {
		case "no", "skip":
			return
		case "yes":
		case "all":
			isOroverwriteAll = true
		}
	}
	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		logger.Fatalf("fail create file %v,err:%v", filePath, err)
	}
	_, err = f.WriteString(content)
	if err != nil {
		logger.Fatalf("fail create file  %v,err:%v", filePath, err)
	}
	//判断文件后缀 进行格式化 todo 多种格式化
	if strings.HasSuffix(filePath, ".go") {
		//go fmt
		utils.FormatSourceCode(filePath)
	}
	logger.Success(fmt.Sprintf("create file:%v", filePath))
}

//创建service
func generatorServiceFromProto(fileName string) (*proto.GeneratorProto, error) {
	var g proto.GeneratorProto
	err := g.Generator(fileName)
	if err != nil {
		return nil, err
	}
	return &g, err
}

func genService(appPath string, protoPaths string) {
	t := tmp.New("Service") //创建一个模板
	t.Funcs(tmp.FuncMap{
		"tmp": ServiceTemplPath,
	})
	p, err := t.Parse(serviceTmpl)
	if err != nil {
		panic(err)
	}
	err = path.Walk(protoPaths, func(paths string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if !info.IsDir() {
			g, err := generatorServiceFromProto(paths)
			if err != nil {
				return err
			}
			servicePath := path.Join(appPath, "service", g.Package)
			arr := strings.Split(paths, "/")
			servicePackPath := ""
			for k, v := range arr {
				if v == "proto" {
					servicePackPath = strings.Join(arr[k:len(arr)-1], "/")
				}
			}
			//建某一个proto的文件夹
			createAllDir(servicePath)
			for _, v := range g.Service {
				stemp := serviceTemp{
					ServiceName: v.ServiceName,
					Imports:     v.Imports,
					PackPath:    servicePackPath,
					Package:     g.Package,
					PackageName: strings.ToLower(v.ServiceName),
					Rpc:         v.Rpc,
				}
				var content bytes.Buffer
				err = p.Execute(&content, stemp)
				if err != nil {
					panic(err)
				}
				//建立某一个service的文件
				writeFile(path.Join(servicePath, strings.ToLower(v.ServiceName))+".go", content.String())
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func ServiceTemplPath(str string, str2 string) string {
	if strings.Contains(str, ".") {
		return str
	}
	return str2 + "." + str
}
