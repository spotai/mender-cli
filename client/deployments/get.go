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
	"encoding/json"
	"fmt"
	"github.com/mendersoftware/mender-cli/log"
	"github.com/pkg/errors"
	"net/http"
)

type Deployment struct {
	Name         string
	Created      string
	Finished     string
	ArtifactName string `json:"artifact_name"`
	Artifacts    []string
	Status       string
}

func (c *Client) Get(
	token string,
	id string,
) (*Deployment, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.deploymentsURL, id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create post request: %s", err)
	}
	req.Header.Set("Authorization", "Bearer "+string(token))

	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %s", err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		if rsp.StatusCode == http.StatusUnauthorized {
			log.Verbf("GET request sent to '%s' failed with status %d", req.Host, rsp.StatusCode)
			return nil, errors.New("Unauthorized. Please Login first")
		}
		return nil, errors.New(
			fmt.Sprintf("GET request sent to '%s' failed with status %d", req.Host, rsp.StatusCode),
		)
	}
	buf := json.NewDecoder(rsp.Body)
	var dep Deployment
	err = buf.Decode(&dep)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %s", err)
	}
	return &dep, nil
}
