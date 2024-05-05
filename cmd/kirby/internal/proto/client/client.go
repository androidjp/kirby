package client

import (
	"fmt"
	"github.com/androidjp/kirby/cmd/kirby/internal/base"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// CmdClient represents the source command.
var CmdClient = &cobra.Command{
	Use:   "client",
	Short: "Generate the proto client code",
	Long:  "Generate the proto client code. Example: kirby proto client helloworld.proto",
	Run:   run,
}

var protoPath string

func init() {
	if protoPath = os.Getenv("KIRBY_PROTO_PATH"); protoPath == "" {
		protoPath = "./third_party"
	}
	CmdClient.Flags().StringVarP(&protoPath, "proto_path", "p", protoPath, "proto path")
}

func run(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please enter the protoFileName file or dir")
		return
	}
	var (
		err           error
		protoFileName = strings.TrimSpace(args[0])
	)
	// 找到这些可执行文件，如果没有，则执行upgrade
	if err = look("protoc-gen-go", "protoc-gen-go-grpc", "protoc-gen-go-http", "protoc-gen-go-errors", "protoc-gen-openapi", "protoc-gen-validate", "protoc-go-inject-tag"); err != nil {
		fmt.Println("look err:", err.Error())
		cmd := exec.Command("kirby", "upgrade")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			fmt.Println(err)
			return
		}
	}

	// 有.proto文件，则生成对应文件
	if strings.HasSuffix(protoFileName, ".proto") {
		err = generate(protoFileName, args)
	} else {
		err = walk(protoFileName, args)
	}
	if err != nil {
		fmt.Println(err)
	}
}

func walk(dir string, args []string) error {
	if dir == "" {
		dir = "."
	}
	// 遍历目录下的所有.proto文件，忽略掉 third_party目录 以及 非.proto文件
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if ext := filepath.Ext(path); ext != ".proto" || strings.HasPrefix(path, "third_party") {
			return nil
		}
		return generate(path, args)
	})
}

func generate(proto string, args []string) error {
	// 1. 生成protoc语句
	input := []string{
		"--proto_path=.",
	}
	if pathExists(protoPath) {
		input = append(input, "--proto_path="+protoPath)
	}
	inputExt := []string{
		"--proto_path=" + base.KirbyMod(),
		"--proto_path=" + filepath.Join(base.KirbyMod(), "third_party"),
		"--go_out=paths=source_relative:.",
		"--go-grpc_out=paths=source_relative:.",
		"--go-http_out=paths=source_relative:.",
		"--go-errors_out=paths=source_relative:.",
		"--openapi_out=paths=source_relative:.",
	}
	input = append(input, inputExt...)

	var needInjectTag bool
	protoBytes, err := os.ReadFile(proto)
	if err == nil && len(protoBytes) > 0 {
		// 如果proto文件中包含import "validate/validate.proto"，则添加--validate_out参数
		if ok, _ := regexp.Match(`\n[^/]*(import)\s+"validate/validate.proto"`, protoBytes); ok {
			input = append(input, "--validate_out=lang=go,paths=source_relative:.")
		}
		if ok, _ := regexp.Match(`@inject_tag`, protoBytes); ok {
			needInjectTag = true
		}
	}

	input = append(input, proto)
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			input = append(input, a)
		}
	}

	// 2. 执行protoc语句
	fd := exec.Command("protoc", input...)
	fmt.Println("背后执行指令：", fd.String())
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = "."
	if err = fd.Run(); err != nil {
		return err
	}
	// 3. 如果发现有 inject-tag，则顺便执行inject-tag
	if needInjectTag {
		injectFD := exec.Command("protoc-go-inject-tag", "-input="+strings.ReplaceAll(proto, ".proto", "")+".pb.go")
		injectFD.Stdout = os.Stdout
		injectFD.Stderr = os.Stderr
		injectFD.Dir = "."
		if err = injectFD.Run(); err != nil {
			return err
		}
	}
	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func look(name ...string) error {
	// 去环境变量Path的目录中查找对应可执行文件
	for _, n := range name {
		if _, err := exec.LookPath(n); err != nil {
			return err
		}
	}
	return nil
}
