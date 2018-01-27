package httptest

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github/Guazi-inc/seed/cmd/command"
	"github/Guazi-inc/seed/cmd/command/version"
	"github/Guazi-inc/seed/logger"
	"github/Guazi-inc/seed/utils"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

var jsonDataMap map[string]interface{}

var CmdHttptest = &commands.Command{
	UsageLine: "httptest",
	Short:     "set up a http server for test",
	Long: `
Run http server fot test,this server will supervise the filesystem of the application for any changes, and recompile/restart it.
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    RunHttptest,
}

var (
	port  string
	path  string
	style string
)

func init() {
	fs := flag.NewFlagSet("httptest", flag.ContinueOnError)
	fs.StringVar(&path, "pkg", "./test-fixtures", "test-fixtures path")
	fs.StringVar(&port, "p", "8090", "local http server test port")
	fs.StringVar(&style, "style", "json", "fixtures style")
	CmdHttptest.Flag = *fs
	commands.AvailableCommands = append(commands.AvailableCommands, CmdHttptest)
}

func RunHttptest(cmd *commands.Command, args []string) int {
	if err := cmd.Flag.Parse(args); err != nil {
		logger.Log.Fatalf("Error while parsing flags: %v", err.Error())
	}
	filePath := path + "/" + style
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		logger.Log.Errorf("please use httptest -pkg to set fixtures path to fix err :%v", err)
		return 1
	}
	jsonDataMap = make(map[string]interface{}, len(files))
	for _, v := range files {
		ret, err := readFile(filePath + "/" + v.Name())
		if err != nil {
			logger.Log.Errorf("%s %v+", v.Name(), err)
			continue
		}
		jsonDataMap[v.Name()] = ret
	}
	var wt watcher
	utils.NewWatcher([]string{filePath}, []string{}, wt)
	//读取配置文件，获取端口号
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	ports := fmt.Sprintf(":%s", port)
	srv := http.Server{Addr: ports}
	http.HandleFunc("/", Handle)
	go func() {
		logger.Log.Infof("http test server start at %s", ports)
		if err := srv.ListenAndServe(); err != nil {
			logger.Log.Infof("server listen: %s\n", err)
		}
	}()
	<-stopChan // wait for SIGINT
	logger.Log.Info("Shutting down server...")
	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	logger.Log.Info("Server gracefully stopped")
	return 0
}

type watcher string

func (wt watcher) Exec(paths []string, files []string, name string) {
	temp, err := readFile(name)
	if err != nil {
		logger.Log.Errorf("exec err %v", err)
	}
	arr := strings.Split(name, "/")
	str := arr[len(arr)-1]
	jsonDataMap[str] = temp
	logger.Log.Infof("%s change and  save success", str)
}

//对请求进行处理
func Handle(w http.ResponseWriter, r *http.Request) {
	var res map[string]interface{}
	host, urlP := splitPath(r.URL.Path)
	logger.Log.Infof("request path : %s", r.Method, host, urlP)
	fileName := host + ".json"
	if temp, ok := jsonDataMap[fileName]; ok {
		res = temp.(map[string]interface{})
	} else {
		logger.Log.Errorf("no %s exist", host)
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
			logger.Log.Errorf("err: %v", err)
		}
		arr = string(date)
		if strings.Contains(arr, "{") { //认为是json请求，就拼凑请求中的内容
			arr = json2Str(arr)
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
	logger.Log.Infof("key: %v", fmt.Sprintf("%v", keys))
	//对fixtures中的参数进行匹配，并返回对应的response
	if temp, ok := res[urlP]; ok { //判断时候是否是*
		jsonToResponse(w, temp.(map[string]interface{}), keys)
	} else {
		//
		str := strings.Replace(urlP, "/", "", 1)
		if temp, ok = res[str]; ok {
			jsonToResponse(w, temp.(map[string]interface{}), keys)
		} else {
			logger.Log.Errorf("no match url ,url should is : %s", urlP)
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
			data, err := retMarshal(v)
			if err != nil {
				panic(err)
			}
			logger.Log.Infof("response is : %v", string(data))
			w.Write(data)
			return
		}
	}
	logger.Log.Errorf("no match key,key should is : %v", keys)
}

func readFile(fileName string) (map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Log.Errorf("read file err: %v", err)
		return nil, err
	}

	var ret map[string]interface{}
	if err := json.Unmarshal(bytes, &ret); err != nil {
		logger.Log.Errorf("read file json unmarshal err: %v", err)
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

func json2Str(jsonData string) string {
	arr := make([]rune, 0)
	isT := false
	for _, v := range jsonData {
		if v == int32(34) || v == int32(123) || v == int32(125) {
			continue
		}
		if v == int32(58) {
			v = int32(61)
		}
		if v == int32(44) {
			v = int32(38)
		}
		if v == int32(92) {
			isT = true
		}
		arr = append(arr, v)
	}
	str := string(arr)
	//处理转义
	if isT {
		str = strings.Replace(str, `\u003c`, "<", -1)
		str = strings.Replace(str, `\u003e`, ">", -1)
		str = strings.Replace(str, `\u0026`, "&", -1)
	}
	return str
}

func retMarshal(s interface{}) ([]byte, error) {
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
