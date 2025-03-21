package spotifind

import (
	"context"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/time/rate"
	"net/url"
	"strings"
	"sync"
	"time"
)

type SpotifindChan chan PlaylistInfo
type ProgressChan chan Progress

// spotifyPlaylistTrackLimit - max number of tracks official SpotifyAPI can provide.
const spotifyPlaylistTrackLimit = 1000

type Spotifind struct {
	auth        *spotifyauth.Authenticator
	client      *spotify.Client
	limiter     *rate.Limiter
	errHandling ErrorHandling

	visitedPlaylists map[string]bool
	visitedMutex     sync.Mutex
	progressMutex    sync.Mutex

	totalPlaylists int
	donePlaylists  int

	marketsLeft []string
	queriesLeft []string
	ignores     []string

	ctx context.Context
}

// NewSpotifind creates a new Spotifind instance.
// configAuth - Spotify API credentials.
// retry - if true, the client will retry the request if it fails.
// Do not use retry, if you wish to implement your own retry logic.
func NewSpotifind(configAuth SpotifindAuth, retry bool) (*Spotifind, error) {
	ctx := context.Background()

	c := &clientcredentials.Config{
		ClientID:     configAuth.ClientID,
		ClientSecret: configAuth.ClientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := c.Token(ctx)
	if err != nil {
		return nil, err
	}

	auth := spotifyauth.New(spotifyauth.WithClientID(configAuth.ClientID), spotifyauth.WithClientSecret(configAuth.ClientSecret))
	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient, spotify.WithRetry(retry))

	return &Spotifind{
		auth:             auth,
		ctx:              ctx,
		client:           client,
		limiter:          rate.NewLimiter(rate.Every(time.Second), 1),
		visitedPlaylists: make(map[string]bool),
	}, nil
}

