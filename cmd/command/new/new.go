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

            ├── main.go
            ├── {{"conf"|foldername}}
            │     └── app.conf
            ├── {{"controllers"|foldername}}
            │     └── default.go
            ├── {{"models"|foldername}}
            ├── {{"routers"|foldername}}
            │     └── router.go
            ├── {{"tests"|foldername}}
            │     └── default_test.go
            ├── {{"static"|foldername}}
            │     └── {{"js"|foldername}}
            │     └── {{"css"|foldername}}
            │     └── {{"img"|foldername}}
            └── {{"views"|foldername}}
                  └── index.tpl

`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    CreateApp,
}

var appconf = `appname = {{.Appname}}
httpport = 8080
runmode = dev
`

var maingo = `package main

import (
	_ "{{.Appname}}/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

`
var router = `package routers

import (
	"{{.Appname}}/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
`

var test = `package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"runtime"
	"path/filepath"
	_ "{{.Appname}}/routers"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}


// TestBeego is a sample to run an endpoint test
func TestBeego(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestBeego", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	                So(w.Code, ShouldEqual, 200)
	        })
	        Convey("The Result Should Not Be Empty", func() {
	                So(w.Body.Len(), ShouldBeGreaterThan, 0)
	        })
	})
}

`

var controllers = `package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
`

var reloadJsClient = `function b(a){var c=new WebSocket(a);c.onclose=function(){setTimeout(function(){b(a)},2E3)};c.onmessage=function(){location.reload()}}try{if(window.WebSocket)try{b("ws://localhost:12450/reload")}catch(a){console.error(a)}else console.log("Your browser does not support WebSockets.")}catch(a){console.error("Exception during connecting to Reload:",a)};
`

func init() {
	commands.AvailableCommands = append(commands.AvailableCommands, CmdNew)
}

func CreateApp(cmd *commands.Command, args []string) int {
	output := cmd.Out()
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

	os.MkdirAll(apppath, 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "conf"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "conf")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "controllers"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "controllers")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "models"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "models")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "routers"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "routers")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "tests"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "tests")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "static"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "static")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "static", "js"), 0755)
	utils.WriteToFile(path.Join(apppath, "static", "js", "reload.min.js"), reloadJsClient)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "static", "js")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "static", "css"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "static", "css")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "static", "img"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "static", "img")+string(path.Separator), "\x1b[0m")
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "views")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "views"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "conf", "app.conf"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "conf", "app.conf"), strings.Replace(appconf, "{{.Appname}}", path.Base(args[0]), -1))

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "controllers", "default.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "controllers", "default.go"), controllers)

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "views", "index.tpl"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "views", "index.tpl"), indextpl)

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "routers", "router.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "routers", "router.go"), strings.Replace(router, "{{.Appname}}", packpath, -1))

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "tests", "default_test.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "tests", "default_test.go"), strings.Replace(test, "{{.Appname}}", packpath, -1))

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "main.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "main.go"), strings.Replace(maingo, "{{.Appname}}", packpath, -1))

	logger.Log.Success("New application successfully created!")
	return 0
}
