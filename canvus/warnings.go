package canvus

import (
	"log"
	"os"
	"sync"
)

// APIWarning represents a known API limitation or issue.
type APIWarning struct {
	Code        string // Short identifier (e.g., "NOTE_TITLE_NOT_EXPOSED")
	Description string // Human-readable description
	Workaround  string // Suggested workaround, if any
	IssueURL    string // Link to GitLab/GitHub issue tracking the fix
}

// Known API limitations and issues.
// These are documented in CANVUS_API_ISSUES_REPORT.md
var (
	// WarningNoteTitleNotExposed indicates that Note widget titles cannot be read or updated via API.
	WarningNoteTitleNotExposed = APIWarning{
		Code:        "NOTE_TITLE_NOT_EXPOSED",
		Description: "Note widget 'title' field is not exposed by the Canvus API. Title values in requests are ignored and responses will not include the title.",
		Workaround:  "Use the 'name' field instead for identifying notes.",
		IssueURL:    "https://gitlab.multitaction.com/swrd/conan/canvus/canvus-app/-/issues/38",
	}

	// WarningVideoInputTitleNotExposed indicates that VideoInput widget titles cannot be read or updated via API.
	WarningVideoInputTitleNotExposed = APIWarning{
		Code:        "VIDEOINPUT_TITLE_NOT_EXPOSED",
		Description: "VideoInput widget 'title' field is not exposed by the Canvus API. Title values in requests are ignored and responses will not include the title.",
		Workaround:  "No workaround available. Await Canvus API fix.",
		IssueURL:    "https://gitlab.multitaction.com/swrd/conan/canvus/canvus-app/-/issues/13",
	}

	// WarningPDFSizeBug indicates that PDF widget size changes may not work as expected.
	WarningPDFSizeBug = APIWarning{
		Code:        "PDF_SIZE_BUG",
		Description: "PDF widget size changes via PATCH update the bounding box but the actual PDF content stays at its original size, creating a visual disconnect.",
		Workaround:  "Avoid resizing PDFs via API. Delete and recreate if different size is needed.",
		IssueURL:    "https://gitlab.multitaction.com/swrd/conan/canvus/canvus-app/-/issues/15",
	}

	// WarningImageAspectRatioNotPreserved indicates that Image widget resizing does not preserve aspect ratio.
	WarningImageAspectRatioNotPreserved = APIWarning{
		Code:        "IMAGE_ASPECT_RATIO_NOT_PRESERVED",
		Description: "Image widget size changes via PATCH do not preserve aspect ratio. Content will stretch/distort to match exact requested dimensions.",
		Workaround:  "Calculate correct aspect-ratio-preserving dimensions before calling UpdateImage.",
		IssueURL:    "https://gitlab.multitaction.com/swrd/conan/canvus/canvus-app/-/issues/39",
	}

	// WarningVideoAspectRatioNotPreserved indicates that Video widget resizing does not preserve aspect ratio.
	WarningVideoAspectRatioNotPreserved = APIWarning{
		Code:        "VIDEO_ASPECT_RATIO_NOT_PRESERVED",
		Description: "Video widget size changes via PATCH do not preserve aspect ratio. Content will stretch/distort to match exact requested dimensions.",
		Workaround:  "Calculate correct aspect-ratio-preserving dimensions before calling UpdateVideo.",
		IssueURL:    "https://gitlab.multitaction.com/swrd/conan/canvus/canvus-app/-/issues/39",
	}
)

// warningLogger handles API limitation warnings.
// Can be disabled via DisableAPIWarnings() or by setting CANVUS_SDK_DISABLE_WARNINGS=1.
var (
	warningLogger  = log.New(os.Stderr, "[canvus-sdk WARNING] ", log.LstdFlags)
	warningsOnce   sync.Once
	warningsIssued = make(map[string]bool)
	warningsMu     sync.Mutex
	warningsEnabled = true
)

func init() {
	// Check environment variable to disable warnings
	if os.Getenv("CANVUS_SDK_DISABLE_WARNINGS") == "1" {
		warningsEnabled = false
	}
}

// DisableAPIWarnings disables all API limitation warnings from the SDK.
// This is useful for production environments where warnings have been acknowledged.
func DisableAPIWarnings() {
	warningsMu.Lock()
	defer warningsMu.Unlock()
	warningsEnabled = false
}

// EnableAPIWarnings enables API limitation warnings from the SDK.
func EnableAPIWarnings() {
	warningsMu.Lock()
	defer warningsMu.Unlock()
	warningsEnabled = true
}

// SetWarningLogger sets a custom logger for API warnings.
// Pass nil to use the default stderr logger.
func SetWarningLogger(logger *log.Logger) {
	warningsMu.Lock()
	defer warningsMu.Unlock()
	if logger == nil {
		warningLogger = log.New(os.Stderr, "[canvus-sdk WARNING] ", log.LstdFlags)
	} else {
		warningLogger = logger
	}
}

// warnOnce emits a warning only once per warning code during the lifetime of the process.
// This prevents excessive logging when the same problematic operation is called repeatedly.
func warnOnce(warning APIWarning) {
	warningsMu.Lock()
	defer warningsMu.Unlock()

	if !warningsEnabled {
		return
	}

	if warningsIssued[warning.Code] {
		return
	}
	warningsIssued[warning.Code] = true

	warningLogger.Printf("%s: %s", warning.Code, warning.Description)
	if warning.Workaround != "" {
		warningLogger.Printf("  Workaround: %s", warning.Workaround)
	}
	if warning.IssueURL != "" {
		warningLogger.Printf("  Issue: %s", warning.IssueURL)
	}
}

// warnAlways emits a warning every time it's called.
// Use for operations where the user should be reminded each time.
func warnAlways(warning APIWarning) {
	warningsMu.Lock()
	defer warningsMu.Unlock()

	if !warningsEnabled {
		return
	}

	warningLogger.Printf("%s: %s", warning.Code, warning.Description)
}

// ResetWarnings clears the record of issued warnings, allowing them to be shown again.
// This is primarily useful for testing.
func ResetWarnings() {
	warningsMu.Lock()
	defer warningsMu.Unlock()
	warningsIssued = make(map[string]bool)
}
