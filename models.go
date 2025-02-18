package spotifind

// searchCriteria is an idicator, that a playlist has a curator contacts
var searchCriteria = []string{
	"@",
	"submithub",
	"sbmt",
	"submit",
	"submissions",
}

const SpotifyOwnerID = "spotify"

type SpotifindAuth struct {
	ClientID     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	RedirectURI  string `yaml:"redirectUri"`
}

type SpotifyOwner struct {
	ID          string
	DisplayName string
}

type Playlist struct {
	ID             string
	Name           string
	Collaborative  bool
	Description    string
	Styles         []string
	ExternalURLs   map[string]string
	Owner          SpotifyOwner
	TracksTotal    int
	FollowersTotal int
	Contacts       []string
	Region         string
}

type Progress struct {
	Total int
	Done  int
}

type PlaylistInfo struct {
	Playlist Playlist
	Progress Progress
}
