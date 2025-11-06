package mysqldevicetokenadapter

// ListTokensByOptedInUsers returns all tokens for users who opted in to notifications
func (a *DeviceTokenAdapter) ListTokensByOptedInUsers() ([]string, error) {
	query := `SELECT DISTINCT dt.device_token 
			  FROM device_tokens dt 
			  INNER JOIN users u ON dt.user_id = u.id 
			  WHERE u.opt_status = 1`

	rows, err := a.db.GetDB().Query(query)
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
