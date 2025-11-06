package mysqldevicetokenadapter

// ListTokensByUserIDIfOptedIn returns tokens for a specific user if they opted in
func (a *DeviceTokenAdapter) ListTokensByUserIDIfOptedIn(userID int64) ([]string, error) {
	query := `SELECT DISTINCT dt.device_token 
			  FROM device_tokens dt 
			  INNER JOIN users u ON dt.user_id = u.id 
			  WHERE u.id = ? AND u.opt_status = 1`

	rows, err := a.db.GetDB().Query(query, userID)
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
