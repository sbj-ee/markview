// Package update is MarkView's signed in-app updater logic, following
// zwriter's flow: read the latest GitHub release, pick the platform package,
// download SHA256SUMS and SHA256SUMS.minisig, verify the minisign signature
// with the embedded public key, then download the package and check its
// SHA-256 against the signed list. No GUI here, so it is unit-tested.
package update

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// LatestReleaseURL is GitHub's public "latest release" endpoint for MarkView.
const LatestReleaseURL = "https://api.github.com/repos/sbj-ee/markview/releases/latest"

// Asset is one release file.
type Asset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
	Size int64  `json:"size"`
}

// Release is the part of GitHub's release JSON the updater uses.
type Release struct {
	TagName     string  `json:"tag_name"`
	Name        string  `json:"name"`
	Body        string  `json:"body"`
	HTMLURL     string  `json:"html_url"`
	PublishedAt string  `json:"published_at"`
	Draft       bool    `json:"draft"`
	Prerelease  bool    `json:"prerelease"`
	Assets      []Asset `json:"assets"`
}

// Version is the tag without its leading "v".
func (r *Release) Version() string { return strings.TrimPrefix(r.TagName, "v") }

// PackageName is the release asset the updater installs on goos/goarch:
// markview_<v>_<arch>.deb on Linux (scripts/build-deb.sh) and
// MarkView-<v>.dmg on macOS (scripts/build-dmg.sh). Empty elsewhere.
func PackageName(version, goos, goarch string) string {
	switch goos {
	case "linux":
		return fmt.Sprintf("markview_%s_%s.deb", version, goarch)
	case "darwin":
		return fmt.Sprintf("MarkView-%s.dmg", version)
	}
	return ""
}

// Find returns the asset called name.
func (r *Release) Find(name string) (Asset, bool) {
	for _, a := range r.Assets {
		if a.Name == name {
			return a, true
		}
	}
	return Asset{}, false
}

// ParseRelease decodes GitHub's release JSON.
func ParseRelease(data []byte) (*Release, error) {
	var r Release
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("unreadable release JSON: %w", err)
	}
	if r.TagName == "" {
		return nil, errors.New("release has no tag")
	}
	return &r, nil
}

// ChecksumFor returns the lower-case hex SHA-256 listed for name in a
// sha256sum-style file ("<hex>  <name>" or "<hex> *<name>"), or "".
func ChecksumFor(sums []byte, name string) string {
	for _, line := range strings.Split(string(sums), "\n") {
		f := strings.Fields(strings.TrimSpace(line))
		if len(f) != 2 || len(f[0]) != 64 {
			continue
		}
		if strings.TrimPrefix(f[1], "*") == name {
			if _, err := hex.DecodeString(f[0]); err == nil {
				return strings.ToLower(f[0])
			}
		}
	}
	return ""
}

// ExpectedSignedComment is the trusted comment the release workflow signs
// SHA256SUMS with. Requiring it binds the signature to one release, so an
// older release's signed SHA256SUMS can't be replayed under a newer tag.
func ExpectedSignedComment(tag string) string {
	return "markview " + tag + " SHA256SUMS"
}

// TrustedAssetURL accepts only the exact GitHub download URL for asset name
// of release tag.
func TrustedAssetURL(raw, tag, name string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Host != "github.com" || u.User != nil || u.Port() != "" {
		return false
	}
	return u.Path == "/sbj-ee/markview/releases/download/"+tag+"/"+name
}

// TrustedAssetRedirect accepts the hosts GitHub redirects release downloads
// to. A download redirected anywhere else is refused.
func TrustedAssetRedirect(u *url.URL) bool {
	return u != nil && u.Scheme == "https" && u.User == nil && u.Port() == "" &&
		(u.Host == "release-assets.githubusercontent.com" || u.Host == "objects.githubusercontent.com")
}

// TrustedDownloadURL accepts only https GitHub release download URLs for
// this repository (GitHub then redirects to its object store).
func TrustedDownloadURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Host != "github.com" {
		return false
	}
	return strings.HasPrefix(u.Path, "/sbj-ee/markview/releases/download/")
}

// Client fetches release metadata and assets.
type Client struct {
	HTTP      *http.Client
	APIURL    string // default LatestReleaseURL
	PublicKey string // default SigningPublicKey
	AllowURL  func(string) bool
	// AllowRedirect decides whether an asset download may follow a
	// redirect; default TrustedAssetRedirect.
	AllowRedirect func(*url.URL) bool
	UserAgent     string
	OnProgress    func(name string, got, total int64)
}

func (c *Client) get(ctx context.Context, u string) (*http.Response, error) {
	return c.do(ctx, c.HTTP, u)
}

// getAsset downloads a release asset, following redirects only to trusted
// asset hosts.
func (c *Client) getAsset(ctx context.Context, u string) (*http.Response, error) {
	hc := *c.HTTP
	allow := c.AllowRedirect
	if allow == nil {
		allow = TrustedAssetRedirect
	}
	hc.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return errors.New("too many redirects")
		}
		if !allow(req.URL) {
			return fmt.Errorf("redirected to %s, which is not a GitHub release-asset host", req.URL.Redacted())
		}
		return nil
	}
	return c.do(ctx, &hc, u)
}

func (c *Client) do(ctx context.Context, hc *http.Client, u string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	ua := c.UserAgent
	if ua == "" {
		ua = "MarkView-Updater"
	}
	req.Header.Set("User-Agent", ua)
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			return nil, errors.New("no published MarkView release was found on GitHub")
		}
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
			return nil, errors.New("GitHub's rate limit was reached; try again later")
		}
		return nil, fmt.Errorf("GitHub returned HTTP %d", resp.StatusCode)
	}
	return resp, nil
}

