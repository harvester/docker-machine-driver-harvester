package harvester

import (
	"encoding/base64"
	"fmt"
	"hash/fnv"

	"k8s.io/apimachinery/pkg/util/validation"
)

// formatLabelValue returns v if it meets the standards for a Kubernetes
// label value. Otherwise, it returns a hash which meets the requirements.
// ref: https://github.com/kubernetes-sigs/cluster-api/blob/010af7f92f98ead971742d347d258a8786d5d57c/util/labels/format/helpers.go
func formatLabelValue(v string) (string, error) {
	// a valid Kubernetes label value must:
	// - be less than 64 characters long.
	// - be an empty string OR consist of alphanumeric characters, '-', '_' or '.'.
	// - start and end with an alphanumeric character
	if len(validation.IsValidLabelValue(v)) == 0 {
		return v, nil
	}

	hasher := fnv.New32a()
	if _, err := hasher.Write([]byte(v)); err != nil {
		return "", err
	}

	// use base64 URL encoding to avoid special characters like '/' in the generated
	// hash
	return fmt.Sprintf("hash_%s_z",
		base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))), nil
}
