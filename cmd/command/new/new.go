package new

import (
	"fmt"
	"os"
	path "path/filepath"
	"strings"

	"github/Guazi-inc/seed/logger/color"
	"github/Guazi-inc/seed/utils"
	"github/Guazi-inc/seed/logger"
	"github/Guazi-inc/seed/cmd/command"
	"github/Guazi-inc/seed/cmd/command/version"
)

var CmdNew = &commands.Command{
	UsageLine: "new [appname]",
	Short:     "Creates a Grpc Golang app",
	Long: `
Creates a Golang application for the given app name in the current directory.

  The command 'new' creates a folder named [appname] and generates the following structure:

            ├── {{"cmd"|foldername}}
            │     └── {{"consumer"|foldername}}
            │           └── main.go
            │     └── {{"grpcserver"|foldername}}
            │           └── main.go
            ├── {{"databases"|foldername}}
            │     └── init-tables.py
            │     └── init.sql
            ├── {{"fixtures"|foldername}}
            │     └── {{"apply"|foldername}}
            │           └── user.yml
            ├── {{"med"|foldername}}
            │     └──  med.yml
            │     └── vars.yml
            ├── {{"model"|foldername}}
            │     └── {{"user"|foldername}}
            │           └── user.go
            │           └── user_test.go
            ├── {{"service"|foldername}}
            │     └── {{"preaudit"|foldername}}
            │           └── service.go
            ├── .gitignore
            ├── .gitlab-ci.yml
            ├── README.md
            ├── gip.yml
            ├── gometalinter.json
            ├── requirements.txt
            └── validate.sh

`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    CreateApp,
}

var cmdConsumerMain = `package main

import (
	"avro/finance"
	"fmt"

	"golang.guazi-corp.com/finance/go-common/etcd"
	"golang.guazi-corp.com/znkf/guazi-avro"
)

func main() {

	etcd.EtcdAddr = "etcdv3.guazi-cloud.com:80"
	etcd.Init("finance")

	err := gzavro.ConsumeAvroMessage(&finance.FactDayholeRepay{}, "dayhole_test", true, func(decodedRecord interface{}) error {
		fmt.Println("receive data:")
		_, ok := decodedRecord.(*finance.FactDayholeRepay)
		if !ok {
			fmt.Println("record assert error")
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

}
`

var cmdGrpcserverMain = `package main

import (
	"golang.guazi-corp.com/finance/{{.Appname}}/service/preaudit"
	"golang.guazi-corp.com/finance/go-common/config"
	"golang.guazi-corp.com/znkf/common/server"
	"google.golang.org/grpc"
)

func main() {
	var middlewares []grpc.UnaryServerInterceptor
	server.StartGRPCServerWithCustom(":5000", "finance", middlewares, func(mainConfig *config.MainConfig, server *grpc.Server) {
		//Register GRPC Servers

		preaudit.Register(server)
	})
}
`
var initTables = `"""
init-tables
"""
`

var initSql = `"""
init-sql
"""`
var fixturesApplyUser = `- id: 1
  user_name: "张三"
  id_card_encrypt: "id"
  phone_encrypt: "phone"
  id_checked: false
  phone_checked: false
`

var medMed = `# repo 模块，配置在registry里面的repo名字和在k8s里面的namespace名
repo:
  name: {{.Appname}}           # repo name, 一般一个git一个name，不可和其他组的name重复
  project: finance         # project name，一般一个组一个name，或者一个大组一个，各组不可重复
  namespace: default    # namespace，一般一个组一个，各组不可重复，以后会按namespace赋予不同的权限

# prepare 模块, 用于拉依赖
prepare:
- name: prepare
  version: v1.01                                            # version，每次添加新的依赖之后需要重新prepare时修改
  image: znkf/common-go-1.10:v1.6       # build 依赖的基础镜像，由medusa团队提供和维护
  workdir: /go/src/golang.guazi-corp.com/finance/{{.Appname}}     # 工作目录，代码放置地方，在$GOPATH下，按照自己的git路径放置
  copy:
  - requirements.txt /go/src/golang.guazi-corp.com/finance/{{.Appname}}/
  run:                                                      # 拉依赖包
  - unlink /etc/localtime && ln -s /usr/share/zoneinfo/Etc/GMT-8 /etc/localtime
  - gip install -v requirements.txt


# build 模块，用于编译二进制文件
build:
- name: build                                               # build 镜像名，用于release copy编译好的二进
  base: prepare
  workdir: /go/src/golang.guazi-corp.com/finance/{{.Appname}}
  ignore:                                                   # 不用copy的文件, 需要copy的文件越多，build越慢
  - vendor/*
  - tmp/*
  copy_from:
  - library/avro-schema.build.staging-release:latest /avro /go/src/avro
  copy:
  - . /go/src/golang.guazi-corp.com/finance/{{.Appname}}/
  run:
  - git clean -df
  - cd /tools; git pull;cd -
  - /tools/med/prepare.sh
  - /tools/med/build-image.sh grpcserver
  - /tools/med/build-image.sh consumer

- name: release
  image: library/ubuntu:14.04.4                              # 运行环境基础镜像，由medusa团队提供和维护
  copy_from:                                                # copy build好的二进制文件和相应的配置文件
  # 结构为：build 模块镜像名   build 镜像位置   运行环境位置
  - build /grpcserver /med/grpcserver
  - build /consumer /med/consumer
  run:
  - TZ='Asia/Shanghai'; export TZ;

test:
  - name: validate
    base: build
    env:
      PROJECT: "{{.Appname}}"
      GROUP: "finance"
    command: "gometalinter ./... --config=gometalinter.json"
  - name: test
    base: build
    env:
      PROJECT: "{{.Appname}}"
      GROUP: "finance"
      ETCD_ADDR: "etcd2v3.guazi-cloud.com:80"
    command: "/tools/med/test-coverage.sh"

# deploy 模块，用于配置部署信息
deploy:
- name: grpc
  base: release
  command: /med/grpcserver -listen :80 -etcd_addr {{ etcd_adrr }}
  replicas: 1
  labels:
    app: grpcserver
  domains: {{.Appname}}
  rules:
    - port: 80
      node_port: 32000
      name: grpc
`

