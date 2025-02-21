package spotifind

type Spotifinder interface {
	Connector
	Search
}

type Connector interface {
	Reconnect(configAuth SpotifindAuth) error
	Continue(ch SpotifindChan, p ProgressChan) error
}

type Search interface {
	SearchPlaylistAllMarkets(ch SpotifindChan, p ProgressChan, q, ignore []string) error
	SearchPlaylistForMarket(ch SpotifindChan, p ProgressChan, market string, q, ignore []string) error
	SearchPlaylistPopular(ch SpotifindChan, p ProgressChan, q, ignore []string) error
	SearchPlaylistUnpopular(ch SpotifindChan, p ProgressChan, q, ignore []string) error
}
