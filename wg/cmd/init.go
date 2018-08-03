// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"net"
	"time"

	"github.com/sdeoras/wg/wg/proto"
	"github.com/sdeoras/wg/wg/server"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a new wait group",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		//start grpc server on localhost
		lis, err := net.Listen("tcp", "localhost:7001")
		if err != nil {
			return err
		}

		globalCtx, cancel := context.WithTimeout(context.Background(), time.Hour)
		defer cancel()

		srv, ctx := server.New(context.Background())

		s := grpc.NewServer()
		proto.RegisterExecServer(s, srv)
		reflection.Register(s)

		errChan := make(chan error)

		// start server in a goroutine and wait
		go func() {
			err := s.Serve(lis)
			if err != nil {
				errChan <- err
			}
		}()

		// wait till server implementor cancels the context it returned
		// or server returns an error, whichever happens first
		select {
		case <-globalCtx.Done():
		case <-ctx.Done():
			time.Sleep(time.Millisecond * 250)
		case err := <-errChan:
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
