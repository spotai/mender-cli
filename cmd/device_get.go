package cmd

import (
	"errors"
	"fmt"
	"github.com/mendersoftware/mender-cli/client/devices"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deviceGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a device by hostname",
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewDeviceGetCmd(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

func init() {
}

type DeviceGetCmd struct {
	server     string
	skipVerify bool
	token      string
	hostname   string
}

func NewDeviceGetCmd(cmd *cobra.Command, args []string) (*DeviceGetCmd, error) {
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

	return &DeviceGetCmd{
		server:     server,
		token:      token,
		skipVerify: skipVerify,
		hostname:   args[0],
	}, nil
}

func (c *DeviceGetCmd) Run() error {

	client := devices.NewClient(c.server, c.skipVerify)
	return client.GetDeviceByHostname(c.token, c.hostname)
}
