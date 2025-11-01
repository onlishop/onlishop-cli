package system

import (
	"os"
	"path"
)

// GetOnlishopCliCacheDir returns the base cache directory for onlishop-c
func GetOnlishopCliCacheDir() string {
	if dir := os.Getenv("ONLISHOP_CLI_CACHE_DIR"); dir != "" {
		return dir
	}

	cacheDir, _ := os.UserCacheDir()

	return path.Join(cacheDir, "onlishop-cli")
}
