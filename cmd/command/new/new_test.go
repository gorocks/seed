package new

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"
)

func TestCreateApp(t *testing.T) {
	//err := os.MkdirAll("/Users/luan/go/src/github.com/Guazi-inc/seed/path/to/dir", os.FileMode(0755))
	//t.Log(err)
	fmt.Printf("%q\n", strings.Split("/home/m_ta/src", "/"))

}
func TestCreateFile(t *testing.T) {
	str := `#path=/men/med.yml
asdfasfasdf
`
	reg, _ := regexp.Compile(`^#path=.+\n`)
	loc := reg.FindIndex([]byte(str))
	if len(loc) > 0 { //存在自定义path
		realPath := strings.Split(str, "=")[1]
		log.Println(realPath)
	} else {
		realPath := strings.Split(strings.Split(tempPath, template)[1], ".template")[0]
		careateFile(cmd, appPath, realPath, string(data))
	}

}
