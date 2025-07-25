package main

import (
	goflag "flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	utilflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"

	"open-cluster-management.io/ocm/pkg/cmd/hub"
	"open-cluster-management.io/ocm/pkg/version"
)

func main() {
	pflag.CommandLine.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	logs.AddFlags(pflag.CommandLine)
	logs.InitLogs()
	defer logs.FlushLogs()

	command := newServerCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1) //nolint:gocritic
	}
}

// newServerCommand creates a new command for the server.
// The server command is responsible for starting the server on the hub.
// The server serves the resources transport between the hub and the managed clusters.
func newServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Server for transport resources",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(1)
		},
	}

	if v := version.Get().String(); len(v) == 0 {
		cmd.Version = "<unknown>"
	} else {
		cmd.Version = v
	}

	cmd.AddCommand(hub.NewGRPCServerCommand())

	return cmd
}
