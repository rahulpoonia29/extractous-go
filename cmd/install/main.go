// go run github.com/rahulpoonia29/extractous-go/cmd/install@latest
package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
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

var (
	verbose bool
	client  = http.DefaultClient
)

func main() {
	var platforms platformList
	var listPlatforms, downloadAll bool

	flag.Var(&platforms, "platform", "Specify a platform to download (e.g., linux_amd64). Can be used multiple times.")
	flag.BoolVar(&listPlatforms, "list", false, "List available platforms from the latest release and exit.")
	flag.BoolVar(&downloadAll, "all", false, "Download all available platforms from the latest release.")
	flag.BoolVar(&verbose, "v", false, "Verbose logging")
	flag.Parse()

	// use logging for errors and info (timestamps)
	log.SetFlags(0) // keep messages clean
	infof("Fetching Extractous FFI release information from GitHub...")

	availablePlatforms, err := getAvailablePlatforms()
	if err != nil {
		fatalf("Error retrieving available platforms: %v", err)
	}

	if listPlatforms {
		printAvailablePlatforms(availablePlatforms)
		return
	}

	platformsToDownload := determinePlatformsToDownload(platforms, downloadAll, availablePlatforms)
	if len(platformsToDownload) == 0 {
		infof("No platforms selected for download.")
		infof("Available platforms (run with --list to view):")
		printAvailablePlatforms(availablePlatforms)
		infof("To install for this machine run without flags, or pass --platform for the platform you want.")
		return
	}

	infof("Platforms selected for installation: %s", strings.Join(platformsToDownload, ", "))

	for _, platform := range platformsToDownload {
		archiveURL, ok := availablePlatforms[platform]
		if !ok {
			log.Printf("Warning: Platform '%s' not found in latest release. Skipping.", platform)
			continue
		}

		infof("Downloading release for platform: %s", platform)

		archivePath, err := downloadFileWithRetries(archiveURL, 3)
		if err != nil {
			fatalf("Failed to download asset for %s: %v", platform, err)
		}
		// ensure cleanup of downloaded archive
		defer os.Remove(archivePath)

		archiveFormat := "tar.gz"
		if strings.HasSuffix(archiveURL, ".zip") {
			archiveFormat = "zip"
		}

		if err := extractArchive(archivePath, nativeDir, platform, archiveFormat); err != nil {
			// attempt cleanup of partial extraction
			destPath := filepath.Join(nativeDir, platform)
			_ = os.RemoveAll(destPath)
			fatalf("Failed to extract archive for %s: %v", platform, err)
		}
		infof("Libraries for %s extracted to ./%s/%s", platform, nativeDir, platform)
	}

	infof("Installation completed successfully.")
}

