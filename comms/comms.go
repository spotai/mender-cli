package comms

import (
	"fmt"
	"github.com/mendersoftware/mender-cli/client/deployments"
	"golang.org/x/mod/semver"
	"sort"
	"strings"
)

const (
	versionKey = "data_partition.comms.version"
)

// LatestArtifact returns the name of the artifact with the most recent comms version
func LatestArtifact(list *deployments.ArtifactsList) (string, error) {
	// Make a list of pointers to avoid copying maps during sorting
	var arts []*deployments.ArtifactData
	for i, art := range list.Artifacts {
		prov, ok := art.ArtifactProvides[versionKey]
		if ok {
			if !strings.HasPrefix(prov, "v") {
				prov = "v" + prov
			}
			if !semver.IsValid(prov) {
				fmt.Println("warning:", prov, "is not a valid semver version, skipping...")
				continue
			}
			arts = append(arts, &list.Artifacts[i])
			fmt.Println("candidate:", list.Artifacts[i].Name)
		}
	}
	if len(arts) == 0 {
		return "", fmt.Errorf("no comms artifacts found")
	}
	sort.Slice(arts, func(i, j int) bool {
		v1 := arts[i].ArtifactProvides[versionKey]
		if !strings.HasPrefix(v1, "v") {
			v1 = "v" + v1
		}
		v2 := arts[j].ArtifactProvides[versionKey]
		if !strings.HasPrefix(v2, "v") {
			v2 = "v" + v2
		}
		return semver.Compare(v1, v2) > 0
	})
	return arts[0].Name, nil
}
