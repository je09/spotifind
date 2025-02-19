package cli

import (
	"fmt"
	"github.com/je09/spotifind"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/rand"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "spotyfind2",
	Short: "Spotifind2 is a CLI tool for searching Spotify playlists",
}
var spotifindHandler *SpotifyHandler
var configs []spotifind.SpotifindAuth

type credsType struct {
	Credits []spotifind.SpotifindAuth `yaml:"credits"`
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	c := ".spotyfind2.yml"
	viper.SetConfigType("yaml")
	viper.SetConfigName(c)

	//home, _ := os.UserHomeDir()
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("spotyfind2")
	viper.BindEnv("spotify_client_id")
	viper.BindEnv("spotify_client_secret")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading configs file:", err)
		os.Exit(1)
	}
	var creds credsType
	err := viper.Unmarshal(&creds)
	configs = creds.Credits

	// randomize configs order
	for i := range configs {
		j := rand.Intn(i + 1)
		configs[i], configs[j] = configs[j], configs[i]

	}

	if err != nil {
		return
	}
	spotifindHandler, err = NewSpotifyHandler()
	if err != nil {
		panic(err)
	}
}
