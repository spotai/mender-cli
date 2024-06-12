package devices

import (
	"encoding/json"
	"fmt"
	"github.com/mendersoftware/mender-cli/client"
)

type attribute struct {
	Name  string
	Scope string
	Value any
}

type inventoryData struct {
	Id         string `json:"id"`
	Attributes []attribute
}

func (c *Client) GetDeviceByHostname(token string, hostname string) (*inventoryData, error) {
	url := fmt.Sprintf("%s/devices?hostname=%s", c.deviceInventoryURL, hostname)
	body, err := client.DoGetRequest(token, url, c.client)
	if err != nil {
		return nil, err
	}

	var devices []inventoryData
	err = json.Unmarshal(body, &devices)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %s", err)
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("no device found")
	}
	return &devices[0], nil
}

func (c *Client) PrintDeviceByHostname(token string, hostname string) error {
	device, err := c.GetDeviceByHostname(token, hostname)
	if err != nil {
		return err
	}
	fmt.Println("id:", device.Id)
	for _, attr := range device.Attributes {
		fmt.Printf("%s (%s): %v\n", attr.Name, attr.Scope, attr.Value)
	}
	return nil
}
