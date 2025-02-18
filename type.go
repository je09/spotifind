package spotifind

type SpotifyAPI interface {
	SearchPlaylistAllMarkets(ch SpotifindChan, p ProgressChan, q, ignore []string) error
	SearchPlaylistForMarket(ch SpotifindChan, p ProgressChan, market string, q, ignore []string) error
	SearchPlaylistPopular(ch SpotifindChan, p ProgressChan, q, ignore []string) error
	SearchPlaylistUnpopular(ch SpotifindChan, p ProgressChan, q, ignore []string) error
}
