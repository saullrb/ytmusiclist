package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/saullbrandao/ytmusiclist/dependencies"
	"github.com/saullbrandao/ytmusiclist/utils"
)

type UserInput struct {
	playlistUrl  string
	dirName      string
	convertToMp3 bool
}

func main() {
	err := dependencies.EnsureFFMPEG()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		utils.GracefulExit()
	}

	ytdlpPath, err := dependencies.EnsureYTDLP()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		utils.GracefulExit()
	}

	for {
		userInput, err := getUserInput(ytdlpPath)
		if err != nil {
			fmt.Println(err)
			utils.GracefulExit()
		}

		err = downloadPlaylist(ytdlpPath, userInput)
		if err != nil {
			fmt.Println(err)
			utils.GracefulExit()
		}

		fmt.Println("\nPlaylist downloaded!")
	}

}

func getUserInput(ytdlpPath string) (UserInput, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter YouTube playlist URL: ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)

	if url == "" {
		return UserInput{}, errors.New("Error: URL required")
	}

	playlistTitle, err := getPlaylistTitle(ytdlpPath, url)
	if err != nil {
		playlistTitle = "playlist"
	}

	fmt.Printf("Enter the directory name you want to use (Leave empty to use \"%s\"): ", playlistTitle)
	dirName, _ := reader.ReadString('\n')
	dirName = strings.TrimSpace(dirName)

	if dirName == "" {
		dirName = playlistTitle
	}

	if dirName == "" {
		dirName = "."
	}

	fmt.Print("Do you want to convert to MP3? (y/N): ")
	mp3Input, _ := reader.ReadString('\n')
	mp3Input = strings.ToLower(strings.TrimSpace(mp3Input))

	convertToMp3 := false
	if mp3Input == "y" {
		convertToMp3 = true
	}

	return UserInput{playlistUrl: url, dirName: dirName, convertToMp3: convertToMp3}, nil
}

func getPlaylistTitle(ytdlpPath, url string) (string, error) {
	cmd := exec.Command(ytdlpPath, "--print", "%(playlist_title)s", "--flat-playlist", "-i", url)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	title := lines[0]
	title = sanitizeFileName(title)
	return title, nil
}

func sanitizeFileName(name string) string {
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := name
	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}
	return strings.TrimSpace(result)
}
