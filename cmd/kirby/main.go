package main

import (
	"github.com/androidjp/kirby/cmd/kirby/internal/proto"
	"github.com/androidjp/kirby/cmd/kirby/internal/upgrade"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = cobra.Command{
	Use:                        "kirby",
	Aliases:                    []string{"kb"},
	SuggestFor:                 []string{"suggestFor"},
	Short:                      "Kirby: An elegant toolkit for Go microservices.[Short]",
	Long:                       "Kirby: An elegant toolkit for Go microservices.[Long]",
	Example:                    "example: xxxxxxxx",
	ValidArgs:                  []string{"a", "b"},
	ValidArgsFunction:          nil,
	Args:                       nil,
	ArgAliases:                 nil,
	BashCompletionFunction:     "",
	Deprecated:                 "",
	Annotations:                nil,
	Version:                    release,
	PersistentPreRun:           nil,
	PersistentPreRunE:          nil,
	PreRun:                     nil,
	PreRunE:                    nil,
	Run:                        nil,
	RunE:                       nil,
	PostRun:                    nil,
	PostRunE:                   nil,
	PersistentPostRun:          nil,
	PersistentPostRunE:         nil,
	FParseErrWhitelist:         cobra.FParseErrWhitelist{},
	CompletionOptions:          cobra.CompletionOptions{},
	TraverseChildren:           false,
	Hidden:                     false,
	SilenceErrors:              false,
	SilenceUsage:               false,
	DisableFlagParsing:         false,
	DisableAutoGenTag:          false,
	DisableFlagsInUseLine:      false,
	DisableSuggestions:         false,
	SuggestionsMinimumDistance: 0,
}

func init() {
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
