package gitea

import "github.com/go-ini/ini"

func getStringConfigValue(config *ini.File, sectionName string, configName string) string {
	configValue, err := config.Section(sectionName).GetKey(configName)
	if err != nil {
		return ""
	}

	return configValue.String()
}

func (accessor *Accessor) GetStringConfig(sectionName string, configName string) string {
	mainConfigValue := getStringConfigValue(accessor.mainConfig, sectionName, configName)
	if mainConfigValue != "" {
		return mainConfigValue
	}

	return getStringConfigValue(accessor.customConfig, sectionName, configName)
}
