package update

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// AppBundleName is the bundle inside the release DMG (scripts/build-dmg.sh).
const AppBundleName = "MarkView.app"

// RunningBundle returns the .app bundle that contains the executable exe
// (…/MarkView.app/Contents/MacOS/markview), or "" if exe isn't inside one.
func RunningBundle(exe string) string {
	macOS := filepath.Dir(exe)
	contents := filepath.Dir(macOS)
	bundle := filepath.Dir(contents)
	if filepath.Base(macOS) != "MacOS" || filepath.Base(contents) != "Contents" || !strings.HasSuffix(bundle, ".app") {
		return ""
	}
	return bundle
}

// Runner runs a command and returns its combined output. Tests replace it.
type Runner func(timeout time.Duration, name string, args ...string) (string, error)

// ExecRunner runs commands for real.
func ExecRunner(timeout time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// ReplaceBundleFromDMG installs MarkView.app from a verified dmg in place of
// bundle: mount the image read-only, ditto the app next to bundle as
// bundle.new, move bundle aside to bundle.old, move bundle.new into place,
// then delete bundle.old. If the final move fails the old bundle is put
// back, so a failure never leaves MarkView missing.
func ReplaceBundleFromDMG(dmg, bundle, workDir string, run Runner) error {
	if run == nil {
		run = ExecRunner
	}
	mnt := filepath.Join(workDir, "mnt")
	if err := os.MkdirAll(mnt, 0o700); err != nil {
		return err
	}
	if out, err := run(2*time.Minute, "/usr/bin/hdiutil", "attach", "-nobrowse", "-readonly", "-noautoopen", "-mountpoint", mnt, dmg); err != nil {
		return fmt.Errorf("couldn't open the disk image: %v %s", err, out)
	}
	defer run(time.Minute, "/usr/bin/hdiutil", "detach", mnt, "-quiet") //nolint:errcheck

	src := filepath.Join(mnt, AppBundleName)
	if fi, err := os.Stat(src); err != nil || !fi.IsDir() {
		return fmt.Errorf("the disk image has no %s", AppBundleName)
	}
	staged := bundle + ".new"
	old := bundle + ".old"
	os.RemoveAll(staged)
	os.RemoveAll(old)
	defer os.RemoveAll(staged)

	// ditto keeps the bundle's symlinks, permissions and code signature.
	if out, err := run(5*time.Minute, "/usr/bin/ditto", src, staged); err != nil {
		return fmt.Errorf("couldn't copy %s next to %s: %v %s", AppBundleName, bundle, err, out)
	}
	if err := os.Rename(bundle, old); err != nil {
		return fmt.Errorf("couldn't move the old %s aside (no permission?): %v", bundle, err)
	}
	if err := os.Rename(staged, bundle); err != nil {
		os.Rename(old, bundle) // put the old one back
		return fmt.Errorf("couldn't put the new %s in place of %s: %v", AppBundleName, bundle, err)
	}
	os.RemoveAll(old)
	return nil
}
