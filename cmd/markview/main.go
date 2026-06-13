package main

import (
	"flag"
	"fmt"
	"os"

	"fyne.io/fyne/v2/app"
	"github.com/sbj-ee/markview/internal/gui"
	"github.com/sbj-ee/markview/internal/themes"
	"go.uber.org/zap"
)

var (
	version = "1.0.21"
	logger  *zap.Logger
)

func main() {
	// Parse command-line flags
	filePath := flag.String("file", "", "Path to markdown file to open")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Allow the file to be given as a positional argument too, e.g.
	// "markview notes.md" in addition to "markview --file notes.md".
	if *filePath == "" && flag.NArg() > 0 {
		*filePath = flag.Arg(0)
	}

	// Initialize logger
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	if *showVersion {
		fmt.Printf("markview version %s\n", version)
		return
	}

	// Create Fyne application
	myApp := app.NewWithID("com.sbj-ee.markview")
	myApp.SetIcon(themes.AppLogo())

	// Create main window
	window := gui.NewWindow(myApp, logger, version)

	// Load file if provided via command line
	if *filePath != "" {
		window.LoadFile(*filePath)
	}

	window.ShowAndRun()
}