// Reconnect reconnects to the Spotify API with new app credentials. In case there's a network error.
// For educational purposes only!
func (s *Spotifind) Reconnect(configAuth SpotifindAuth) error {
	c := &clientcredentials.Config{
		ClientID:     configAuth.ClientID,
		ClientSecret: configAuth.ClientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := c.Token(s.ctx)
	if err != nil {
		return err
	}

	s.auth = spotifyauth.New(spotifyauth.WithClientID(configAuth.ClientID), spotifyauth.WithClientSecret(configAuth.ClientSecret))
	httpClient := spotifyauth.New().Client(s.ctx, token)
	s.client = spotify.New(httpClient, spotify.WithRetry(true))

	return nil
}

// Continue previous search with the same queries and ignores if there were any errors.
func (s *Spotifind) Continue(ch SpotifindChan, p ProgressChan) error {
	for _, query := range s.queriesLeft {
		for _, market := range s.marketsLeft {
			if err := s.searchPlaylistForMarket(ch, p, query, market, s.ignores); err != nil {
				return err
			}
		}
	}
	close(ch)

	return nil
}

// SearchPlaylistAllMarkets searches for playlists in all markets Spotify supports.
// Used it carefully, as it can take a lot of time.
func (s *Spotifind) SearchPlaylistAllMarkets(ch SpotifindChan, p ProgressChan, q /* queries to search */, ignore []string) error {
	// For progress bar.
	s.totalPlaylists = len(MarketsAll) * spotifyPlaylistTrackLimit * len(q)

	// Caching in case of reconnecting.
	s.marketsLeft = MarketsAll
	s.queriesLeft = q

	for _, query := range q {
		for _, market := range MarketsAll {
			if err := s.searchPlaylistForMarket(ch, p, query, market, ignore); err != nil {
				return err
			}
			s.markMarketAsDone()
		}
		s.markQueryAsDone()
	}
	close(ch)

	return nil
}

// SearchPlaylistForMarket searches for playlists in a specific market (note, that Spotify has a limited number of markets).
// market is a named for the country standardized as ISO 3166-1 alpha-2.
// ignore is a list of strings to ignore in the playlist name and description.
func (s *Spotifind) SearchPlaylistForMarket(ch SpotifindChan, p ProgressChan, market string, q, ignore []string) error {
	// For progress bar.
	s.totalPlaylists = spotifyPlaylistTrackLimit * len(q)

	// Caching in case of reconnecting.
	s.marketsLeft = []string{market}
	s.queriesLeft = q

	for _, query := range q {
		err := s.searchPlaylistForMarket(ch, p, query, market, ignore)
		if err != nil {
			return err
		}
		s.markQueryAsDone()
	}
	close(ch)

	return nil
}

// SearchPlaylistPopular searches for playlists in popular markets only.
// Popular markets include US, GB, DE, FR, etc. Popularity of markets is strictly subjective.
// Probably, the best option to use, when searching for playlists.
// ignore is a list of strings to ignore in the playlist name and description.
func (s *Spotifind) SearchPlaylistPopular(ch SpotifindChan, p ProgressChan, q, ignore []string) error {
	// For progress bar.
	s.totalPlaylists = len(marketsPopular) * spotifyPlaylistTrackLimit * len(q)

	// Caching in case of reconnecting.
	s.marketsLeft = marketsPopular
	s.queriesLeft = q

	for _, query := range q {
		for _, market := range marketsPopular {
			if err := s.searchPlaylistForMarket(ch, p, query, market, ignore); err != nil {
				return err
			}
			s.markMarketAsDone()
		}
		s.markQueryAsDone()
	}
	close(ch)

	return nil
}

// SearchPlaylistUnpopular searches for playlists in unpopular markets only.
// Unpopular markets are everything except popular markets.
// ignore is a list of strings to ignore in the playlist name and description.
func (s *Spotifind) SearchPlaylistUnpopular(ch SpotifindChan, p ProgressChan, q, ignore []string) error {
	// For progress bar.
	s.totalPlaylists = len(marketsUnpopular) * spotifyPlaylistTrackLimit * len(q)

	// Caching in case of reconnecting.
	s.marketsLeft = marketsUnpopular
	s.queriesLeft = q

	for _, query := range q {
		for _, market := range marketsUnpopular {
			if err := s.searchPlaylistForMarket(ch, p, query, market, ignore); err != nil {
				return err
			}
			s.markMarketAsDone()
		}
		s.markQueryAsDone()
	}
	close(ch)

	return nil
}

func (s *Spotifind) searchPlaylistForMarket(ch SpotifindChan, p ProgressChan, q, market string, ignore []string) error {
	s.ignores = ignore

	res, err := s.searchPlaylists(q, market)
	if err != nil {
		return s.errHandling.Handle(err)
	}
	if err = s.processPlaylists(ch, p, res, ignore, market); err != nil {
		return s.errHandling.Handle(err)
	}

	return nil
}

func (s *Spotifind) searchPlaylists(query, market string) (*spotify.SearchResult, error) {
	opts := []spotify.RequestOption{
		spotify.Market(market),
	}
	return s.client.Search(s.ctx, query, spotify.SearchTypePlaylist, opts...)
}

// processPlaylists processes the search results and sends them to the channel if there's any contact info.
func (s *Spotifind) processPlaylists(ch SpotifindChan, p ProgressChan, results *spotify.SearchResult, ignore []string, region string) error {
	// if there are no playlists, we need to skip this market.
	for results.Playlists != nil {
		for _, playlist := range results.Playlists.Playlists {
			s.incrementProgress()
			s.sendProgressToChannel(p, len(results.Playlists.Playlists), int(results.Playlists.Limit))
			if !s.checkPlaylistConditions(playlist, ignore) {
				continue
			}
			if s.hasVisitedPlaylist(playlist.ID.String()) {
				continue
			}
			playlistInfo, err := s.getPlaylistInfo(playlist.ID)
			if err != nil {
				return err
			}
			styles, err := s.getPlaylistStyles(playlistInfo)
			if err != nil {
				return err
			}
			if !s.checkStyleConditions(styles, ignore) {
				continue
			}

			if err := s.sendPlaylistToChannel(ch, playlistInfo, styles, region); err != nil {
				return err
			}
			s.rememberVisitedPlaylist(playlist.ID.String())
		}
		// for some reason next doesn't stop it
		if results.Playlists.Offset+results.Playlists.Limit >= results.Playlists.Total {
			break
		}
		if err := s.client.NextPlaylistResults(s.ctx, results); err != nil {
			break
		}
	}
	if results.Playlists.Total == 0 || results.Playlists == nil {
		// skip this market
		s.sendProgressToChannel(p, 0, 1000)
	}

	return nil
}

// checkPlaylistConditions checks if the playlist has any contact info.
// Criteria for the contact info is described in models.go as searchCriteria string list.
func (s *Spotifind) checkPlaylistConditions(item spotify.SimplePlaylist, ignore []string) bool {
	// Playlists owned by Spotify never possess any contact info, and are on the top of the search results.
	if item.Owner.ID == SpotifyOwnerID {
		return false
	}

	// if the playlist has ignore criteria, we need to skip it.
	for _, criteria := range ignore {
		// Ignore empty strings.
		if criteria == "" {
			continue
		}

		if strings.Contains(strings.ToLower(item.Name), criteria) {
			return false
		}
		if strings.Contains(strings.ToLower(item.Description), criteria) {
			return false
		}
	}

	// if the playlist has search criteria, we need to check if it has any contact info.
	for _, criteria := range searchCriteria {
		if strings.Contains(strings.ToLower(item.Description), criteria) {
			return true
		}
		if strings.Contains(strings.ToLower(item.Name), criteria) {
			return true
		}
	}

	return false
}

// checkStyleConditions checks if the playlist's styles match the ignore criteria.
// If the playlist has any of the ignore criteria, we need to skip it.
func (s *Spotifind) checkStyleConditions(styles, ignore []string) bool {
	for _, style := range styles {
		for _, criteria := range ignore {
			if criteria == "" {
				continue
			}

			if strings.Contains(strings.ToLower(style), criteria) {
				return false
			}
		}
	}

	return true
}

// getContacts gets all the contacts from the playlist description.
// basically, it tries to extract everything that looks like an email, or instragram, twitter, etc handle.
func (s *Spotifind) getContacts(item spotify.SimplePlaylist) []string {
	var res []string
	for _, line := range strings.Split(item.Description, " ") {
		if strings.Contains(line, "@") {
			contact, _ := url.QueryUnescape(line)
			res = append(res, contact)
		}
	}

	return res
}

// getArtistsStyles gets all the styles from the artists in the playlist and returns them sorted by genre incidence.
func (s *Spotifind) getArtistsStyles(artistIDs []spotify.ID) ([]string, error) {
	styleMap := make(map[string]int)

	// chunk the artistIDs
	for i := 0; i < len(artistIDs); i += 50 {
		end := i + 50
		if end > len(artistIDs) {
			end = len(artistIDs)
		}
		artists, err := s.client.GetArtists(s.ctx, artistIDs[i:end]...)
		if err != nil {
			return nil, err
		}
		// for some reason spotify can return nil artists
		if len(artists) == 0 {
			continue
		}
		for _, artist := range artists {
			if artist == nil || len(artist.Genres) == 0 {
				continue
			}

			for _, genre := range artist.Genres {
				styleMap[genre]++
			}
		}
	}

	return SortStyleMap(styleMap), nil
}

// getPlaylistStyles gets all the styles from the playlist and returns them sorted by genre incidence.
func (s *Spotifind) getPlaylistStyles(playlist *spotify.FullPlaylist) ([]string, error) {
	artistIDs := make([]spotify.ID, 0)
	for _, track := range playlist.Tracks.Tracks {
		for _, artist := range track.Track.Artists {
			if artist.ID == "" {
				continue
			}

			artistIDs = append(artistIDs, artist.ID)
		}
	}

	return s.getArtistsStyles(artistIDs)
}

func (s *Spotifind) getPlaylistInfo(id spotify.ID) (*spotify.FullPlaylist, error) {
	opts := []spotify.RequestOption{
		spotify.Limit(100),
	}
	return s.client.GetPlaylist(s.ctx, id, opts...)
}

func (s *Spotifind) sendPlaylistToChannel(ch SpotifindChan, item *spotify.FullPlaylist, styles []string, region string) error {
	contacts := s.getContacts(item.SimplePlaylist)

	ch <- PlaylistInfo{
		Playlist: Playlist{
			ID:            item.ID.String(),
			Name:          item.Name,
			Collaborative: item.Collaborative,
			Description:   item.Description,
			ExternalURLs:  item.ExternalURLs,
			Owner: SpotifyOwner{
				ID:          item.Owner.ID,
				DisplayName: item.Owner.DisplayName,
			},
			TracksTotal:    int(item.Tracks.Total),
			FollowersTotal: int(item.Followers.Count),
			Contacts:       contacts,
			Styles:         styles,
			Region:         region,
		},
		Progress: s.getProgress(),
	}

	return nil
}

func (s *Spotifind) sendProgressToChannel(ch ProgressChan, found, limit int) {
	// if we found less than the limit, we need to subtract the difference
	if found < limit {
		s.setTotalPlaylists(s.totalPlaylists - limit + found)
	}
	ch <- s.getProgress()
}

// probably need to use RWMutex or something like sync.Pool, but I'm just too lazy to rewrite this old code.
// if it ain't broke, don't fix it. right? =)

func (s *Spotifind) rememberVisitedPlaylist(id string) {
	s.visitedMutex.Lock()
	s.visitedPlaylists[id] = true
	s.visitedMutex.Unlock()
}

func (s *Spotifind) hasVisitedPlaylist(id string) bool {
	s.visitedMutex.Lock()
	_, ok := s.visitedPlaylists[id]
	s.visitedMutex.Unlock()

	return ok
}

func (s *Spotifind) incrementProgress() {
	s.progressMutex.Lock()
	s.donePlaylists++
	s.progressMutex.Unlock()
}

func (s *Spotifind) setTotalPlaylists(total int) {
	s.progressMutex.Lock()
	s.totalPlaylists = total
	s.progressMutex.Unlock()
}

func (s *Spotifind) markMarketAsDone() {
	s.progressMutex.Lock()
	if len(s.marketsLeft) == 0 {
		s.progressMutex.Unlock()
		return
	}
	s.marketsLeft = s.marketsLeft[1:]
	s.progressMutex.Unlock()
}

func (s *Spotifind) markQueryAsDone() {
	s.progressMutex.Lock()
	if len(s.queriesLeft) == 0 {
		s.progressMutex.Unlock()
		return
	}
	s.queriesLeft = s.queriesLeft[1:]
	s.progressMutex.Unlock()
}

func (s *Spotifind) getProgress() Progress {
	done, total := 0, 0

	s.progressMutex.Lock()
	done = s.donePlaylists
	total = s.totalPlaylists
	s.progressMutex.Unlock()

	// It's a hack mostly, but it works.
	if done > total {
		total = done
	}

	return Progress{
		Done:  done,
		Total: total,
	}
}
