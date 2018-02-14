package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"unicode"

	"github.com/Guazi-inc/seed/logger"
	"github.com/Guazi-inc/seed/logger/color"
)

// IsExist returns whether a file or directory exists.
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

//CheckGip returns whether gip exists.
func CheckGip() bool {
	cmd := exec.Command("gip", "version")
	bout := bytes.NewBuffer(nil)
	cmd.Stdout = bout
	if err := cmd.Run(); err != nil {
		logger.Warnf("use gip can read https://github.com/caojia/gip")
		logger.Warnf("Error while running gip: %s", err)
		return false
	}
	logger.Infof("gip local version is:%v", bout.String())
	return true
}

//CheckProtoc returns whether protoc exists.
func CheckProtoc() bool {
	cmd := exec.Command("protoc", "--version")
	bout := bytes.NewBuffer(nil)
	cmd.Stdout = bout
	if err := cmd.Run(); err != nil {
		logger.Warnf("protoc can read https://github.com/golang/protobuf")
		logger.Warnf("Error while running protoc: %s", err)
		return false
	}
	logger.Infof("protoc local version is:%v", bout.String())
	return true
}

func ProtocGenGo(protoPath, outPath, iPath string) {
	cmd := exec.Command("protoc", protoPath, fmt.Sprintf("--go_out=plugins=grpc:%v", outPath), "-I", iPath)
	if err := cmd.Run(); err != nil {
		logger.Fatalf("Error while running protoc: %s", err)
	}
	logger.Successf("generator file %s", protoPath)
}

//DoGipInstall gip install from requirements.txt
func DoGipInstall(requirementsPath string) {
	cmd := exec.Command("gip", "install", "-v", requirementsPath)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		logger.Warnf("use gip can read https://github.com/caojia/gip")
		logger.Warnf("Error while running gip: %s", err)
	}
	logger.Info("do gip install success")
}

func InstallAndUpdateGometalinter() bool {
	cmd1 := exec.Command("go", "get", "-u", "github.com/alecthomas/gometalinter")
	if err := cmd1.Run(); err != nil {
		logger.Warnf("use gometalinter  can read https://github.com/alecthomas/gometalinter")
		logger.Warnf("Error while running go get -u : %s", err)
		return false
	}
	cmd2 := exec.Command("gometalinter", " --install")
	if err := cmd2.Run(); err != nil {
		return false
	}
	return true
}

//do gometalinter from json
func DoGometalinterFromJson(filePath string) {
	cmd := exec.Command("gometalinter", "./...", fmt.Sprintf("--config=%s", filePath))
	cmd.Stdout = os.Stdout
	logger.Info("begin gometalinter")
	if err := cmd.Run(); err != nil {
		logger.Info("use gometalinter  can read https://github.com/alecthomas/gometalinter")
		logger.Fatalf("Error while running gometalinter: %s", err)
	}
	logger.Success("gometalinter success")
}

//do gometalinter from ci
func DoGometalinterCI() {
	cmd := exec.Command("gometalinter", "./...", "--cyclo-over=15", "--enable=deadcode", "--enable=errcheck", "--enable=gas", "--enable=goconst", "--enable=gocyclo", "--enable=golint", "--enable=gotype", "--enable=ineffassign", "--enable=interfacer", "--enable=megacheck", "--enable=structcheck", "--enable=unconvert", "--enable=varcheck", "--enable=vet", "--enable=vetshadow", "--enable=gofmt", "--enable=goimports", "--enable=unparam", "--enable=misspell", "--deadline=180s", "--concurrency=4")
	cmd.Stdout = os.Stdout
	logger.Info("begin gometalinter")
	if err := cmd.Run(); err != nil {
		logger.Info("use gometalinter  can read https://github.com/alecthomas/gometalinter")
		logger.Fatalf("Error while running gometalinter: %s", err)
	}
	logger.Success("gometalinter success")
}

// GetGOPATHs returns all paths in GOPATH variable.
func GetGOPATHs() []string {
	goPath := os.Getenv("GOPATH")
	if goPath == "" && strings.Compare(runtime.Version(), "go1.8") >= 0 {
		goPath = defaultGOPATH()
	}
	return filepath.SplitList(goPath)
}

// IsInGOPATH checks whether the path is inside of any GOPATH or not
func IsInGOPATH(thePath string) bool {
	for _, gopath := range GetGOPATHs() {
		if strings.Contains(thePath, filepath.Join(gopath, "src")) {
			return true
		}
	}
	return false
}

// SearchGOPATHs searchs the user GOPATH(s) for the specified application name.
// It returns a boolean, the application's GOPATH and its full path.
func SearchGOPATHs(app string) (bool, string, string) {
	gps := GetGOPATHs()
	if len(gps) == 0 {
		logger.Fatal("GOPATH environment variable is not set or empty")
	}

	// Lookup the application inside the user workspace(s)
	for _, gopath := range gps {
		var currentPath string

		if !strings.Contains(app, "src") {
			gopathsrc := path.Join(gopath, "src")
			currentPath = path.Join(gopathsrc, app)
		} else {
			currentPath = app
		}

		if IsExist(currentPath) {
			return true, gopath, currentPath
		}
	}
	return false, "", ""
}

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func AskForConfirmation() string {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		logger.Fatalf("%s", err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES", "overwrite"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	skipResponses := []string{"s", "skip", "sk"}
	allResponses := []string{"all", "skip all", "sa", "a"}
	ovResponses := []string{"owa", "overwrite all"}

	if containsString(okayResponses, response) {
		return "yes"
	} else if containsString(nokayResponses, response) {
		return "no"
	} else if containsString(skipResponses, response) {
		return "skip"
	} else if containsString(allResponses, response) {
		return "skip all"
	} else if containsString(ovResponses, response) {
		return "overwrite all"
	} else {
		fmt.Println("Please right type  then press enter:")
		return AskForConfirmation()
	}
}

func containsString(slice []string, element string) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	return false
}

