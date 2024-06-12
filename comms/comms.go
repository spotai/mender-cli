package comms

import (
	"fmt"
	"github.com/mendersoftware/mender-cli/client/deployments"
	"golang.org/x/mod/semver"
	"os"
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
			prov = ensureSemverPrefix(prov)
			if !semver.IsValid(prov) {
				_, _ = fmt.Fprintln(os.Stderr, "warning:", prov, "is not a valid semver version, skipping...")
				continue
			}
			arts = append(arts, &list.Artifacts[i])
		}
	}
	if len(arts) == 0 {
		return "", fmt.Errorf("no comms artifacts found")
	}
	sort.Slice(arts, func(i, j int) bool {
		v1 := ensureSemverPrefix(arts[i].ArtifactProvides[versionKey])
		v2 := ensureSemverPrefix(arts[j].ArtifactProvides[versionKey])
		return semver.Compare(v1, v2) > 0
	})
	return arts[0].Name, nil
}

// Ensures that a version string is prefixed with a 'v'
func ensureSemverPrefix(s string) string {
	if !strings.HasPrefix(s, "v") {
		return "v" + s
	}
	return s
}
