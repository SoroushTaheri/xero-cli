package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type CompetitionConfig struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type ChallengeConfig struct {
	Name    string `json:"name"`
	Path    string `json:"user_id"`
	EnterID string `json:"enter_id"`
}

type Config struct {
	AccessToken   string            `json:"access_token"`
	RefreshToken  string            `json:"refresh_token"`
	ParticipantID string            `json:"participant_id"`
	Joined        bool              `json:"joined"`
	Competition   CompetitionConfig `json:"competiton"`
	Challenges    []ChallengeConfig `json:"challenges"`
}

// go 1.18 is here. I WANT GENERICS ... AT ANY COST !!!
func first[T, U any](val T, _ U) T {
	return val
}

var (
	configPath = filepath.Join(first(os.UserHomeDir()), ".xero.conf") // HE HE! Shampoo! I just used a generic function.
)

func New() (*Config, error) {
	var xeroConfig Config

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		// File does not exist. Must create

		// Don't change the following line unless you know what you're doing
		xeroConfig.Competition.Path = "xero-ctf"
		xeroConfig.Competition.Name = "XeroCTF"

		file, err := json.MarshalIndent(xeroConfig, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to write .xero.config on '%s': %v", configPath, err)
		}

		if err := ioutil.WriteFile(configPath, file, 0600); err != nil {
			return nil, fmt.Errorf("failed to write .xero.config on '%s': %v", configPath, err)
		}

		return &xeroConfig, nil
	}

	// File exists
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not read '%s': %v", configPath, err)
	}

	if err := json.Unmarshal([]byte(file), &xeroConfig); err != nil {
		return nil, fmt.Errorf("could not parse JSON from '%s': %v", configPath, err)
	}

	// Don't change the following line unless you know what you're doing
	xeroConfig.Competition.Path = "xero-ctf"
	xeroConfig.Competition.Name = "XeroCTF"

	return &xeroConfig, nil
}

func (c *Config) Sync() error {
	file, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to write (sync) .xero.config on '%s': %v", configPath, err)
	}

	if err := ioutil.WriteFile(configPath, file, 0600); err != nil {
		return fmt.Errorf("failed to write (sync) .xero.config on '%s': %v", configPath, err)
	}

	return nil
}

func (c *Config) SetTokens(accessToken, refreshToken string) error {
	c.AccessToken, c.RefreshToken = accessToken, refreshToken
	return c.Sync()
}

func (c *Config) SetJoined(joined bool) error {
	c.Joined = joined
	return c.Sync()
}

func (c *Config) SetCompetitionConfig(competitionConfig CompetitionConfig) error {
	c.Competition = competitionConfig
	return c.Sync()
}
