package trac

import "path/filepath"

// GetFullPath retrieves the absolute path of a path relative to the root of the Trac installation.
func (accessor *DefaultAccessor) GetFullPath(element ...string) string {
	return filepath.Join(accessor.rootDir, filepath.Join(element...))
}