func infof(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func printAvailablePlatforms(platforms map[string]string) {
	if len(platforms) == 0 {
		fmt.Println("  (no platforms found)")
		return
	}
	names := make([]string, 0, len(platforms))
	for n := range platforms {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("  - %s\n", name)
	}
}

func determinePlatformsToDownload(platforms platformList, downloadAll bool, availablePlatforms map[string]string) []string {
	if downloadAll {
		keys := make([]string, 0, len(availablePlatforms))
		for k := range availablePlatforms {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys
	}

	if len(platforms) > 0 {
		return platforms
	}

	currentPlatform, _ := getPlatformAndFormat()
	if _, ok := availablePlatforms[currentPlatform]; ok {
		return []string{currentPlatform}
	}

	// not found for current platform
	return []string{}
}

func getAvailablePlatforms() (map[string]string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)

	var resp *http.Response
	var err error

	// simple retry here too
	for attempt := 0; attempt < 3; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		wait := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		time.Sleep(wait)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query GitHub API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		remaining := resp.Header.Get("X-RateLimit-Remaining")
		reset := resp.Header.Get("X-RateLimit-Reset")
		if remaining == "0" && reset != "" {
			// parse reset as unix timestamp
			if ts, err := strconv.ParseInt(reset, 10, 64); err == nil {
				resetTime := time.Unix(ts, 0).Local()
				// Round to nearest minute for nicer display
				duration := max(time.Until(resetTime), 0)
				humanWait := fmt.Sprintf("about %d min", int(duration.Minutes()+0.5))

				return nil, fmt.Errorf(
					"GitHub API rate limit exceeded.\nLimit resets at: %s (%s from now)\nTip: set a personal access token to increase your limit",
					resetTime.Format("Mon 2 15:04 MST"),
					humanWait,
				)
			}
		}
		return nil, fmt.Errorf("access forbidden from GitHub API: %s", resp.Status)
	}

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
		if !strings.HasPrefix(asset.Name, "extractous-ffi-") {
			continue
		}
		// skip checksum assets like .sha256
		if strings.HasSuffix(asset.Name, ".sha256") || strings.HasSuffix(asset.Name, ".sha256.txt") {
			if verbose {
				log.Printf("Skipping checksum asset: %s", asset.Name)
			}
			continue
		}
		after := strings.TrimPrefix(asset.Name, "extractous-ffi-")
		name := strings.TrimSuffix(after, ".zip")
		name = strings.TrimSuffix(name, ".tar.gz")
		name = strings.TrimSuffix(name, ".tgz")
		platforms[name] = asset.DownloadURL
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

// downloadFileWithRetries will try a few times and show a progress bar.
func downloadFileWithRetries(url string, attempts int) (string, error) {
	var lastErr error
	for i := 1; i <= attempts; i++ {
		if i > 1 {
			// backoff
			backoff := time.Duration(i*i) * time.Second
			if verbose {
				log.Printf("Retrying in %s...", backoff)
			}
			time.Sleep(backoff)
		}
		path, err := downloadFile(url)
		if err == nil {
			return path, nil
		}
		lastErr = err
		if verbose {
			log.Printf("Attempt %d/%d failed: %v", i, attempts, err)
		}
	}
	return "", fmt.Errorf("download failed after %d attempts: %w", attempts, lastErr)
}

func downloadFile(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := client.Do(req)
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

	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionSetWidth(30),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionShowCount(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	if _, err = io.Copy(io.MultiWriter(tmpFile, bar), resp.Body); err != nil {
		return "", err
	}
	println("")

	return tmpFile.Name(), nil
}

func extractArchive(src, dest, platform, format string) error {
	destPath := filepath.Join(dest, platform)
	if err := os.MkdirAll(destPath, 0o755); err != nil {
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

// prevent zip-slip and path traversal by resolving absolute paths
func safeJoin(dest, name string) (string, error) {
	absDest, err := filepath.Abs(dest)
	if err != nil {
		return "", err
	}
	cleanName := filepath.Clean(strings.ReplaceAll(name, "\\", string(os.PathSeparator)))
	joined := filepath.Join(absDest, cleanName)
	absJoined, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}
	// allow the file to be exactly the dest dir or inside it
	if absJoined == absDest || strings.HasPrefix(absJoined, absDest+string(os.PathSeparator)) {
		return absJoined, nil
	}
	return "", fmt.Errorf("illegal file path outside destination: %s", name)
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// use forward slashes in zip entries; convert for local FS
		fname := filepath.FromSlash(f.Name)
		targetPath, err := safeJoin(dest, fname)
		if err != nil {
			return err
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, f.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
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

		// Clean header name to avoid path traversal
		name := header.Name
		if name == "" {
			continue
		}
		targetPath, err := safeJoin(dest, name)
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0o755); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
				return err
			}
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		case tar.TypeSymlink, tar.TypeLink:
			// skip symlinks for safety
			if verbose {
				log.Printf("Skipping symlink: %s", header.Name)
			}
		default:
			if verbose {
				log.Printf("Skipping unknown tar entry type %c for %s", header.Typeflag, header.Name)
			}
		}
	}
}
