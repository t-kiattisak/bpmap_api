package worker

import (
	"pbmap_api/src/pkg/config"
)

// StartBackgroundJobs starts background jobs. Returns a cleanup function.
// Adapter/data_sync were removed; add new jobs here when needed.
func StartBackgroundJobs(cfg *config.Config) func() {
	_ = cfg
	return func() {}
}
