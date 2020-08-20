// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.
package trac

// GetStringConfig retrieves a value from the Trac config as a string.
func (accessor *DefaultAccessor) GetStringConfig(sectionName string, configName string) string {
	configValue, err := accessor.config.Section(sectionName).GetKey(configName)
	if err != nil {
		return ""
	}

	return configValue.String()
}
