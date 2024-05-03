package add

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Proto struct {
	Name        string
	Path        string
	Service     string
	Package     string
	GoPackage   string
	JavaPackage string
}

func (p *Proto) Generate() error {
	// 1. 根据template，得到文本内容
	body, err := p.execute()
	if err != nil {
		return err
	}
	// 2. 获取当前工作目录
	wd, wdErr := os.Getwd()
	if wdErr != nil {
		panic(wdErr)
	}
	to := filepath.Join(wd, p.Path)

	// 3. 不存在则创建目录，并让这些目录的权限变为：当前用户可读写执行。
	if _, statErr := os.Stat(to); os.IsNotExist(statErr) {
		if mkdirErr := os.MkdirAll(to, 0o700); mkdirErr != nil {
			return mkdirErr
		}
	}

	// 4. 判断目标proto文件是否存在，如果存在则跳过，否则新创建proto文件
	// TODO 考虑加上 -force选项，让proto文件可以重新创建
	name := filepath.Join(to, p.Name)
	if _, statErr := os.Stat(name); !os.IsNotExist(statErr) {
		return fmt.Errorf("%s already exists", name)
	}
	return os.WriteFile(name, body, 0o644)
}

func (p *Proto) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("proto").Parse(strings.TrimSpace(protoTemplate))
	if err != nil {
		return nil, err
	}
	if eErr := tmpl.Execute(buf, p); eErr != nil {
		return nil, eErr
	}
	return buf.Bytes(), nil
}
