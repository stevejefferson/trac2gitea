// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"github.com/go-ini/ini"
	"github.com/stevejefferson/trac2gitea/log"
)

func getStringConfigValue(config *ini.File, sectionName string, configName string) string {
	if config == nil {
		return ""
	}

	configValue, err := config.Section(sectionName).GetKey(configName)
	if err != nil {
		return ""
	}

	return configValue.String()
}

// GetStringConfig retrieves a value from the Gitea config as a string.
func (accessor *DefaultAccessor) GetStringConfig(sectionName string, configName string) string {
	mainConfigValue := getStringConfigValue(accessor.mainConfig, sectionName, configName)
	if mainConfigValue != "" {
		log.Debug("found value in Gitea main config section=%s, name=%s, value=%s", sectionName, configName, mainConfigValue)
		return mainConfigValue
	}

	customConfigValue := getStringConfigValue(accessor.customConfig, sectionName, configName)
	log.Debug("found value in Gitea custom config section=%s, name=%s, value=%s", sectionName, configName, customConfigValue)
	return customConfigValue
}
