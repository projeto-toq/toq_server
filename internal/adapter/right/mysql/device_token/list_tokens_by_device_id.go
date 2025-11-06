package mysqldevicetokenadapter

// ListTokensByDeviceID returns all tokens for a specific device
func (a *DeviceTokenAdapter) ListTokensByDeviceID(userID int64, deviceID string) ([]string, error) {
	query := `SELECT device_token FROM device_tokens WHERE user_id = ? AND device_id = ?`

	rows, err := a.db.GetDB().Query(query, userID, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}
