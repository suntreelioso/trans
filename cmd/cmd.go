package cmd

import (
	"fmt"
	"log"
	"net"
	"trans/client"
	"trans/common"
	"trans/server"

	"github.com/spf13/cobra"
)

const version = "0.0.1"

var Cmd *cobra.Command

func init() {
	rootCmd := &cobra.Command{
		Version:           version,
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
		Use:               "trans",
		Long:              "trans - a simple file transfer tool",
		Run: func(cmd *cobra.Command, args []string) {
			addr := parseAddrFlag(cmd)
			path := parsePathFlag(cmd)
			filenames, err := cmd.Flags().GetStringSlice("get")
			getAll, _ := cmd.Flags().GetBool("get-all")
			if len(filenames) > 0 {
				if err != nil {
					common.ExitWithError(err)
				}
				client.DownloadFile(addr.String(), path, filenames)
			} else if getAll {
				client.DownloadAllFile(addr.String(), path)
			} else {
				client.ListFiles(addr.String())
			}
		},
	}

	{
		cmd := &cobra.Command{
			Use:   "server",
			Short: "server mode",
			Run: func(cmd *cobra.Command, args []string) {
				addr := parseAddrFlag(cmd)
				path := parsePathFlag(cmd)
				log.Printf("working directory: %v", path)
				server.StartServer(addr.String(), path)
			},
		}
		cmd.Flags().StringP("addr", "a", "0.0.0.0:8080", "listen address")
		rootCmd.AddCommand(cmd)
	}

	rootCmd.Flags().StringP("addr", "a", "127.0.0.1:8080", "server address")
	rootCmd.PersistentFlags().StringP("path", "p", "", "path to share, default: current working directory")
	rootCmd.Flags().BoolP("list", "l", true, "list all files")
	rootCmd.Flags().StringSliceP("get", "g", nil, "get one or more files, e.g. -g 1.txt,2.mp3")
	rootCmd.Flags().BoolP("get-all", "G", false, "get all files")
	Cmd = rootCmd
}

func parseAddrFlag(cmd *cobra.Command) *net.TCPAddr {
	addr, err := cmd.Flags().GetString("addr")
	if err != nil {
		common.ExitWithError(err)
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		common.ExitWithError(err)
	}
	return tcpAddr
}

func parsePathFlag(cmd *cobra.Command) string {
	path, _ := cmd.Flags().GetString("path")
	if path == "" {
		path = common.GetCwd()
	}
	if !common.DirIsExist(path) {
		common.ExitWithError(fmt.Errorf("the path does not exist or is not a directory: %v", path))
	}
	return path
}
