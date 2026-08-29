package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/codewithsorabh/adbot/internal/version"
)

const (
	RepoOwner       = "codewithsorabh"
	RepoName        = "adbot"
	ReleaseAPIURL   = "https://api.github.com/repos/codewithsorabh/adbot/releases/latest"
	FallbackVersion = "https://raw.githubusercontent.com/codewithsorabh/adbot/main/VERSION"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// CheckForUpdate checks if a newer version is available.
// If verbose is true, outputs status information; if false, reports only if an update is found.
func CheckForUpdate(verbose bool) {
	current := version.Get()
	latestTag, downloadURL, err := FetchLatestReleaseInfo()
	if err != nil {
		if verbose {
			fmt.Printf("[ERROR] Failed to check for updates: %v\n", err)
		}
		return
	}

	cleanLatest := strings.TrimPrefix(latestTag, "v")
	cleanCurrent := strings.TrimPrefix(current, "v")

	if cleanLatest == "" || cleanLatest == cleanCurrent || current == "dev" {
		if verbose {
			fmt.Printf("Already up to date (%s)\n", version.FullInfo())
		}
		return
	}

	if !verbose {
		fmt.Printf("A new version of ADBOT is available: v%s -> v%s\n", cleanCurrent, cleanLatest)
		fmt.Println("Run 'adbot --check-update' to update.")
		fmt.Println()
		return
	}

	fmt.Printf("Updating ADBOT: v%s -> v%s\n", cleanCurrent, cleanLatest)

	if downloadURL == "" {
		fmt.Printf("[ERROR] No prebuilt binary found for %s/%s in release %s\n", runtime.GOOS, runtime.GOARCH, latestTag)
		return
	}

	if err := DownloadAndReplaceSelf(downloadURL); err != nil {
		fmt.Printf("[ERROR] Failed to update ADBOT: %v\n", err)
		return
	}

	fmt.Printf("ADBOT updated successfully to v%s\n", cleanLatest)
	fmt.Println("Run 'adbot' again to continue.")
}

// FetchLatestReleaseInfo queries GitHub Releases API for the latest release tag and download URL.
func FetchLatestReleaseInfo() (string, string, error) {
	client := &http.Client{
		Timeout: 7 * time.Second,
	}

	req, err := http.NewRequest("GET", ReleaseAPIURL, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("User-Agent", "adbot/"+version.Get())
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			resp.Body.Close()
		}
		// Fallback to raw VERSION file
		tag, fbErr := fetchFallbackVersion(client)
		if fbErr != nil {
			if err != nil {
				return "", "", err
			}
			return "", "", fmt.Errorf("HTTP error: %s", resp.Status)
		}
		return tag, "", nil
	}
	defer resp.Body.Close()

	var rel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", "", err
	}

	targetName := fmt.Sprintf("adbot-%s-%s-%s", rel.TagName, runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		targetName += ".exe"
	}

	var downloadURL string
	for _, asset := range rel.Assets {
		if asset.Name == targetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	return rel.TagName, downloadURL, nil
}

func fetchFallbackVersion(client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", FallbackVersion, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "adbot/"+version.Get())

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fallback request failed: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

// DownloadAndReplaceSelf downloads a binary from url and safely replaces the currently running executable.
func DownloadAndReplaceSelf(downloadURL string) error {
	client := &http.Client{
		Timeout: 90 * time.Second,
	}

	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "adbot/"+version.Get())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(execPath), ".adbot-update-*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return err
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	if err := os.Chmod(tmpPath, 0755); err != nil {
		os.Remove(tmpPath)
		return err
	}

	if err := os.Rename(tmpPath, execPath); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return nil
}
