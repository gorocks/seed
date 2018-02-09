package proto

import (
	"errors"
	"os"
	"strings"

	"github.com/emicklei/proto"
)

const (
	GOOGLEEMPTY = "github.com/golang/protobuf/ptypes/empty"
)

type GeneratorProto struct {
	Package string      `json:"package"`
	Service []*GService `json:"service"`
}
type GService struct {
	ServiceName string   `json:"service_name"`
	Imports     []string `json:"imports"` //这个service下用到的import
	Rpc         []*GFunc `json:"rpc"`
}

type GFunc struct {
	FunName  string `json:"fun_name"`
	Request  string `json:"request"`
	Response string `json:"response"`
}

//简单的解析proto,只是按照字符串解析，为自己想要的结构，而不是打包为go文件
func (g *GeneratorProto) Generator(filePath string) error {

	if !strings.HasSuffix(filePath, ".proto") {
		return errors.New("invalid file name")
	}
	//匹配proto里service 中rpc的内容
	r, err := os.Open(filePath)
	if err != nil {
		return err
	}
	p, err := proto.NewParser(r).Parse()
	if err != nil {
		return err
	}
	var gservices []*GService
	importMap := make(map[string]string)

	for _, v := range p.Elements {
		switch value := v.(type) {
		case *proto.Package:
			g.Package = value.Name
		case *proto.Service:
			gservive := &GService{
				ServiceName: value.Name,
			}
			var gfuns []*GFunc
			ims := make(map[string]string)
			for _, j := range value.Elements {
				if rpc, ok := j.(*proto.RPC); ok {
					gfun := &GFunc{
						FunName:  rpc.Name,
						Request:  rpc.RequestType,
						Response: rpc.ReturnsType,
					}
					arr1 := strings.Split(rpc.ReturnsType, ".")
					if len(arr1) >= 2 {
						if _, ok := ims[arr1[0]]; !ok {
							ims[arr1[0]] = importMap[arr1[0]]
						}
					}
					arr2 := strings.Split(rpc.RequestType, ".")
					if len(arr2) >= 2 {
						if _, ok := ims[arr2[0]]; ok {
							ims[arr2[0]] = importMap[arr2[0]]
						}
					}
					gfuns = append(gfuns, gfun)
				}
			}
			gservive.Rpc = gfuns
			imp := make([]string, 0)
			for _, v := range ims {
				imp = append(imp, v)
			}
			gservive.Imports = imp
			gservices = append(gservices, gservive)
		case *proto.Import:
			path := value.Filename
			//处理path
			arr := strings.Split(path, "/")
			importMap[arr[len(arr)-2]] = strings.Join(arr[:len(arr)-1], "/")
		}
	}
	g.Service = gservices
	return nil
}
