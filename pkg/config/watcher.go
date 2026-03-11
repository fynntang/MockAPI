package config

import (
	"os"
	"time"
)

// Watcher watches the config file for changes and reloads
type Watcher struct {
	configFile string
	onChange   func(*Config)
	stopCh     chan struct{}
}

// NewWatcher creates a new config file watcher
func NewWatcher(configFile string, onChange func(*Config)) *Watcher {
	return &Watcher{
		configFile: configFile,
		onChange:   onChange,
		stopCh:     make(chan struct{}),
	}
}

// Start begins watching for config changes using polling
func (w *Watcher) Start() {
	go func() {
		var lastMod time.Time
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-w.stopCh:
				return
			case <-ticker.C:
				info, err := os.Stat(w.configFile)
				if err != nil {
					continue
				}
				if info.ModTime().After(lastMod) {
					if !lastMod.IsZero() {
						// File changed, reload
						newCfg, err := Load(w.configFile)
						if err == nil && w.onChange != nil {
							w.onChange(newCfg)
						}
					}
					lastMod = info.ModTime()
				}
			}
		}
	}()
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	close(w.stopCh)
}