var medVars = `dev:
  etcd_adrr : "10.16.11.144:2479,10.16.11.145:2479,10.16.11.143:2479"

online:
  etcd_adrr : "10.16.11.144:2579,10.16.11.145:2579,10.16.11.143:2579"
`
var modelUser = `
// Code generated by model_gen
package user

`
var modelUserTest = `
package user
//test for moed
`

var serverServer = `
package preaudit

import (
	"proto/finance/service/apply"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type preAudit struct{}

func (p *preAudit) Get(ctx context.Context, in *apply.GetRequest) (*apply.GetResponse, error) {
	return &apply.GetResponse{Passed: false}, nil
}

// Register grpc services
func Register(server *grpc.Server) {
	apply.RegisterPreauditServer(server, &preAudit{})
}

`
var gitignore = `.dockerignore
.idea/
.med/
`
var gitlabCi = `
stages:
- prepare
- build
- validate
- test
- push

prepare:
  stage: prepare
  script: med prepare -n prepare

build:
  stage: build
  script:
  - med build -n build
  - med build -n release

validate:
  stage: validate
  script: med test -n validate

test:
  stage: test
  script: med test -n test

push:
   stage: push
   script: med push -n grpc
   artifacts:
     paths:
     - .med/deploy_grpc.yml
`
var readMe = `
read me
`
var gometalinter= `
{
  "Cyclo": 15,
  "Enable": [
    "deadcode",
    "errcheck",
    "gas",
    "goconst",
    "gocyclo",
    "golint",
    "gotype",
    "ineffassign",
    "interfacer",
    "megacheck",
    "structcheck",
    "unconvert",
    "varcheck",
    "vet",
    "vetshadow",
    "gofmt",
    "goimports",
    "unparam",
    "misspell"
  ],
  "Deadline": "120s",
  "Concurrency": 4
}
`

var gip = `
imports:
  - package: golang.guazi-corp.com/finance/go-common
    version: master
    repo: git+ssh://git@git.guazi-corp.com/finance/go-common
    global: true
  - package: golang.guazi-corp.com/finance/go-rule-engine
    version: master
    repo: git+ssh://git@git.guazi-corp.com/finance/go-rule-engine
    global: true
  - package: golang.guazi-corp.com/znkf/common
    version: master
    repo: git+ssh://git@git.guazi-corp.com/znkf/common
    global: true
  - package: golang.guazi-corp.com/finance/data-soup
    version: master
    repo: git+ssh://git@git.guazi-corp.com/finance/data-soup
    global: true
  - package: golang.guazi-corp.com/znkf/process-soup
    version: master
    repo: git+ssh://git@git.guazi-corp.com/znkf/process-soup
    global: true
  - package: golang.org/x/net
    version: d1e1b351919c6738fdeb9893d5c998b161464f0c
    repo: https://github.com/golang/net
  - package: gopkg.in/redis.v4
    version: 4b0862b5fd0a5ae4e63c76476a64655752d6031b
    repo: https://github.com/go-redis/redis
  - package: github.com/caojia/go-orm
    repo: https://github.com/caojia/go-orm
`
var requirement= `
https://github.com/Guazi-inc/go-avro#f8eb3232ed9f7385fb5d91e3c6a6006df016767c,github.com/Guazi-inc/go-avro
`

var validate = `
#!/usr/bin/env bash
gometalinter ./... --config=gometalinter.json
`

func init() {
	commands.AvailableCommands = append(commands.AvailableCommands, CmdNew)
}

func CreateApp(cmd *commands.Command, args []string) int {

	if len(args) != 1 {
		logger.Log.Fatal("Argument [appname] is missing")
	}

	apppath, packpath, err := utils.CheckEnv(args[0])
	if err != nil {
		logger.Log.Fatalf("%s", err)
	}

	if utils.IsExist(apppath) {
		logger.Log.Errorf(colors.Bold("Application '%s' already exists"), apppath)
		logger.Log.Warn(colors.Bold("Do you want to overwrite it? [Yes|No] "))
		if !utils.AskForConfirmation() {
			os.Exit(2)
		}
	}

	logger.Log.Info("Creating application...")
	return CreateGolangApp(cmd,apppath,packpath)
}

