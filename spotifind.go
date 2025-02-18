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

const spotifyPlaylistTrackLimit = 1000

type Spotifind struct {
	auth        *spotifyauth.Authenticator
	ctx         context.Context
	client      *spotify.Client
	limiter     *rate.Limiter
	errHandling ErrorHandling

	visitedPlaylists map[string]bool
	visitedMutex     sync.Mutex
	progressMutex    sync.Mutex

	totalPlaylists int
	donePlaylists  int
}

func NewSpotifind(configAuth SpotifindAuth, retry bool) (*Spotifind, error) {
	ctx := context.Background()

	creds := &clientcredentials.Config{
		ClientID:     configAuth.ClientID,
		ClientSecret: configAuth.ClientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := creds.Token(ctx)
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

func (s *Spotifind) SearchPlaylistAllMarkets(ch SpotifindChan, p ProgressChan, q, ignore []string) error {
	s.totalPlaylists = len(MarketsAll) * spotifyPlaylistTrackLimit * len(q)
	for _, query := range q {
		for _, market := range MarketsAll {
			if err := s.searchPlaylistForMarket(ch, p, query, market, ignore); err != nil {
				return err
			}
		}
	}
	close(ch)

	return nil
}

func (s *Spotifind) SearchPlaylistForMarket(ch SpotifindChan, p ProgressChan, market string, q, ignore []string) error {
	s.totalPlaylists = spotifyPlaylistTrackLimit * len(q)
	for _, query := range q {
		err := s.searchPlaylistForMarket(ch, p, query, market, ignore)
		if err != nil {
			return err
		}
	}
	close(ch)

	return nil
}

func (s *Spotifind) SearchPlaylistPopular(ch SpotifindChan, p ProgressChan, q, ignore []string) error {
	s.totalPlaylists = len(marketsPopular) * spotifyPlaylistTrackLimit * len(q)
	for _, query := range q {
		for _, market := range marketsPopular {
			if err := s.searchPlaylistForMarket(ch, p, query, market, ignore); err != nil {
				return err
			}
		}
	}
	close(ch)

	return nil
}

func (s *Spotifind) SearchPlaylistUnpopular(ch SpotifindChan, p ProgressChan, q, ignore []string) error {
	s.totalPlaylists = len(marketsUnpopular) * spotifyPlaylistTrackLimit * len(q)

	for _, query := range q {
		for _, market := range marketsUnpopular {
			if err := s.searchPlaylistForMarket(ch, p, query, market, ignore); err != nil {
				return err
			}
		}
	}
	close(ch)

	return nil
}

func (s *Spotifind) searchPlaylistForMarket(ch SpotifindChan, p ProgressChan, q, market string, ignore []string) error {
	res, err := s.searchPlaylists(q, market)
	if err != nil {
		return s.errHandling.Handle(err)
	}
	s.processPlaylists(ch, p, res, ignore, market)

	return nil
}

func (s *Spotifind) searchPlaylists(query, market string) (*spotify.SearchResult, error) {
	opts := []spotify.RequestOption{
		spotify.Market(market),
	}
	return s.client.Search(s.ctx, query, spotify.SearchTypePlaylist, opts...)
}

func (s *Spotifind) processPlaylists(ch SpotifindChan, p ProgressChan, results *spotify.SearchResult, ignore []string, region string) {
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

			s.sendPlaylistToChannel(ch, playlist, region)
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
}

func (s *Spotifind) checkPlaylistConditions(item spotify.SimplePlaylist, ignore []string) bool {
	if item.Owner.ID == SpotifyOwnerID {
		return false
	}

	for _, criteria := range ignore {
		if strings.Contains(strings.ToLower(item.Name), criteria) {
			return false
		}
		if strings.Contains(strings.ToLower(item.Description), criteria) {
			return false
		}
	}

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

func (s *Spotifind) getContacts(item spotify.SimplePlaylist) []string {
	// get all the lines started from @ from string
	var res []string
	for _, line := range strings.Split(item.Description, " ") {
		// todo add submithub support
		if strings.Contains(line, "@") {
			// TODO fix later
			contact, _ := url.QueryUnescape(line)
			res = append(res, contact)
		}
	}

	return res
}

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

func (s *Spotifind) getPlaylistStyles(playlist *spotify.FullPlaylist) ([]string, error) {
	artistIDs := make([]spotify.ID, 0)
	for _, track := range playlist.Tracks.Tracks {
		for _, artist := range track.Track.Artists {
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

func (s *Spotifind) sendPlaylistToChannel(ch SpotifindChan, item spotify.SimplePlaylist, region string) {
	contacts := s.getContacts(item)

	playlistInfo, _ := s.getPlaylistInfo(item.ID)
	styles, _ := s.getPlaylistStyles(playlistInfo)
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
			FollowersTotal: int(playlistInfo.Followers.Count),
			Contacts:       contacts,
			Styles:         styles,
			Region:         region,
		},
		Progress: s.getProgress(),
	}
}

func (s *Spotifind) sendProgressToChannel(ch ProgressChan, found, limit int) {
	// if we found less than the limit, we need to subtract the difference
	if found < limit {
		s.setTotalPlaylists(s.totalPlaylists - limit + found)
	}
	ch <- s.getProgress()
}

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

func (s *Spotifind) getProgress() Progress {
	done, total := 0, 0

	s.progressMutex.Lock()
	done = s.donePlaylists
	total = s.totalPlaylists
	s.progressMutex.Unlock()

	return Progress{
		Done:  done,
		Total: total,
	}
}