// snake string, XxYy to xx_yy
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if !k && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || !k) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// camelCase converts a _ delimited string to camel case
// e.g. very_important_person => VeryImportantPerson
func CamelCase(in string) string {
	tokens := strings.Split(in, "_")
	for i := range tokens {
		tokens[i] = strings.Title(strings.Trim(tokens[i], " "))
	}
	return strings.Join(tokens, "")
}

// formatSourceCode formats source files
func FormatSourceCode(filename string) {
	cmd := exec.Command("gofmt", "-w", filename)
	if err := cmd.Run(); err != nil {
		logger.Warnf("Error while running gofmt: %s", err)
	}
}

// __FILE__ returns the file name in which the function was invoked
func FILE() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}

// __LINE__ returns the line number at which the function was invoked
func LINE() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}

func SeedFuncMap() template.FuncMap {
	return template.FuncMap{
		"trim":       strings.TrimSpace,
		"bold":       colors.Bold,
		"headline":   colors.MagentaBold,
		"foldername": colors.RedBold,
		"endline":    EndLine,
		"tmpltostr":  TmplToString,
	}
}
func EndLine() string {
	return "\n"
}

// TmplToString parses a text template and return the result as a string.
func TmplToString(tmpl string, data interface{}) string {

	t := template.New("tmpl").Funcs(SeedFuncMap())
	template.Must(t.Parse(tmpl))

	var doc bytes.Buffer
	err := t.Execute(&doc, data)
	if err != nil {
		logger.Fatalf("Error while TmplToString: %s", err)
	}

	return doc.String()
}

func Tmpl(text string, data interface{}) {
	output := colors.NewColorWriter(os.Stderr)

	t := template.New("Usage").Funcs(SeedFuncMap())
	template.Must(t.Parse(text))

	err := t.Execute(output, data)
	if err != nil {
		logger.Error(err.Error())
	}
}

func CheckEnv(appname string) (apppath, packpath string, err error) {
	gps := GetGOPATHs()
	if len(gps) == 0 {
		logger.Fatal("GOPATH environment variable is not set or empty")
	}
	currpath, _ := os.Getwd()
	currpath = filepath.Join(currpath, appname)
	for _, gpath := range gps {
		gsrcpath := filepath.Join(gpath, "src")
		if strings.HasPrefix(strings.ToLower(currpath), strings.ToLower(gsrcpath)) {
			packpath = strings.Replace(currpath[len(gsrcpath)+1:], string(filepath.Separator), "/", -1)
			return currpath, packpath, nil
		}
	}

	// In case of multiple paths in the GOPATH, by default
	// we use the first path
	gopath := gps[0]

	logger.Warn("You current workdir is not inside $GOPATH/src.")

	gosrcpath := filepath.Join(gopath, "src")
	apppath = filepath.Join(gosrcpath, appname)

	if _, e := os.Stat(apppath); !os.IsNotExist(e) {
		err = fmt.Errorf("Cannot create application without removing '%s' first", apppath)
		logger.Errorf("Path '%s' already exists", apppath)
		return
	}
	packpath = strings.Join(strings.Split(apppath[len(gosrcpath)+1:], string(filepath.Separator)), "/")
	return
}

func PrintErrorAndExit(message, errorTemplate string) {
	Tmpl(fmt.Sprintf(errorTemplate, message), nil)
	os.Exit(2)
}

// GoCommand executes the passed command using Go tool
func GoCommand(command string, args ...string) error {
	allargs := []string{command}
	allargs = append(allargs, args...)
	goBuild := exec.Command("go", allargs...)
	goBuild.Stderr = os.Stderr
	return goBuild.Run()
}

func Json2Str(jsonData string) string {
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

// SplitQuotedFields is like strings.Fields but ignores spaces
// inside areas surrounded by single quotes.
// To specify a single quote use backslash to escape it: '\''
func SplitQuotedFields(in string) []string {
	type stateEnum int
	const (
		inSpace stateEnum = iota
		inField
		inQuote
		inQuoteEscaped
	)
	state := inSpace
	r := []string{}
	var buf bytes.Buffer

	for _, ch := range in {
		switch state {
		case inSpace:
			if ch == '\'' {
				state = inQuote
			} else if !unicode.IsSpace(ch) {
				buf.WriteRune(ch)
				state = inField
			}

		case inField:
			if ch == '\'' {
				state = inQuote
			} else if unicode.IsSpace(ch) {
				r = append(r, buf.String())
				buf.Reset()
			} else {
				buf.WriteRune(ch)
			}

		case inQuote:
			if ch == '\'' {
				state = inField
			} else if ch == '\\' {
				state = inQuoteEscaped
			} else {
				buf.WriteRune(ch)
			}

		case inQuoteEscaped:
			buf.WriteRune(ch)
			state = inQuote
		}
	}

	if buf.Len() != 0 {
		r = append(r, buf.String())
	}

	return r
}

func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		return filepath.Join(home, "go")
	}
	return ""
}

func GetUsefulPath(path string, field string, isContaint bool) string {
	arr := strings.Split(path, "/")
	for k, v := range arr {
		if v == field {
			if !isContaint {
				return strings.Join(arr[k+1:len(arr)-1], "/")
			}
			return strings.Join(arr[k:len(arr)-1], "/")
		}
	}
	return ""
}

func RmDuplicate(list []string) []string {
	x := make([]string, 0)
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}
