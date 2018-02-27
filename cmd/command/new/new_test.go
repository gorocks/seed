package new

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strings"
	"testing"
	tmp "text/template"

	"github.com/Guazi-inc/seed/cmd/command/generator/proto"
)

func TestCreateFile(t *testing.T) {
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
	protoP := "/Users/luan/go/src/protobuf-schema/proto/finance/service/borrow/borrow.proto"
	g := proto.GeneratorProto{}
	err := g.Generator(protoP)
	if err != nil {
		log.Println(err)
		return
	}
	arr := strings.Split(protoP, "/")

	type sr struct {
		ServiceName string         `json:"service_name"`
		Imports     []string       `json:"imports"` //这个service下用到的import
		Rpc         []*proto.GFunc `json:"rpc"`
		Package     string         `json:"package"`
		PackageName string         `json:"package_name"`
		PackPath    string         `json:"pack_path"`
	}
	tp := tmp.New("Service") //创建一个模板
	tp.Funcs(tmp.FuncMap{
		"tmp": ServiceTemplPath,
	})
	p, err := tp.Parse(serviceTmpl)
	if err != nil {
		panic(err)
	}
	for _, v := range g.Service {
		a, _ := json.Marshal(v.Imports)
		log.Println(string(a))
		s := sr{
			ServiceName: v.ServiceName,
			Imports:     v.Imports,
			Rpc:         v.Rpc,
			Package:     g.Package,
			PackageName: strings.ToLower(v.ServiceName),
		}
		for k, v := range arr {
			if v == "proto" {
				s.PackPath = strings.Join(arr[k:len(arr)-1], "/")
			}
		}
		var content bytes.Buffer
		err = p.Execute(&content, s)
		if err != nil {
			panic(err)
		}
		//建立某一个service的文件
		log.Println(content.String())
	}
}

func Test_isNetZip(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"t1", args{[]byte("https://github.com/Guazi-inc/seed/archive/master.zip")}, true},
		{"t2", args{[]byte("master.zip")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNetZip(tt.args.b); got != tt.want {
				t.Errorf("isNetZip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseZip(t *testing.T) {
	tf := tempFileName()
	type args struct {
		url    string
		toPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"t1", args{"https://github.com/Guazi-inc/seed/archive/master.zip", tf}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parseZip(tt.args.url, tt.args.toPath); (err != nil) != tt.wantErr {
				t.Errorf("parseZip() error = %v, wantErr %v", err, tt.wantErr)
			}
			defer os.RemoveAll(tf)
		})
	}
}
