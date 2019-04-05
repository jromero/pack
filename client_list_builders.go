package pack

import (
	"context"
	"github.com/buildpack/pack/builder"
	"github.com/buildpack/pack/style"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type BuilderListing struct {
	ImageSummary types.ImageSummary
	BuilderInfo *BuilderInfo
}

func (c *Client) ListBuilders(ctx context.Context) ([]*BuilderListing, error) {
	var (
		result []*BuilderListing
		err error
	)

	entries, err := c.docker.ImageList(ctx, types.ImageListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", "io.buildpacks.stack.id"),
			filters.Arg("label", builder.MetadataLabel),
		),
	})

	if err != nil {
		return nil, err
	}

	for _, imageSummary := range entries {
		if len(imageSummary.RepoTags) >= 1 {
			info, err := c.InspectBuilder(imageSummary.ID, true)
			if err != nil {
				c.logger.Error("failed to get builder info of %s: %s", style.Symbol(imageSummary.ID), err.Error())
			} else if info != nil {
				result = append(result, &BuilderListing{
					ImageSummary: imageSummary,
					BuilderInfo: info,
				})
			}
		}
	}

	return result, nil
}