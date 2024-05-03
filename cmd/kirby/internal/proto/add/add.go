package add

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
)

var CmdAdd = &cobra.Command{
	Use:   "add",
	Short: "Add a proto API template",
	Long:  "Add a proto API template. Example: kratos proto add helloworld/v1/hello.proto",
	Run:   run,
}

func run(_ *cobra.Command, args []string) {
	// TODO:
	// 1. 解析参数
	if len(args) == 0 {
		fmt.Println("Please enter the proto file or directory")
		return
	}
	input := args[0]
	n := strings.LastIndex(input, "/")
	if n == -1 {
		fmt.Println("The proto path needs to be hierarchical.")
		return
	}
	path := input[:n]
	fileName := input[n+1:]
	pkgName := strings.ReplaceAll(path, "/", ".")

	fmt.Println("--------------------------------")
	fmt.Printf("包名: %s\n", pkgName)
	fmt.Printf("proto文件名: %s\n", fileName)
	fmt.Printf("proto文件名: %s\n", fileName)
	fmt.Println("--------------------------------")

	// 2. 生成模板
	p := &Proto{
		Name:        fileName,
		Path:        path,
		Package:     pkgName,
		GoPackage:   goPackage(path),
		JavaPackage: javaPackage(pkgName),
		Service:     serviceName(fileName),
	}
	if err := p.Generate(); err != nil {
		fmt.Println(err)
		return
	}
}

func serviceName(name string) string {
	return toUpperCamelCase(strings.Split(name, ".")[0])
}

func toUpperCamelCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = cases.Title(language.Und, cases.NoLower).String(s)
	return strings.ReplaceAll(s, " ", "")
}

func javaPackage(name string) string {
	return name
}

func goPackage(path string) string {
	s := strings.Split(path, "/")
	mName := modName()
	fmt.Println("modName: ", mName)
	return modName() + "/" + path + ";" + s[len(s)-1]
}

func modName() string {
	// 依次尝试读取：当前上下文的go.mod、上一级目录的go.mod
	modBytes, err := os.ReadFile("go.mod")
	if err != nil {
		if modBytes, err = os.ReadFile("../go.mod"); err != nil {
			return ""
		}
	}
	return modfile.ModulePath(modBytes)
}
