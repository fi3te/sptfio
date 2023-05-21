package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"gopkg.in/yaml.v3"
)

const callbackPattern = "[a-z\\-]*"
const portPattern = "[1-9][0-9]*"

var redirectURLRegexp *regexp.Regexp

func init() {
	redirectURLPattern := fmt.Sprintf("^http://localhost:(%s)(/%s)$", portPattern, callbackPattern)
	redirectURLRegexp = regexp.MustCompile(redirectURLPattern)
}

type Config struct {
	ClientID      string `yaml:"client-id"`
	InputFilePath string `yaml:"input-file-path"`
	RedirectURL   string `yaml:"redirect-url"`
	PlaylistName  string `yaml:"playlist-name"`
}

func (cfg *Config) GetPort() int {
	submatch := redirectURLRegexp.FindStringSubmatch(cfg.RedirectURL)
	port, _ := strconv.Atoi(submatch[1])
	return port
}

func (cfg *Config) GetCallback() string {
	submatch := redirectURLRegexp.FindStringSubmatch(cfg.RedirectURL)
	return submatch[2]
}

func ReadConfig() (*Config, error) {
	log.Println("Reading and validating config...")

	file, err := os.Open("config.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	err = decoder.Decode(&config)

	if err == nil {
		err = config.validate()
	}
	if err != nil {
		return nil, err
	}

	log.Println("Read valid config.")

	return &config, nil
}

func (cfg *Config) validate() error {
	if cfg.ClientID == "" {
		return errors.New("attribute 'client ID' must be specified")
	}
	if cfg.InputFilePath == "" {
		return errors.New("attribute 'input file path' must be specified")
	}
	if cfg.RedirectURL == "" {
		return errors.New("attribute 'redirect URL' must be specified")
	}
	if !redirectURLRegexp.MatchString(cfg.RedirectURL) {
		return fmt.Errorf("attribute 'redirect URL' does not match pattern '%s'", redirectURLRegexp)
	}
	port := cfg.GetPort()
	if port < 1024 || port > 65535 {
		return errors.New("port of redirect URL must be in interval [1024, 65535]")
	}
	if cfg.PlaylistName == "" {
		return errors.New("attribute 'playlist name' must be specified")
	}

	return nil
}
