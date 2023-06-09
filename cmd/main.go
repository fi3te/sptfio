package main

import (
	"context"
	"log"
	"time"

	"github.com/fi3te/sptfio/pkg/client"
	"github.com/fi3te/sptfio/pkg/config"
	"github.com/fi3te/sptfio/pkg/io"
	"github.com/zmb3/spotify/v2"
)

func main() {
	ctx := context.Background()

	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	sc, err := client.Login(ctx, cfg)
	if err != nil {
		panic(err)
	}

	user, err := sc.CurrentUser(ctx)
	if err != nil {
		log.Panicf("Cannot get current user: %v\n", err)
	}
	log.Printf("You are logged in as '%s'.", user.ID)

	log.Printf("Searching for playlist with name '%s'...\n", cfg.PlaylistName)
	playlist, err := client.FindPlaylist(ctx, sc, user.ID, cfg.PlaylistName)
	if err != nil {
		log.Printf("Playlist was not found (%v). Creating playlist with name '%s'...", err.Error(), cfg.PlaylistName)
		playlist, err = client.CreatePrivatePlaylist(ctx, sc, user.ID, cfg.PlaylistName, "")
	}
	if err != nil {
		log.Panicf("Cannot find or create playlist: %v\n", err)
	}

	trackIds, err := client.GetTrackIdsOfPlaylist(ctx, sc, playlist.Tracks)
	if err != nil {
		log.Panicf("Cannot fetch tracks: %v\n", err)
	}
	log.Printf("Number of tracks in playlist: %d\n", len(trackIds))

	lines, err := io.ReadLineByLine(cfg.InputFilePath, true)
	if err != nil {
		log.Panicf("Cannot read input file: %v\n", err)
	}

	for line := range *lines {
		lineIsID := client.IsSpotifyID(line)
		var track *spotify.FullTrack
		if lineIsID {
			log.Printf("Fetching track information for id '%s'...\n", line)
			track, err = client.GetTrack(ctx, sc, line)
		} else {
			log.Printf("Finding best matching track for query '%s'...\n", line)
			track, err = client.FindBestMatchingTrack(ctx, sc, line)
		}
		if err != nil {
			log.Printf("Cannot find track: %v\n", err)
			continue
		}

		if containsId(trackIds, track.ID) {
			log.Printf("Playlist already contains the track '%s'. Skipping...", track.Name)
			continue
		}

		err = client.AddTrack(ctx, sc, playlist, track)
		if err != nil {
			log.Printf("Cannot add track to playlist: %v\n", err)
		} else {
			log.Printf("Added track '%s' to playlist.\n", track.Name)
		}

		trackIds = append(trackIds, track.ID)

		time.Sleep(100 * time.Millisecond)
	}

	outputLines := make([]string, len(trackIds))
	for i := range trackIds {
		outputLines[i] = string(trackIds[i])
	}
	log.Println("Writing output file...")
	err = io.WriteLines(cfg.OutputFilePath, outputLines)
	if err != nil {
		log.Panicf("Cannot write output file: %v\n", err)
	}

	log.Println("Done.")
}

func containsId(ids []spotify.ID, id spotify.ID) bool {
	for _, value := range ids {
		if value == id {
			return true
		}
	}
	return false
}
