package container

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/t-hiroyoshi/tsumiki/version"
)

type DockerPayload struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current uint16 `json:"current"`
		Total   uint16 `json:"total"`
	} `json:"progressDetail"`
}

type containerClient struct {
	client *client.Client
}

var tagRegex = regexp.MustCompile(`(?P<name>[a-z0-9]+)(?P<version>:[._-][a-z0-9]+)*`)

func parseTag(tag string) (string, string, error) {
	p := tagRegex.FindAllStringSubmatch(tag, -1)
	fmt.Println(p)
	return "name", "version", nil
}

func NewContainerClient() (*containerClient, error) {
	cli, err := client.NewEnvClient()
	return &containerClient{client: cli}, err
}

func (c *containerClient) Host() string {
	return c.client.DaemonHost()
}

func (c *containerClient) ImagePull(ctx context.Context, tag string) (*version.PackageInfo, error) {
	loader, err := c.client.ImagePull(ctx, tag, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	defer loader.Close()

	payload := DockerPayload{}

	scanner := bufio.NewScanner(loader)
	for scanner.Scan() {
		json.Unmarshal(scanner.Bytes(), &payload)
		fmt.Printf("\t%+v\n", payload)
	}

	is, err := c.client.ImageList(ctx, types.ImageListOptions{})
	imageID := ""
	for _, s := range is {
		installed := false
		for _, t := range s.RepoTags {
			if t == tag {
				installed = true
			}
		}
		if installed {
			imageID = s.ID
			break
		}
	}

	return &version.PackageInfo{
		Tag:        tag,
		DockerHost: c.client.DaemonHost(),
		ImageID:    imageID,
	}, nil
}

func (c *containerClient) ImageRemove(ctx context.Context, info version.PackageInfo) error {
	_, err := c.client.ImageRemove(ctx, info.ImageID, types.ImageRemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}
