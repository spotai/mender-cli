package cmd

import (
	"errors"
	"fmt"
	"github.com/mendersoftware/mender-cli/client/deployments"
	"github.com/mendersoftware/mender-cli/client/devices"
	"github.com/mendersoftware/mender-cli/comms"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

var ensureCmd = &cobra.Command{
	Use:   "ensure",
	Short: "Ensure a device has the latest comms version installed",
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewEnsureCmd(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

func init() {
	waitMinutes = ensureCmd.Flags().IntP(argWaitMinutes, "w", 32, "minutes to wait for deployment to finish (0 = don't wait)")
}

type EnsureCmd struct {
	server     string
	skipVerify bool
	token      string
	hostname   string
}

func NewEnsureCmd(cmd *cobra.Command, args []string) (*EnsureCmd, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("command requires one argument (the hostname)")
	}
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

	return &EnsureCmd{
		server:     server,
		token:      token,
		skipVerify: skipVerify,
		hostname:   args[0],
	}, nil
}

func (c *EnsureCmd) Run() error {

	devCli := devices.NewClient(c.server, c.skipVerify)
	dev, err := devCli.GetDeviceByHostname(c.token, c.hostname)
	if err != nil {
		return fmt.Errorf("failed to get device: %s", err)
	}
	depCli := deployments.NewClient(c.server, c.skipVerify)
	list, err := depCli.ListArtifacts(c.token)
	if err != nil {
		return fmt.Errorf("failed to list artifacts: %s", err)
	}
	artifactName, version, err := comms.LatestArtifactNameAndVersion(list)
	if err != nil {
		return fmt.Errorf("failed to find latest artifact: %s", err)
	}
	fmt.Println("latest artifact name and version:", artifactName, "/", version)
	found := false
	for _, attr := range dev.Attributes {
		if attr.Name == comms.VersionKey {
			found = true
			if attr.Value == version {
				fmt.Println("device already has the latest version -- ok")
				return nil
			}
			fmt.Println("device has out-of-date version:", attr.Value)
			break
		}
	}
	if !found {
		fmt.Println("device does not yet have the artifact")
	}
	depName := fmt.Sprintf("comms-%s", c.hostname)
	deploymentId, err := depCli.Create(c.token, depName, artifactName, dev.Id)
	if err != nil {
		return fmt.Errorf("failed to create deployment: %s", err)
	}
	fmt.Println("deployment id:", deploymentId)
	if *waitMinutes > 0 {
		for i := 0; i < *waitMinutes; i++ {
			dep, err := depCli.Get(c.token, deploymentId)
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
