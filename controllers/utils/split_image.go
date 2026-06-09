package utils

import (
	"strings"

	docker "github.com/distribution/reference"
	"github.com/pkg/errors"
)

// SplitImage split image string to image and tag
func SplitImage(fullImage string) (string, string) {
	// ReferenceRegexp.FindStringSubmatch return an array with following content
	// [<original string>, <[domain/]image>, <tag>, <digest>]
	result := docker.ReferenceRegexp.FindStringSubmatch(strings.ToLower(fullImage))

	// Default empty values
	image := ""
	tag := ""

	switch len(result) {
	case 1, 2:
		// Full image has only domain and image, or only image, without tag. For example:
		// hub.docker.io:17000/prom/prometheus
		// hub.docker.io/prom/prometheus
		// prom/prometheus
		// So as result can return original string
		image = result[0]
	case 3, 4:
		// Full image contains domain, image and tag, image and tag. For example:
		// hub.docker.io:17000/prom/prometheus:v2.1.0
		// hub.docker.io/prom/prometheus:v2.1.0
		// prom/prometheus:v2.1.0
		image = result[1]
		tag = result[2]
	}

	return image, tag
}

// SplitPathImage split image string to path and tag
func SplitPathImage(fullImage string) (path string, tag string, err error) {
	if fullImage == "" {
		return "", "", nil
	}
	var parsed docker.Named
	parsed, err = docker.ParseNamed(fullImage)
	if err != nil {
		return "", "", err
	}

	path = docker.Path(parsed)

	tagged, ok := parsed.(docker.Tagged)
	if !ok {
		return "", "", errors.New("can't parse image tag")
	}
	tag = tagged.Tag()

	return
}