// Latest fetches the latest published release.
func (c *Client) Latest(ctx context.Context) (*Release, error) {
	u := c.APIURL
	if u == "" {
		u = LatestReleaseURL
	}
	resp, err := c.get(ctx, u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	return ParseRelease(data)
}

func (c *Client) fetchSmall(ctx context.Context, tag string, a Asset) ([]byte, error) {
	if !c.allowed(a.URL, tag, a.Name) {
		return nil, fmt.Errorf("refusing to download %s from %s", a.Name, a.URL)
	}
	resp, err := c.getAsset(ctx, a.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(io.LimitReader(resp.Body, 1<<20))
}

func (c *Client) allowed(u, tag, name string) bool {
	if c.AllowURL != nil {
		return c.AllowURL(u)
	}
	return TrustedAssetURL(u, tag, name)
}

// DownloadVerified downloads SHA256SUMS and SHA256SUMS.minisig, verifies the
// signature against the embedded public key before fetching anything else,
// then downloads pkgName into dir and checks its SHA-256 against the signed
// list. It returns the verified package path; on any failure nothing is left
// behind for an installer to pick up.
func (c *Client) DownloadVerified(ctx context.Context, r *Release, pkgName, dir string) (string, error) {
	path, _, err := c.DownloadVerifiedSum(ctx, r, pkgName, dir)
	return path, err
}

// DownloadVerifiedSum is DownloadVerified that also returns the signed
// SHA-256 of the package, for re-checking it just before install.
func (c *Client) DownloadVerifiedSum(ctx context.Context, r *Release, pkgName, dir string) (string, string, error) {
	pub := c.PublicKey
	if pub == "" {
		pub = SigningPublicKey
	}
	pkg, ok := r.Find(pkgName)
	if !ok {
		return "", "", fmt.Errorf("release %s has no %s", r.TagName, pkgName)
	}
	sumsA, ok1 := r.Find("SHA256SUMS")
	sigA, ok2 := r.Find("SHA256SUMS.minisig")
	if !ok1 || !ok2 {
		return "", "", fmt.Errorf("release %s is not signed (SHA256SUMS or SHA256SUMS.minisig missing)", r.TagName)
	}
	sums, err := c.fetchSmall(ctx, r.TagName, sumsA)
	if err != nil {
		return "", "", fmt.Errorf("downloading SHA256SUMS: %w", err)
	}
	sig, err := c.fetchSmall(ctx, r.TagName, sigA)
	if err != nil {
		return "", "", fmt.Errorf("downloading SHA256SUMS.minisig: %w", err)
	}
	comment, err := VerifyMinisignComment(pub, sums, sig)
	if err != nil {
		return "", "", fmt.Errorf("SHA256SUMS signature is not valid: %w", err)
	}
	if want := ExpectedSignedComment(r.TagName); comment != want {
		return "", "", fmt.Errorf("SHA256SUMS is signed for %q, not %q", comment, want)
	}
	want := ChecksumFor(sums, pkgName)
	if want == "" {
		return "", "", fmt.Errorf("%s is not listed in the signed SHA256SUMS", pkgName)
	}
	if !c.allowed(pkg.URL, r.TagName, pkg.Name) {
		return "", "", fmt.Errorf("refusing to download %s from %s", pkg.Name, pkg.URL)
	}
	resp, err := c.getAsset(ctx, pkg.URL)
	if err != nil {
		return "", "", fmt.Errorf("downloading %s: %w", pkgName, err)
	}
	defer resp.Body.Close()
	path := filepath.Join(dir, filepath.Base(pkgName))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return "", "", err
	}
	h := sha256.New()
	var got int64
	buf := make([]byte, 64<<10)
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := f.Write(buf[:n]); werr != nil {
				f.Close()
				os.Remove(path)
				return "", "", werr
			}
			h.Write(buf[:n])
			got += int64(n)
			if c.OnProgress != nil {
				c.OnProgress(pkgName, got, resp.ContentLength)
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			f.Close()
			os.Remove(path)
			return "", "", fmt.Errorf("downloading %s: %w", pkgName, rerr)
		}
	}
	if err := f.Close(); err != nil {
		os.Remove(path)
		return "", "", err
	}
	if have := hex.EncodeToString(h.Sum(nil)); have != want {
		os.Remove(path)
		return "", "", fmt.Errorf("%s does not match the signed SHA256SUMS (expected %s, got %s)", pkgName, want, have)
	}
	return path, want, nil
}

// StageVerified copies the package at src into dir while hashing it and
// returns the copy's path only if it still matches the signed SHA-256. The
// installer is handed that private, freshly re-verified copy.
func StageVerified(src, wantSHA256, dir string) (string, error) {
	in, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer in.Close()
	dst := filepath.Join(dir, filepath.Base(src))
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	_, cerr := io.Copy(io.MultiWriter(out, h), in)
	if err := out.Close(); cerr == nil {
		cerr = err
	}
	if cerr != nil {
		os.Remove(dst)
		return "", cerr
	}
	if have := hex.EncodeToString(h.Sum(nil)); have != wantSHA256 {
		os.Remove(dst)
		return "", fmt.Errorf("%s changed after it was verified (expected %s, got %s)", filepath.Base(src), wantSHA256, have)
	}
	return dst, nil
}

// InstallCommand is how a verified .deb is installed on Linux: through
// pkexec apt (the user is asked for their password). macOS installs in
// place with ReplaceBundleFromDMG instead, so it has no command.
func InstallCommand(goos, pkgPath string) []string {
	if goos == "linux" {
		return []string{"/usr/bin/pkexec", "/usr/bin/apt", "install", "-y", pkgPath}
	}
	return nil
}
