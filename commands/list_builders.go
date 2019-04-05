package commands

import (
	"fmt"
	"github.com/buildpack/pack/logging"
	"github.com/buildpack/pack/style"
	"github.com/spf13/cobra"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

func ListBuilders(logger *logging.Logger, client PackClient) *cobra.Command {
	ctx := createCancellableContext()

	cmd := &cobra.Command{
		Use:   "list-builders",
		Short: "List builders",
		Args:  cobra.ExactArgs(0),
		RunE: logError(logger, func(cmd *cobra.Command, args []string) error {
			listings, err := client.ListBuilders(ctx)
			if err != nil {
				return err
			}

			logger.Info("Local builders:\n")
			sort.Slice(listings, func(i, j int) bool {
				// TODO handle empty tags
				return listings[i].ImageSummary.RepoTags[0] < listings[j].ImageSummary.RepoTags[0]
			})

			tw := tabwriter.NewWriter(logger.RawWriter(), 10, 10, 5, ' ', tabwriter.TabIndent)

			for _, v := range listings {
				t := time.Unix(v.ImageSummary.Created, 0)
				_, _ = tw.Write([]byte(fmt.Sprintf(
					"\t%s\t%s\t%s\t\n",
					strings.TrimPrefix(v.ImageSummary.ID, "sha256:")[:12],
					style.Symbol(v.ImageSummary.RepoTags[0]),
					style.Surpressed(t.Format("Jan _2 15:04:05")),
				)))
			}

			_ = tw.Flush()

			logger.Info("")
			suggestBuilders(logger)

			return nil
		}),
	}

	AddHelpFlag(cmd, "list-builders")
	return cmd
}
