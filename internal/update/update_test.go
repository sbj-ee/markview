package update

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testPub = "RWSLgkbIFSEuRS5xAEv9bP6Ma6NWSf7wuuDwLKeT++2GUiloG7OAZjVS"

func read(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestBlake2b512KnownAnswer(t *testing.T) {
	// RFC 7693 Appendix A.
	d := blake2b512([]byte("abc"))
	want := "ba80a53f981c4d0d6a2797b69f12f6e94c212f14685ac4b74b12bb6fdbffa2d17d87c5392aab792dc252d5de4533cc9518d38aa8dbf1925ab92386edd4009923"
	if got := hex.EncodeToString(d[:]); got != want {
		t.Fatalf("blake2b(abc) = %s", got)
	}
	long := blake2b512([]byte(strings.Repeat("a", 1000)))
	if hex.EncodeToString(long[:]) == hex.EncodeToString(d[:]) {
		t.Fatal("multi-block hash collided")
	}
}

// Fixtures made by the real minisign tool (copied from zwriter's tests).
func TestVerifyRealMinisignOutput(t *testing.T) {
	if err := VerifyMinisign(testPub, read(t, "SHA256SUMS"), read(t, "SHA256SUMS.minisig")); err != nil {
		t.Fatalf("prehashed: %v", err)
	}
	if err := VerifyMinisign(testPub, read(t, "SHA256SUMS.legacy"), read(t, "SHA256SUMS.legacy.minisig")); err != nil {
		t.Fatalf("legacy: %v", err)
	}
}

func TestVerifyRejects(t *testing.T) {
	sums, sig := read(t, "SHA256SUMS"), read(t, "SHA256SUMS.minisig")
	tampered := append([]byte{}, sums...)
	tampered[0] ^= 1
	if VerifyMinisign(testPub, tampered, sig) == nil {
		t.Error("tampered SHA256SUMS accepted")
	}
	if VerifyMinisign(SigningPublicKey, sums, sig) == nil {
		t.Error("signature by another key accepted")
	}
	badTC := strings.Replace(string(sig), "timestamp:0", "timestamp:1", 1)
	if VerifyMinisign(testPub, sums, []byte(badTC)) == nil {
		t.Error("edited trusted comment accepted")
	}
	if VerifyMinisign(testPub, sums, []byte("garbage")) == nil {
		t.Error("garbage accepted")
	}
	if VerifyMinisign("not a key", sums, sig) == nil {
		t.Error("bad key accepted")
	}
}

func TestEmbeddedKeyIsTheSharedReleaseKey(t *testing.T) {
	pk, err := base64.StdEncoding.DecodeString(SigningPublicKey)
	if err != nil || len(pk) != 42 {
		t.Fatal("embedded key does not decode")
	}
	id := make([]byte, 8)
	for i := 0; i < 8; i++ {
		id[i] = pk[9-i]
	}
	if got := strings.ToUpper(hex.EncodeToString(id)); got != "DBC46B23C3AF0EF3" {
		t.Fatalf("key id %s, want DBC46B23C3AF0EF3", got)
	}
}

func TestChecksumFor(t *testing.T) {
	h := strings.Repeat("ab", 32)
	sums := []byte(h + "  markview_1.1.0_amd64.deb\n" + strings.Repeat("cd", 32) + " *MarkView-1.1.0.dmg\n")
	if ChecksumFor(sums, "markview_1.1.0_amd64.deb") != h {
		t.Error("deb line")
	}
	if ChecksumFor(sums, "MarkView-1.1.0.dmg") != strings.Repeat("cd", 32) {
		t.Error("binary-mode line")
	}
	if ChecksumFor(sums, "other") != "" {
		t.Error("unlisted file")
	}
}

func TestPackageNameAndInstall(t *testing.T) {
	if PackageName("1.1.0", "linux", "amd64") != "markview_1.1.0_amd64.deb" ||
		PackageName("1.1.0", "darwin", "arm64") != "MarkView-1.1.0.dmg" ||
		PackageName("1.1.0", "windows", "amd64") != "" {
		t.Error("package names")
	}
	if strings.Join(InstallCommand("linux", "/tmp/x.deb"), " ") != "/usr/bin/pkexec /usr/bin/apt install -y /tmp/x.deb" {
		t.Error("linux install command")
	}
	if InstallCommand("darwin", "/tmp/x.dmg")[0] != "/usr/bin/open" {
		t.Error("mac open command")
	}
}

func TestTrustedDownloadURL(t *testing.T) {
	ok := "https://github.com/sbj-ee/markview/releases/download/v1.1.0/SHA256SUMS"
	if !TrustedDownloadURL(ok) {
		t.Error(ok)
	}
	for _, bad := range []string{"http://github.com/sbj-ee/markview/releases/download/v1/x",
		"https://evil.example/sbj-ee/markview/releases/download/v1/x",
		"https://github.com/someone/markview/releases/download/v1/x"} {
		if TrustedDownloadURL(bad) {
			t.Error(bad)
		}
	}
}

// signMinisign writes a prehashed minisign signature with a throwaway key.
func signMinisign(priv ed25519.PrivateKey, keyID []byte, msg []byte, trusted string) []byte {
	d := blake2b512(msg)
	s := ed25519.Sign(priv, d[:])
	sig := append(append([]byte("ED"), keyID...), s...)
	global := ed25519.Sign(priv, append(append([]byte{}, s...), trusted...))
	return []byte("untrusted comment: test\n" + base64.StdEncoding.EncodeToString(sig) +
		"\ntrusted comment: " + trusted + "\n" + base64.StdEncoding.EncodeToString(global) + "\n")
}

func mockGitHub(t *testing.T, pkg []byte, listed []byte, signWith ed25519.PrivateKey, keyID []byte) (*httptest.Server, *Release) {
	sum := sha256.Sum256(listed)
	sums := []byte(hex.EncodeToString(sum[:]) + "  markview_1.1.0_amd64.deb\n")
	sig := signMinisign(signWith, keyID, sums, "timestamp:0 file:SHA256SUMS hashed")
	mux := http.NewServeMux()
	mux.HandleFunc("/SHA256SUMS", func(w http.ResponseWriter, _ *http.Request) { w.Write(sums) })
	mux.HandleFunc("/SHA256SUMS.minisig", func(w http.ResponseWriter, _ *http.Request) { w.Write(sig) })
	mux.HandleFunc("/pkg", func(w http.ResponseWriter, _ *http.Request) { w.Write(pkg) })
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	r := &Release{TagName: "v1.1.0", Assets: []Asset{
		{Name: "markview_1.1.0_amd64.deb", URL: srv.URL + "/pkg"},
		{Name: "SHA256SUMS", URL: srv.URL + "/SHA256SUMS"},
		{Name: "SHA256SUMS.minisig", URL: srv.URL + "/SHA256SUMS.minisig"},
	}}
	return srv, r
}

func TestDownloadVerifiedFlow(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	keyID := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	pubB64 := base64.StdEncoding.EncodeToString(append(append([]byte("Ed"), keyID...), pub...))
	pkg := []byte(strings.Repeat("fake markview package\n", 500))
	client := func(srv *httptest.Server, key string) *Client {
		return &Client{HTTP: srv.Client(), PublicKey: key, AllowURL: func(u string) bool { return strings.HasPrefix(u, srv.URL) }}
	}

	srv, rel := mockGitHub(t, pkg, pkg, priv, keyID)
	dir := t.TempDir()
	path, err := client(srv, pubB64).DownloadVerified(context.Background(), rel, "markview_1.1.0_amd64.deb", dir)
	if err != nil {
		t.Fatalf("good release refused: %v", err)
	}
	if b, _ := os.ReadFile(path); string(b) != string(pkg) {
		t.Fatal("package bytes differ")
	}

	// Package doesn't match the signed checksum: refused, nothing left.
	srv2, rel2 := mockGitHub(t, append(pkg, 'x'), pkg, priv, keyID)
	dir2 := t.TempDir()
	if _, err := client(srv2, pubB64).DownloadVerified(context.Background(), rel2, "markview_1.1.0_amd64.deb", dir2); err == nil {
		t.Fatal("mismatched package accepted")
	}
	if ents, _ := os.ReadDir(dir2); len(ents) != 0 {
		t.Fatal("rejected package left on disk")
	}

	// Signed by another key: refused before the package is fetched.
	_, other, _ := ed25519.GenerateKey(nil)
	srv3, rel3 := mockGitHub(t, pkg, pkg, other, keyID)
	if _, err := client(srv3, pubB64).DownloadVerified(context.Background(), rel3, "markview_1.1.0_amd64.deb", t.TempDir()); err == nil ||
		!strings.Contains(err.Error(), "signature") {
		t.Fatalf("wrong-key signature accepted: %v", err)
	}

	// Unsigned release: refused.
	rel4 := &Release{TagName: "v1.1.0", Assets: rel.Assets[:1]}
	if _, err := client(srv, pubB64).DownloadVerified(context.Background(), rel4, "markview_1.1.0_amd64.deb", t.TempDir()); err == nil {
		t.Fatal("unsigned release accepted")
	}

	// Default URL policy rejects the mock (non-GitHub) host.
	c := &Client{HTTP: srv.Client(), PublicKey: pubB64}
	if _, err := c.DownloadVerified(context.Background(), rel, "markview_1.1.0_amd64.deb", t.TempDir()); err == nil {
		t.Fatal("untrusted host accepted")
	}
}

func TestLatest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"tag_name":"v1.1.0","html_url":"https://github.com/sbj-ee/markview/releases/tag/v1.1.0","assets":[]}`))
	}))
	defer srv.Close()
	r, err := (&Client{HTTP: srv.Client(), APIURL: srv.URL}).Latest(context.Background())
	if err != nil || r.Version() != "1.1.0" {
		t.Fatalf("%v %v", r, err)
	}
}
