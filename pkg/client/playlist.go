package client

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
)

func FindPlaylist(ctx context.Context, sc *spotify.Client, userID string, playlistName string) (*spotify.FullPlaylist, error) {
	playlistPage, err := sc.GetPlaylistsForUser(ctx, userID)

	if err != nil {
		return nil, err
	}
	for {
		for _, playlist := range playlistPage.Playlists {
			if playlist.Name == playlistName {
				return sc.GetPlaylist(ctx, playlist.ID)
			}
		}

		err = sc.NextPage(ctx, playlistPage)
		if err != nil {
			if err == spotify.ErrNoMorePages {
				break
			}
			return nil, err
		}
	}

	return nil, fmt.Errorf("playlist with name '%s' does not exist", playlistName)
}

func CreatePrivatePlaylist(ctx context.Context, sc *spotify.Client, userID string, name string, description string) (*spotify.FullPlaylist, error) {
	return sc.CreatePlaylistForUser(ctx, userID, name, description, false, false)
}

func GetTrackIdsOfPlaylist(ctx context.Context, sc *spotify.Client, page spotify.PlaylistTrackPage) ([]spotify.ID, error) {
	var trackIds []spotify.ID
	for {
		for _, track := range page.Tracks {
			trackIds = append(trackIds, track.Track.ID)
		}

		err := sc.NextPage(ctx, &page)
		if err != nil {
			if err == spotify.ErrNoMorePages {
				break
			}
			return nil, err
		}
	}
	return trackIds, nil
}

func AddTrack(ctx context.Context, sc *spotify.Client, playlist *spotify.FullPlaylist, track *spotify.FullTrack) error {
	_, err := sc.AddTracksToPlaylist(ctx, playlist.ID, track.ID)
	if err != nil {
		return err
	}
	return nil
}
