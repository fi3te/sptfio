package client

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
)

func FindBestMatchingTrack(ctx context.Context, sc *spotify.Client, query string) (*spotify.FullTrack, error) {
	searchResult, err := sc.Search(ctx, query, spotify.SearchTypeTrack)
	if err != nil {
		return nil, err
	}
	trackPage := searchResult.Tracks
	numberOfResults := trackPage.Limit
	if numberOfResults == 0 {
		return nil, fmt.Errorf("cannot find track for query '%s'", query)
	}
	return &trackPage.Tracks[0], nil
}

func GetTrack(ctx context.Context, sc *spotify.Client, id string) (*spotify.FullTrack, error) {
	return sc.GetTrack(ctx, spotify.ID(id))
}
