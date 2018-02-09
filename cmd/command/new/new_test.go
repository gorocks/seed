package new

import (
	"bytes"
	"log"
	"testing"
	tmp "text/template"

	"strings"

	"github.com/Guazi-inc/seed/cmd/command/generator/proto"
	"github.com/gin-gonic/gin/json"
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
	protoP := "/Users/luan/go/src/protobuf-schema/proto/finance/service/repay/repay_offline_apply.proto"
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