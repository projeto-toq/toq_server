package mysqldevicetokenadapter

// RemoveAllByUserID deletes all device tokens for a user
func (a *DeviceTokenAdapter) RemoveAllByUserID(userID int64) error {
	query := `DELETE FROM device_tokens WHERE user_id = ?`
	_, err := a.db.GetDB().Exec(query, userID)
	return err
}
