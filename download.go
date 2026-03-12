package main

import (
	"fmt"
	"os"
	"os/exec"
)

func downloadPlaylist(ytdlpPath string, userInput UserInput) error {
	if err := os.MkdirAll(userInput.dirName, 0755); err != nil {
		return err
	}

	audioFormat := "mp3"
	if !userInput.convertToMp3 {
		audioFormat = "best"
	}

	cmd := exec.Command(
		ytdlpPath,
		"-x",
		"--audio-format",
		audioFormat,
		"--quiet",
		"--progress",
		"--no-simulate",
		"--print",
		"Downloading item %(playlist_index)s of %(playlist_count)s - %(title)s",
		"--download-archive",
		userInput.dirName+"/downloaded.txt",
		"-o",
		userInput.dirName+"/%(title)s [%(id)s].%(ext)s",
		"--no-post-overwrites",
		userInput.playlistUrl,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("\nDownloading playlist to \"%s\" folder\n", userInput.dirName)
	return cmd.Run()
}
