package cmd

import (
	"errors"
	"fmt"
	"github.com/mendersoftware/mender-cli/client/deployments"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deploymentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a deployment",
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewDeploymentCreateCmd(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

const (
	argDeploymentName = "name"
	argArtifactName   = "artifact"
	argDeviceId       = "device"
)

func init() {
	deploymentCreateCmd.Flags().StringP(argDeploymentName, "n", "", "deployment name")
	deploymentCreateCmd.Flags().StringP(argArtifactName, "a", "", "artifact name")
	deploymentCreateCmd.Flags().StringP(argDeviceId, "d", "", "device id")
}

type DeploymentCreateCmd struct {
	server         string
	skipVerify     bool
	token          string
	deploymentName string
	artifactName   string
	deviceId       string
}

func NewDeploymentCreateCmd(cmd *cobra.Command, args []string) (*DeploymentCreateCmd, error) {

	server := viper.GetString(argRootServer)
	if server == "" {
		return nil, errors.New("no server specified")
	}

	skipVerify, err := cmd.Flags().GetBool(argRootSkipVerify)
	if err != nil {
		return nil, err
	}

	token, err := getAuthToken(cmd)
	if err != nil {
		return nil, err
	}

	deploymentName, err := cmd.Flags().GetString(argDeploymentName)
	if err != nil {
		return nil, err
	}
	artifactName, err := cmd.Flags().GetString(argArtifactName)
	if err != nil {
		return nil, err
	}
	deviceId, err := cmd.Flags().GetString(argDeviceId)
	if err != nil {
		return nil, err
	}
	return &DeploymentCreateCmd{
		server:         server,
		token:          token,
		skipVerify:     skipVerify,
		deploymentName: deploymentName,
		artifactName:   artifactName,
		deviceId:       deviceId,
	}, nil
}

func (c *DeploymentCreateCmd) Run() error {

	client := deployments.NewClient(c.server, c.skipVerify)
	deploymentId, err := client.Create(c.token, c.deploymentName, c.artifactName, c.deviceId)
	if err != nil {
		return err
	}
	fmt.Println("deployment id:", deploymentId)
	return nil
}
