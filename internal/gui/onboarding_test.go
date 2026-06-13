package gui

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"go.uber.org/zap"
)

// TestShouldShowOnboardingOnlyOnce verifies onboarding is offered on the first
// launch and suppressed thereafter.
func TestShouldShowOnboardingOnlyOnce(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()

	w := NewWindow(app, zap.NewNop(), "test")

	if !w.shouldShowOnboarding() {
		t.Fatal("expected onboarding on first launch")
	}
	if w.shouldShowOnboarding() {
		t.Error("onboarding should not be shown again after the first launch")
	}
}
