// Copyright 2023 Northern.tech AS
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	    http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.
package cmd

import (
	"errors"
	"fmt"
	"github.com/mendersoftware/mender-cli/comms"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mendersoftware/mender-cli/client/deployments"
)

const ()

var commsArtifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Show the artifact with the latest comms version",
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewCommsArtifactCmd(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

func init() {
	// commsArtifactCmd.Flags().IntP(argDetailLevel, "d", 0, "artifacts list detail level [0..3]")
}

type CommsArtifactCmd struct {
	server     string
	skipVerify bool
	token      string
}

func NewCommsArtifactCmd(cmd *cobra.Command, args []string) (*CommsArtifactCmd, error) {
	server := viper.GetString(argRootServer)
	if server == "" {
		return nil, errors.New("No server")
	}

	skipVerify, err := cmd.Flags().GetBool(argRootSkipVerify)
	if err != nil {
		return nil, err
	}

	token, err := getAuthToken(cmd)
	if err != nil {
		return nil, err
	}

	return &CommsArtifactCmd{
		server:     server,
		token:      token,
		skipVerify: skipVerify,
	}, nil
}

func (c *CommsArtifactCmd) Run() error {

	client := deployments.NewClient(c.server, c.skipVerify)
	list, err := client.ListArtifacts(c.token)
	if err != nil {
		return err
	}
	latest, err := comms.LatestArtifact(list)
	if err != nil {
		return err
	}
	fmt.Println("latest comms artifact is:", latest)
	return nil
}
