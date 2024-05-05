package upgrade

import (
	"fmt"
	"github.com/androidjp/kirby/cmd/kirby/internal/base"
	"github.com/spf13/cobra"
)

var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the kirby tools",
	Long:  "Upgrade the kirby tools. Example: kirby upgrade",
	Run:   Run,
}

func Run(_ *cobra.Command, _ []string) {
	err := base.GoInstall(
		"github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest",
		"github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest",
		"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
		"github.com/google/gnostic/cmd/protoc-gen-openapi@latest",
		"github.com/envoyproxy/protoc-gen-validate@latest",
		"github.com/favadi/protoc-go-inject-tag@latest",
	)
	if err != nil {
		fmt.Println(err)
	}
}
