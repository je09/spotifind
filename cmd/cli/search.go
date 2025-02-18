package main

import "github.com/spf13/cobra"

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for Spotify playlists by query",
	Run: func(cmd *cobra.Command, args []string) {
		queries, _ := cmd.Flags().GetStringSlice("queries")
		ignore, _ := cmd.Flags().GetStringSlice("ignore")
		market, _ := cmd.Flags().GetString("market")
		popular, _ := cmd.Flags().GetBool("popular")
		unpopular, _ := cmd.Flags().GetBool("unpopular")
		exportPath, _ := cmd.Flags().GetString("export")

		// queries and exportPath are required
		if len(queries) == 0 || exportPath == "" {
			cmd.Help()
			return
		}
		if len(ignore) > 0 {
			rootCmd.Printf(Green+"ignoring playlists descriptions containing: %v\n"+Reset, ignore)
		}

		spotifindHandler.Csv.Path = exportPath
		spotifindHandler.KnownPlaylists, _ = spotifindHandler.Csv.ReadFromFile()
		rootCmd.Printf(Green+"exporting playlists to: %s\n"+Reset, exportPath)
		rootCmd.Printf(Green+"search queries: %s\n"+Reset, queries)

		switch {
		case market != "":
			rootCmd.Printf(Blue+"searching playlists for %s market only\n"+Reset, market)
			spotifindHandler.SearchPlaylistForMarket(market, queries, ignore)
		case popular:
			rootCmd.Printf(Blue + "searching popular playlists only\n" + Reset)
			spotifindHandler.SearchPlaylistPopular(queries, ignore)
			break
		case unpopular:
			rootCmd.Printf(Blue + "searching unpopular playlists only\n" + Reset)
			spotifindHandler.SearchPlaylistUnpopular(queries, ignore)
			break
		default:
			rootCmd.Printf(Blue + "warning! searching all the markets!\n" + Reset)
			spotifindHandler.SearchPlaylistAllMarkets(queries, ignore)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringSliceP("queries", "q", nil, "Search queries (required)")
	searchCmd.Flags().StringP("export", "e", "", "Export playlists to CSV file (required)")
	searchCmd.Flags().StringSliceP("ignore", "i", nil, "Ignore playlists containing these words")
	searchCmd.Flags().StringP("market", "m", "", "Market to search in")
	searchCmd.Flags().BoolP("popular", "p", false, "Search in popular regions")
	searchCmd.Flags().BoolP("unpopular", "u", false, "Search in unpopular regions")
	searchCmd.Flags().SortFlags = false
}
