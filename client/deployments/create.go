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
package deployments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mendersoftware/mender-cli/log"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type newDeployment struct {
	Name         string   `json:"name"`
	ArtifactName string   `json:"artifact_name"`
	Devices      []string `json:"devices"`
}

func (c *Client) Create(
	token string,
	name string,
	artifactName string,
	deviceId string,
) (string, error) {
	nd := newDeployment{
		Name:         name,
		ArtifactName: artifactName,
		Devices:      []string{deviceId},
	}
	data, err := json.Marshal(&nd)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body: %s", err)
	}
	reader := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPost, c.deploymentsURL, reader)
	if err != nil {
		return "", fmt.Errorf("failed to create post request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+string(token))

	rsp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %s", err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusCreated {
		if rsp.StatusCode == http.StatusUnauthorized {
			log.Verbf("deployment creation sent to '%s' failed with status %d", req.Host, rsp.StatusCode)
			return "", errors.New("Unauthorized. Please Login first")
		} else if rsp.StatusCode == http.StatusConflict {
			log.Verbf("deployment creation sent to '%s' failed with status %d", req.Host, rsp.StatusCode)
			return "", errors.New("unfinished deployment with same name and targeting same device already exists")
		}
		return "", errors.New(
			fmt.Sprintf("deployment creation sent to '%s' failed with status %d", req.Host, rsp.StatusCode),
		)
	}
	loc := rsp.Header.Get("location")
	if loc == "" {
		return "", fmt.Errorf("response does not contain a Location header")
	}
	elems := strings.Split(loc, "/")
	deploymentId := elems[len(elems)-1]
	return deploymentId, nil
}
