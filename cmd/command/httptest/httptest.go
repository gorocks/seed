package httptest

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Guazi-inc/seed/cmd/command"
	"github.com/Guazi-inc/seed/cmd/command/version"
	"github.com/Guazi-inc/seed/logger"
	"github.com/Guazi-inc/seed/utils"
)

var CmdHttptest = &commands.Command{
	UsageLine: "httptest -pkg=[test-fixtures path]",
	Short:     "set up a http server for test",
	Long: `
Run http server fot test,this server will supervise the filesystem of the application for any changes, and recompile/restart it.
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    RunHttptest,
}

type HttpJsonMock struct {
	jsonResponse map[string]interface{}
	port         string
	fxituresPath string
	style        string
}

var hm HttpJsonMock

func init() {
	fs := flag.NewFlagSet("httptest", flag.ContinueOnError)
	fs.StringVar(&hm.fxituresPath, "pkg", "./test-fixtures", "test-fixtures path")
	fs.StringVar(&hm.port, "p", "8090", "local http server test port")
	fs.StringVar(&hm.style, "style", "json", "fixtures style")
	CmdHttptest.Flag = *fs
	commands.AvailableCommands = append(commands.AvailableCommands, CmdHttptest)
}

func RunHttptest(cmd *commands.Command, args []string) int {
	if err := cmd.Flag.Parse(args); err != nil {
		logger.Fatalf("Error while parsing flags: %v", err.Error())
	}
	filePath := hm.fxituresPath + "/" + hm.style
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		logger.Fatalf("please use httptest -pkg to set fixtures path to fix err :%v", err)
	}
	hm.jsonResponse = make(map[string]interface{})
	for _, v := range files {
		ret, err := readFileToMap(filePath + "/" + v.Name())
		if err != nil {
			logger.Errorf("%s %v+", v.Name(), err)
			continue
		}
		hm.jsonResponse[v.Name()] = ret
	}
	utils.NewWatcher([]string{filePath}, []string{}, &hm)
	hm.RunServer()
	return 0
}

func (hm *HttpJsonMock) RunServer() {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	ports := fmt.Sprintf(":%s", hm.port)
	srv := http.Server{Addr: ports}
	http.HandleFunc("/", hm.Handle)
	go func() {
		logger.Infof("http test server start at %s", ports)
		if err := srv.ListenAndServe(); err != nil {
			logger.Infof("server listen: %s\n", err)
		}
	}()
	<-stopChan
	logger.Info("Shutting down server...")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	logger.Info("Server gracefully stopped")
}

func (hm *HttpJsonMock) Exec(paths []string, files []string, name string) {
	temp, err := readFileToMap(name)
	if err != nil {
		logger.Errorf("exec err %v", err)
	}
	arr := strings.Split(name, "/")
	str := arr[len(arr)-1]
	hm.jsonResponse[str] = temp
	logger.Infof("%s change and  save success", str)
}

//对请求进行处理
func (hm *HttpJsonMock) Handle(w http.ResponseWriter, r *http.Request) {
	var res map[string]interface{}
	host, urlP := splitPath(r.URL.Path)
	logger.Infof("request path : %s", r.Method, host, urlP)
	fileName := host + ".json"
	if temp, ok := hm.jsonResponse[fileName]; ok {
		res = temp.(map[string]interface{})
	} else {
		logger.Errorf("no %s exist", host)
		return
	}
	//读取request 中的参数,判断请求的方式
	arr := ""
	switch r.Method {
	case "GET":
		//对get请求中的参数进行匹配
		pathArr := strings.Split(r.RequestURI, "?")
		if len(pathArr) > 1 {
			arr = pathArr[1]
		}
	case "POST":
		//对post请求进行处理
		date, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Errorf("err: %v", err)
		}
		arr = string(date)
		if strings.Contains(arr, "{") { //认为是json请求，就拼凑请求中的内容
			arr = utils.Json2Str(arr)
		}
	default:
		w.Write([]byte("unsupported " + r.Method + " method"))
	}
	sli := strings.Split(arr, "&")
	key := make([]string, 0)
	for _, v := range sli { //对请求体中的参数进行过滤
		if !(strings.Contains(v, "signature") || strings.Contains(v, "time") || strings.Contains(v, "appkey") || strings.Contains(v, "expires") || strings.Contains(v, "nonce")) {
			key = append(key, v)
		}
	}
	keys := strings.Replace(strings.Join(key, "&"), " ", "", -1)
	logger.Infof("key: %v", fmt.Sprintf("%v", keys))
	//对fixtures中的参数进行匹配，并返回对应的response
	if temp, ok := res[urlP]; ok { //判断时候是否是*
		jsonToResponse(w, temp.(map[string]interface{}), keys)
	} else {
		//
		str := strings.Replace(urlP, "/", "", 1)
		if temp, ok = res[str]; ok {
			jsonToResponse(w, temp.(map[string]interface{}), keys)
		} else {
			logger.Errorf("no match url ,url should is : %s", urlP)
		}

	}
}

//解析json文件
func jsonToResponse(w http.ResponseWriter, arr map[string]interface{}, keys string) {
	//遍历匹配key
	for k, v := range arr {
		isKey := true
		if k != keys { //如果不能直接匹配就判断包含关系
			arr := strings.Split(k, "&")
			for _, v := range arr {
				if !strings.Contains(keys, v) {
					isKey = false
					break
				}
			}
		}
		if isKey {
			data, err := httpCommonRet(v)
			if err != nil {
				panic(err)
			}
			logger.Infof("response is : %v", string(data))
			w.Write(data)
			return
		}
	}
	logger.Errorf("no match key,key should is : %v", keys)
}

func readFileToMap(fileName string) (map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Errorf("read file err: %v", err)
		return nil, err
	}

	var ret map[string]interface{}
	if err := json.Unmarshal(bytes, &ret); err != nil {
		logger.Errorf("read file json unmarshal err: %v", err)
		return nil, err
	}

	return ret, nil
}

func splitPath(path string) (string, string) {
	arr := strings.Split(path, "/")[1:]
	s := ""
	p := ""
	for k, v := range arr {
		if k == 0 {
			s = v
		} else {
			p = p + "/" + v
		}
	}
	return s, p
}

func httpCommonRet(s interface{}) ([]byte, error) {
	date := map[string]interface{}{
		"code":    0,
		"message": "succeed",
		"data":    s,
	}
	data, err := json.Marshal(date)
	if err != nil {
		return nil, err
	}
	return data, nil
}