func CreateGolangApp(cmd *commands.Command,apppath, packpath string)int{
	output := cmd.Out()
	//创建工程总文件夹
	os.MkdirAll(apppath, 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")

	//创建cmd的目录及目录文件
	os.Mkdir(path.Join(apppath, "cmd"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "cmd")+string(path.Separator), "\x1b[0m")

	os.Mkdir(path.Join(apppath, "cmd", "consumer"), 0755)
	utils.WriteToFile(path.Join(apppath, "cmd", "consumer", "main.go"), cmdConsumerMain)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "cmd", "consumer","main.go")+string(path.Separator), "\x1b[0m")

	os.Mkdir(path.Join(apppath, "cmd", "grpcserver"), 0755)
	utils.WriteToFile(path.Join(apppath, "cmd", "grpcserver", "main.go"), strings.Replace(cmdGrpcserverMain, "{{.Appname}}", packpath, -1))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "cmd", "grpcserver","main.go")+string(path.Separator), "\x1b[0m")

	//创建databases目录及目录文件
	os.Mkdir(path.Join(apppath, "databases"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "databases")+string(path.Separator), "\x1b[0m")

	utils.WriteToFile(path.Join(apppath, "databases", "init-tables.py"), initTables)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "databases","init-tables.py")+string(path.Separator), "\x1b[0m")

	utils.WriteToFile(path.Join(apppath, "databases", "init.sql"), initSql)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "databases","init.sql")+string(path.Separator), "\x1b[0m")

	//创建fixtures目录及目录文件
	os.Mkdir(path.Join(apppath, "fixtures"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "fixtures")+string(path.Separator), "\x1b[0m")

	os.Mkdir(path.Join(apppath, "fixtures","apply"), 0755)
	utils.WriteToFile(path.Join(apppath, "fixtures", "apply", "user.yml"), fixturesApplyUser)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "fixtures","apply","user.yml")+string(path.Separator), "\x1b[0m")

	//创建med目录及目录文件
	os.Mkdir(path.Join(apppath, "med"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "med")+string(path.Separator), "\x1b[0m")

	utils.WriteToFile(path.Join(apppath, "med", "med.yml"), strings.Replace(medMed, "{{.Appname}}", packpath, -1))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "med","med.yml")+string(path.Separator), "\x1b[0m")

	utils.WriteToFile(path.Join(apppath, "med", "vars.yml"), medVars)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "med","vars.yml")+string(path.Separator), "\x1b[0m")

	//创建model目录及目录文件
	os.Mkdir(path.Join(apppath, "model"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "model")+string(path.Separator), "\x1b[0m")

	os.Mkdir(path.Join(apppath, "model","user"), 0755)
	utils.WriteToFile(path.Join(apppath, "model", "user", "user.go"), modelUser)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "model","user","user.go")+string(path.Separator), "\x1b[0m")

	utils.WriteToFile(path.Join(apppath, "model", "user", "user_test.go"),modelUserTest)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "model","user","user_test.go")+string(path.Separator), "\x1b[0m")

	//创建service目录及目录文件
	os.Mkdir(path.Join(apppath, "service"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "service")+string(path.Separator), "\x1b[0m")

	os.Mkdir(path.Join(apppath, "service","preaudit"), 0755)
	utils.WriteToFile(path.Join(apppath, "service", "preaudit", "service.go"), serverServer)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "service","preaudit","service.go")+string(path.Separator), "\x1b[0m")

	//创建.gitignore
	utils.WriteToFile(path.Join(apppath, ".gitignore"), gitignore)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, ".gitignore")+string(path.Separator), "\x1b[0m")

	//创建.gitlab-ci.yml
	utils.WriteToFile(path.Join(apppath, ".gitlab-ci.yml"), gitlabCi)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, ".gitlab-ci.yml")+string(path.Separator), "\x1b[0m")

	//创建README.md
	utils.WriteToFile(path.Join(apppath, "README.md"), readMe)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "README.md")+string(path.Separator), "\x1b[0m")

	//创建gip.yml
	utils.WriteToFile(path.Join(apppath, "gip.yml"), gip)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "gip.yml")+string(path.Separator), "\x1b[0m")

	//创建gometalinter.json
	utils.WriteToFile(path.Join(apppath, "gometalinter.json"), gometalinter)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "gometalinter.json")+string(path.Separator), "\x1b[0m")

	//创建requirements.txt
	utils.WriteToFile(path.Join(apppath, "requirements.txt"), requirement)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "requirements.txt")+string(path.Separator), "\x1b[0m")

	//创建validate.sh
	utils.WriteToFile(path.Join(apppath, "validate.sh"), validate)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "validate.sh")+string(path.Separator), "\x1b[0m")

	logger.Log.Success("New application successfully created!")
	return 0
}