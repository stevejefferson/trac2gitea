// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "path/filepath"

// GetFullPath retrieves the absolute path of a path relative to the root of the Trac installation.
func (accessor *DefaultAccessor) GetFullPath(element ...string) string {
	return filepath.Join(accessor.rootDir, filepath.Join(element...))
}
