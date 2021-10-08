package data

func QuerySystemConfig(configKey string) (configValue string, err error) {

	row := Db.QueryRow("SELECT config_value FROM system_config WHERE config_key = $1", configKey)
	err = row.Scan(&configValue)
	return
}
