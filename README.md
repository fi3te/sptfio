# Sptfio

## Usage

1. Visit developer.spotify.com and create an app.
2. Enter the client ID and redirect URL in the config.yml file.
3. Specify the name of a playlist to read from and/or write to.
4. Specify the path to a file with songs (search strings or Spotify IDs; empty file if no songs are to be added).
5. Specify the path to a file where the Spotify IDs of all songs of the playlist are to be written to.
6. Run the app with `go run .\cmd\main.go` and follow the instructions.