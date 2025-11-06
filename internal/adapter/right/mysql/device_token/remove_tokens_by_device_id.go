package mysqldevicetokenadapter

// RemoveTokensByDeviceID deletes all tokens for a specific device
func (a *DeviceTokenAdapter) RemoveTokensByDeviceID(userID int64, deviceID string) error {
	query := `DELETE FROM device_tokens WHERE user_id = ? AND device_id = ?`
	_, err := a.db.GetDB().Exec(query, userID, deviceID)
	return err
}
