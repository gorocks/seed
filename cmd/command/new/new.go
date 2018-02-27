package new

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	path "path/filepath"
	"regexp"
	"strings"
	tmp "text/template"

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
	isGip     bool
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

type serviceTemp struct {
	PackageName string
	PackPath    string
	ServiceName string
	Package     string
	Rpc         []*proto.GFunc
	Imports     []string
}

var (
	isOverwriteAll = false
	isSkipAll      = false
)

func init() {
	fs := flag.NewFlagSet("new", flag.ContinueOnError)
	fs.StringVar(&appName, "n", "", "set a name for application")
	fs.StringVar(&groupName, "g", "finance", "this application belong  with which group")
	fs.StringVar(&protoPath, "pt", "", "proto file path")
	fs.StringVar(&style, "s", "grpcweb", "can choose grpcweb,grpcservice,consumer,all")
	fs.StringVar(&template, "tn", "eipis-apply", "template name,use which template")
	fs.StringVar(&tempPath, "tp", "", "template path")
	fs.BoolVar(&isGip, "gip", false, "do gip install -v requirements.txt")
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
			isOverwriteAll = true
			logger.Info("Overwrite all file...")
		}
	}

	logger.Info("Creating application...")
	return CreateFile(tempPath, appPath)
}

const pathSeparator = string(os.PathSeparator)

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// remove additional prefix directory name
		destPath := path.Join(dest, f.Name[strings.Index(f.Name, pathSeparator)+1:])

		if f.FileInfo().IsDir() {
			os.MkdirAll(destPath, os.ModePerm)
		} else {
			// create f's parent directory if not exists
			if err = os.MkdirAll(destPath[:strings.LastIndex(destPath, pathSeparator)], os.ModePerm); err != nil {
				return err
			}

			destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer destFile.Close()

			if _, err = io.Copy(destFile, rc); err != nil {
				return err
			}
		}
	}
	return nil
}

// modified from https://git.io/vA6BQ
func tempFileName() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return path.Join(os.TempDir(), hex.EncodeToString(randBytes))
}

// parseZip will fetch a zip by url, and unzip it to [toPath].
func parseZip(url string, toPath string) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	tmpZip, err := ioutil.TempFile("", "")
	if err != nil {
		return
	}
	// clean up
	defer os.Remove(tmpZip.Name())

	if _, err = tmpZip.Write(content); err != nil {
		return
	}
	if err = tmpZip.Close(); err != nil {
		return
	}

	return unzip(tmpZip.Name(), toPath)
}

func isNetZip(b []byte) bool {
	matched, _ := regexp.Match(`^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)zip$`, b)
	return matched
}

