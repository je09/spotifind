package csv

// that's a shit code for gui and cli applications, don't reuse it ever.

import (
	"encoding/csv"
	"github.com/je09/spotifind"
	"io"
	"os"
	"strconv"
	"strings"
)

type CsvHandler struct {
	Path string
}

func NewCsvHandler(path string) *CsvHandler {
	return &CsvHandler{
		Path: path,
	}
}

func (c *CsvHandler) WriteToFile(playlist spotifind.Playlist) error {
	if c.Path == "" {
		return nil
	}

	file, err := os.OpenFile(c.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 5 styles max
	if len(playlist.Styles) > 5 {
		playlist.Styles = playlist.Styles[:5]
	}

	// Convert the playlist to a slice of strings
	data := []string{
		playlist.Name,
		strconv.Itoa(playlist.FollowersTotal),
		strings.Join(playlist.Styles, "; "),
		strings.Join(playlist.Contacts, "; "),
		playlist.Description,
		playlist.Region,
		playlist.ExternalURLs["spotify"],
	}

	// remove all the comas
	for i, d := range data {
		data[i] = strings.ReplaceAll(d, ",", "")
	}

	if err := writer.Write(data); err != nil {
		return err
	}

	return nil
}

func (c *CsvHandler) ReadFromFile() ([]string /*playlist names*/, error) {
	file, err := os.Open(c.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	playlists := make([]string, 0)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		playlist := spotifind.Playlist{
			Name: record[0],
			ExternalURLs: map[string]string{
				"spotify": record[6],
			},
		}

		playlists = append(playlists, playlist.ExternalURLs["spotify"])
	}

	return playlists, nil
}
