package mysqldevicetokenadapter

// RemoveToken deletes a specific device token
func (a *DeviceTokenAdapter) RemoveToken(userID int64, token string) error {
	query := `DELETE FROM device_tokens WHERE user_id = ? AND device_token = ?`
	_, err := a.db.GetDB().Exec(query, userID, token)
	return err
}
