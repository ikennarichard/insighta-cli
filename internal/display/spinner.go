// internal/display/spinner.go
package display

import (
    "sync"

    "github.com/briandowns/spinner"
)

var (
    currentSpinner *spinner.Spinner
    mu             sync.Mutex
)

func ShowLoader(message string) {
    mu.Lock()
    defer mu.Unlock()

    if currentSpinner != nil {
        currentSpinner.Stop()
    }

    currentSpinner = spinner.New(spinner.CharSets[14], 100)
    currentSpinner.Suffix = " " + message
    currentSpinner.Start()
}

func StopLoader() {
    mu.Lock()
    defer mu.Unlock()

    if currentSpinner != nil {
        currentSpinner.Stop()
        currentSpinner = nil
    }
}