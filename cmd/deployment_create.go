package cmd

import (
	"errors"
	"fmt"
	"github.com/mendersoftware/mender-cli/client/deployments"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
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
	argWaitMinutes    = "wait"
)

var waitMinutes *int

func init() {
	deploymentCreateCmd.Flags().StringP(argDeploymentName, "n", "", "deployment name")
	deploymentCreateCmd.Flags().StringP(argArtifactName, "a", "", "artifact name")
	deploymentCreateCmd.Flags().StringP(argDeviceId, "d", "", "device id")
	waitMinutes = deploymentCreateCmd.Flags().IntP(argWaitMinutes, "w", 32, "minutes to wait for deployment to finish (0 = don't wait)")
	CheckErr(deploymentCreateCmd.MarkFlagRequired(argDeploymentName))
	CheckErr(deploymentCreateCmd.MarkFlagRequired(argArtifactName))
	CheckErr(deploymentCreateCmd.MarkFlagRequired(argDeviceId))
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
	if *waitMinutes > 0 {
		for i := 0; i < *waitMinutes; i++ {
			dep, err := client.Get(c.token, deploymentId)
			if err != nil {
				return fmt.Errorf("failed to get deployment: %s", err)
			}
			if dep.Status == "finished" {
				fmt.Println("deployment finished")
				return nil
			}
			fmt.Println("deployment status is still", dep.Status, "... checking again in one minute")
			time.Sleep(time.Minute)
		}
		return fmt.Errorf("deployment still unfinished after %d minutes", *waitMinutes)
	} else {
		return nil
	}
}