// CreateFile 创建文件
func CreateFile(templatePath string, appPath string) int { // nolint: gocyclo
	if isNetZip([]byte(templatePath)) {
		tf := tempFileName()
		err := parseZip(templatePath, tf)
		if err != nil {
			logger.Fatal(colors.Bold(err.Error()))
		}
		defer os.RemoveAll(tf)
		templatePath = tf
	}
	files, _ := ioutil.ReadDir(templatePath)
	var (
		isTruePath           = false
		isNeedGeneratorProto = false
		isWeb                = false
		isGrpc               = false
		isconsumer           = false
	)
	if len(protoPath) > 0 { //protopath不为空，先生成对应的proto service文件
		isNeedGeneratorProto = true
	}
	switch style {
	case "grpcweb", "web", "gw":
		isWeb = true
	case "grpcservice", "gs":
		isGrpc = true
	case "consumer", "c":
		isconsumer = true
	case "all", "a":
		isWeb = true
		isGrpc = true
	}
	for _, fi := range files {
		if fi.IsDir() && fi.Name() == template { //找到当前目录下名字为template的文件夹
			isTruePath = true
			//创建总项目目录
			createAllDir(appPath)
			if isNeedGeneratorProto { //是否是通过proto生成service
				genService(appPath, protoPath, isWeb, isGrpc)
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
						if isNeedGeneratorProto && v == "service" ||
							!isWeb && v == "grpcweb" ||
							!isGrpc && v == "grpcserver" ||
							!isconsumer && v == "consumer" {
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
	if fileName == "requirements.txt" && isGip {
		doGip(path.Join(fileRPath, fileName))
	}
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
	logger.Success(fmt.Sprintf("Create dir:%v", filePath+string(path.Separator)))
}

//create file
func writeFile(filePath string, content string) {
	if isSkipAll {
		logger.Warnf("Skip %v", filePath)
		return
	}
	if utils.IsExist(filePath) && !isOverwriteAll {
		logger.Errorf(colors.Bold("file '%s' already exists"), filePath)
		logger.Warn(colors.Bold("Do you want to overwrite it , skip it , skip all or overwrite all,yes is just overwrite current file? [yes|overwrite|skip|skip all|overwrite all] "))
		switch utils.AskForConfirmation() {
		case "skip":
			logger.Infof("Skip %v this file", filePath)
			return
		case "yes":
			logger.Infof("Overwrite %v this file", filePath)
		case "skip all":
			isSkipAll = true
			logger.Infof("skip all begin this file %v", filePath)
			return
		case "overWrite all":
			isOverwriteAll = true
			logger.Infof("overwrite all begin current file :%v", filePath)
		}
	}

	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		logger.Fatalf("Fail create file %v,err:%v", filePath, err)
	}
	_, err = f.WriteString(content)
	if err != nil {
		logger.Fatalf("Fail create file  %v,err:%v", filePath, err)
	}
	//判断文件后缀 进行格式化 todo 多种格式化
	if strings.HasSuffix(filePath, ".go") {
		//go fmt
		utils.FormatSourceCode(filePath)
	}
	logger.Success(fmt.Sprintf("Create file:%v", filePath))
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

type service struct {
	AppGoPath      string `json:"app_go_path"`
	ProtoPath      []string
	RegisterServer map[string][]string
	IsWed          bool
	IsGrpc         bool
}

var serviceTmp = `package service

import (
	{{range $k,$v:=.ProtoPath}}
"{{$v}}"
	{{end}}
	"google.golang.org/grpc"
	{{range $k,$v:=.RegisterServer}}
{{$k}}2 "{{$.AppGoPath}}/service/{{$k}}"
	{{end}}
)

{{if .IsWed}}

// RegisterGRPCWebServices RegisterAll grpc web services
func RegisterGRPCWebServices(grpcServer *grpc.Server) {
	{{range $k,$v:=.RegisterServer}}
{{range $i,$j:=$v}}
{{$k}}.Register{{$j}}Server(grpcServer,&{{$k}}2.{{$j}}{})
{{end}}
{{end}}
}
{{end}}

{{if .IsGrpc}}
// RegisterGRPCServices RegisterAll grpc services
func RegisterGRPCServices(grpcServer *grpc.Server) {
{{range $k,$v:=.RegisterServer}}
{{range $i,$j:=$v}}
{{$k}}.Register{{$j}}Server(grpcServer,&{{$k}}2.{{$j}}{})
{{end}}
{{end}}
}
{{end}}

`

func genService(appPath string, protoPaths string, isWed, isGrpc bool) {
	t := tmp.New("Service") //创建一个模板
	t.Funcs(tmp.FuncMap{
		"tmp": ServiceTemplPath,
	})
	p, err := t.Parse(serviceTmpl)
	if err != nil {
		panic(err)
	}
	ser := service{
		AppGoPath: path.Join(utils.GetUsefulPath(appPath, "src", false), appName),
		IsWed:     isWed,
		IsGrpc:    isGrpc,
	}
	registerMap := make(map[string][]string)
	pPath := make([]string, 0)
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

			servicePackPath := utils.GetUsefulPath(paths, "proto", true)
			pPath = append(pPath, servicePackPath)
			//建某一个proto的文件夹
			createAllDir(servicePath)

			sNames := make([]string, 0)
			for _, v := range g.Service {
				stemp := serviceTemp{
					ServiceName: v.ServiceName,
					Imports:     v.Imports,
					PackPath:    servicePackPath,
					Package:     g.Package,
					PackageName: strings.ToLower(v.ServiceName),
					Rpc:         v.Rpc,
				}
				if v.ServiceName != "" {
					sNames = append(sNames, v.ServiceName)
				}
				var content bytes.Buffer
				err = p.Execute(&content, stemp)
				if err != nil {
					panic(err)
				}
				//建立某一个service的文件
				writeFile(path.Join(servicePath, strings.ToLower(v.ServiceName))+".go", content.String())
			}
			if len(sNames) > 0 && len(g.Service) > 0 {
				registerMap[g.Package] = sNames
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	ser.RegisterServer = registerMap
	ser.ProtoPath = utils.RmDuplicate(pPath)
	st := tmp.New("Service-go") //创建一个模板
	sp, err := st.Parse(serviceTmp)
	if err != nil {
		panic(err)
	}
	var content bytes.Buffer
	err = sp.Execute(&content, &ser)
	if err != nil {
		panic(err)
	}
	writeFile(path.Join(appPath, "service", "service")+".go", content.String())

}

func doGip(reqPath string) {
	if utils.CheckGip() {
		utils.DoGipInstall(reqPath)
	}
}

func ServiceTemplPath(str string, str2 string) string {
	if strings.Contains(str, ".") {
		return str
	}
	return str2 + "." + str
}
