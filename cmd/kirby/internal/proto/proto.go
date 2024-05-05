package proto

import (
	"github.com/androidjp/kirby/cmd/kirby/internal/proto/add"
	"github.com/androidjp/kirby/cmd/kirby/internal/proto/client"
	"github.com/spf13/cobra"
)

var CmdProto = &cobra.Command{
	Use:   "proto",
	Short: "Generate the proto files",
	Long:  "Generate the proto files",
}

func init() {
	CmdProto.AddCommand(add.CmdAdd)
	CmdProto.AddCommand(client.CmdClient)
	//CmdProto.AddCommand(server.CmdServer)
}
