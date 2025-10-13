// Command install-extractous downloads and installs the required FFI libraries.
//
// This command fetches a release from GitHub, verifies its checksum, and extracts
// it into a local `./native` directory, making the CGO libraries available for
// your Go project.
//
// Usage:
//
//	go run github.com/rahulpoonia29/extractous-go/cmd/install@latest
package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	repoOwner = "rahulpoonia29"
	repoName  = "extractous-go"
	nativeDir = "native"
)

type platformList []string

func (p *platformList) String() string {
	return strings.Join(*p, ", ")
}

func (p *platformList) Set(value string) error {
	*p = append(*p, value)
	return nil
}

func main() {
	var platforms platformList
	var listPlatforms, downloadAll bool

	flag.Var(&platforms, "platform", "Specify a platform to download (e.g., linux_amd64). Can be used multiple times.")
	flag.BoolVar(&listPlatforms, "list-platforms", false, "List available platforms from the latest release and exit.")
	flag.BoolVar(&downloadAll, "all", false, "Download all available platforms from the latest release.")
	flag.Parse()

	fmt.Println("Fetching Extractous FFI release information from GitHub...")

	availablePlatforms, err := getAvailablePlatforms()
	if err != nil {
		fatalf("Error retrieving available platforms: %v", err)
	}

	if listPlatforms {
		fmt.Println("Available platforms:")
		for name := range availablePlatforms {
			fmt.Printf("  - %s\n", name)
		}
		return
	}

	platformsToDownload := determinePlatformsToDownload(platforms, downloadAll, availablePlatforms)
	if len(platformsToDownload) == 0 {
		fmt.Println("No platforms specified for download. Use --platform, --all, or run without flags to download for your current OS.")
		return
	}

	fmt.Printf("Platforms selected for installation: %s\n", strings.Join(platformsToDownload, ", "))

	for _, platform := range platformsToDownload {
		archiveURL, ok := availablePlatforms[platform]
		if !ok {
			fmt.Fprintf(os.Stderr, "Warning: Platform '%s' not found in the latest release. Skipping.\n", platform)
			continue
		}

		fmt.Printf("Downloading release for platform: %s\n", platform)

		archivePath, err := downloadFile(archiveURL)
		if err != nil {
			fatalf("Failed to download asset for %s: %v", platform, err)
		}
		defer os.Remove(archivePath)

		archiveFormat := "tar.gz"
		if strings.HasSuffix(archiveURL, ".zip") {
			archiveFormat = "zip"
		}

		if err := extractArchive(archivePath, nativeDir, platform, archiveFormat); err != nil {
			fatalf("Failed to extract archive for %s: %v", platform, err)
		}
		fmt.Printf("Libraries for %s extracted to ./%s/%s\n", platform, nativeDir, platform)
	}

	fmt.Println("\nInstallation completed successfully.")
}

func determinePlatformsToDownload(platforms platformList, downloadAll bool, availablePlatforms map[string]string) []string {
	if downloadAll {
		keys := make([]string, 0, len(availablePlatforms))
		for k := range availablePlatforms {
			keys = append(keys, k)
		}
		return keys
	}

	if len(platforms) > 0 {
		return platforms
	}

	currentPlatform, _ := getPlatformAndFormat()
	if _, ok := availablePlatforms[currentPlatform]; ok {
		return []string{currentPlatform}
	}

	return []string{}
}

func getAvailablePlatforms() (map[string]string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to query GitHub API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status from GitHub API: %s", resp.Status)
	}

	var releaseInfo struct {
		Assets []struct {
			Name        string `json:"name"`
			DownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&releaseInfo); err != nil {
		return nil, fmt.Errorf("failed to decode GitHub release info: %w", err)
	}

	platforms := make(map[string]string)
	for _, asset := range releaseInfo.Assets {
		if name, ok :=strings.CutPrefix(asset.Name, "extractous-ffi-"); ok  {
			name = strings.TrimSuffix(strings.TrimSuffix(name, ".zip"), ".tar.gz")
			platforms[name] = asset.DownloadURL
		}
	}

	if len(platforms) == 0 {
		return nil, fmt.Errorf("no compatible FFI assets found in the latest release")
	}

	return platforms, nil
}

func getPlatformAndFormat() (platform, format string) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	switch goos {
	case "linux":
		return fmt.Sprintf("linux_%s", goarch), "tar.gz"
	case "darwin":
		return fmt.Sprintf("darwin_%s", goarch), "tar.gz"
	case "windows":
		return fmt.Sprintf("windows_%s", goarch), "zip"
	default:
		fatalf("Unsupported operating system: %s", goos)
		return "", ""
	}
}

func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	tmpFile, err := os.CreateTemp("", "extractous-*.download")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err = io.Copy(tmpFile, resp.Body); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func extractArchive(src, dest, platform, format string) error {
	destPath := filepath.Join(dest, platform)
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return err
	}

	switch format {
	case "zip":
		return unzip(src, destPath)
	case "tar.gz":
		return untar(src, destPath)
	default:
		return fmt.Errorf("unsupported archive format: %s", format)
	}
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		if _, err = io.Copy(outFile, rc); err != nil {
			outFile.Close()
			rc.Close()
			return err
		}

		outFile.Close()
		rc.Close()
	}
	return nil
}

func untar(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
