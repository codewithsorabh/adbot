package updater

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
)

func TestParseReleaseAsset(t *testing.T) {
	tag := "v1.0.0"
	targetName := "adbot-" + tag + "-" + runtime.GOOS + "-" + runtime.GOARCH
	if runtime.GOOS == "windows" {
		targetName += ".exe"
	}

	mockRel := githubRelease{
		TagName: tag,
		Assets: []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		}{
			{
				Name:               targetName,
				BrowserDownloadURL: "https://example.com/download/" + targetName,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockRel)
	}))
	defer server.Close()

	if len(mockRel.Assets) == 0 || mockRel.Assets[0].Name != targetName {
		t.Fatalf("expected asset %s, got %v", targetName, mockRel.Assets)
	}
}
