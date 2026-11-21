package update

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRunningBundle(t *testing.T) {
	if got := RunningBundle("/Applications/MarkView.app/Contents/MacOS/markview"); got != "/Applications/MarkView.app" {
		t.Errorf("RunningBundle = %q", got)
	}
	for _, exe := range []string{"/usr/bin/markview", "/tmp/MacOS/markview", "/x/Contents/MacOS/markview"} {
		if got := RunningBundle(exe); got != "" {
			t.Errorf("RunningBundle(%q) = %q, want empty", exe, got)
		}
	}
}

// fakeRunner emulates hdiutil attach (by populating the mountpoint with a
// MarkView.app holding marker) and ditto (by copying that marker).
func fakeRunner(t *testing.T, marker string, dittoErr error) (Runner, *[]string) {
	var calls []string
	return func(_ time.Duration, name string, args ...string) (string, error) {
		calls = append(calls, filepath.Base(name)+" "+args[0])
		switch {
		case name == "/usr/bin/hdiutil" && args[0] == "attach":
			mnt := args[len(args)-2]
			app := filepath.Join(mnt, AppBundleName)
			os.MkdirAll(app, 0o755)
			return "", os.WriteFile(filepath.Join(app, "version"), []byte(marker), 0o644)
		case name == "/usr/bin/ditto":
			if dittoErr != nil {
				return "boom", dittoErr
			}
			os.MkdirAll(args[1], 0o755)
			b, _ := os.ReadFile(filepath.Join(args[0], "version"))
			return "", os.WriteFile(filepath.Join(args[1], "version"), b, 0o644)
		}
		return "", nil
	}, &calls
}

func newBundle(t *testing.T, marker string) string {
	bundle := filepath.Join(t.TempDir(), "MarkView.app")
	os.MkdirAll(bundle, 0o755)
	os.WriteFile(filepath.Join(bundle, "version"), []byte(marker), 0o644)
	return bundle
}

func TestReplaceBundleFromDMG(t *testing.T) {
	bundle := newBundle(t, "old")
	run, calls := fakeRunner(t, "new", nil)
	if err := ReplaceBundleFromDMG("x.dmg", bundle, t.TempDir(), run); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(bundle, "version")); string(b) != "new" {
		t.Fatalf("bundle has %q, want new", b)
	}
	for _, leftover := range []string{bundle + ".new", bundle + ".old"} {
		if _, err := os.Stat(leftover); err == nil {
			t.Errorf("%s left behind", leftover)
		}
	}
	if last := (*calls)[len(*calls)-1]; last != "hdiutil detach" {
		t.Errorf("image not detached; calls %v", *calls)
	}
}

func TestReplaceBundleFromDMGKeepsOldOnFailure(t *testing.T) {
	bundle := newBundle(t, "old")
	run, _ := fakeRunner(t, "new", errors.New("ditto failed"))
	if err := ReplaceBundleFromDMG("x.dmg", bundle, t.TempDir(), run); err == nil {
		t.Fatal("failure not reported")
	}
	if b, _ := os.ReadFile(filepath.Join(bundle, "version")); string(b) != "old" {
		t.Fatalf("old bundle damaged: %q", b)
	}
}
