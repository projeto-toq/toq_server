package sessionmysqladapter

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sessionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/session_model"
	sessionrepository "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/session_repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mysqlSessionAdapter struct {
	db *sql.DB
}

func NewMySQLSessionAdapter(db *sql.DB) sessionrepository.SessionRepoPortInterface {
	return &mysqlSessionAdapter{db: db}
}

func (a *mysqlSessionAdapter) CreateSession(ctx context.Context, session sessionmodel.SessionInterface) (sessionmodel.SessionInterface, error) {
	query := `INSERT INTO sessions (user_id, refresh_hash, token_jti, expires_at, absolute_expires_at, created_at, rotated_at, user_agent, ip, device_id, rotation_counter, last_refresh_at, revoked)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := a.db.ExecContext(ctx, query,
		session.GetUserID(),
		session.GetRefreshHash(),
		session.GetTokenJTI(),
		session.GetExpiresAt(),
		session.GetAbsoluteExpiresAt(),
		session.GetCreatedAt(),
		session.GetRotatedAt(),
		session.GetUserAgent(),
		session.GetIP(),
		session.GetDeviceID(),
		session.GetRotationCounter(),
		session.GetLastRefreshAt(),
		session.GetRevoked(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating session")
	}
	id, _ := res.LastInsertId()
	session.SetID(id)
	return session, nil
}

func (a *mysqlSessionAdapter) GetSessionByID(ctx context.Context, id int64) (sessionmodel.SessionInterface, error) {
	query := `SELECT id, user_id, refresh_hash, token_jti, expires_at, absolute_expires_at, created_at, rotated_at, user_agent, ip, device_id, rotation_counter, last_refresh_at, revoked FROM sessions WHERE id = ?`
	row := a.db.QueryRowContext(ctx, query, id)
	return scanSession(row)
}

func (a *mysqlSessionAdapter) GetActiveSessionByRefreshHash(ctx context.Context, hash string) (sessionmodel.SessionInterface, error) {
	query := `SELECT id, user_id, refresh_hash, token_jti, expires_at, absolute_expires_at, created_at, rotated_at, user_agent, ip, device_id, rotation_counter, last_refresh_at, revoked FROM sessions
		WHERE refresh_hash = ? AND revoked = false AND expires_at > UTC_TIMESTAMP()`
	row := a.db.QueryRowContext(ctx, query, hash)
	s, err := scanSession(row)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (a *mysqlSessionAdapter) RevokeSession(ctx context.Context, id int64) error {
	query := `UPDATE sessions SET revoked = true WHERE id = ?`
	_, err := a.db.ExecContext(ctx, query, id)
	if err != nil {
		return status.Errorf(codes.Internal, "error revoking session")
	}
	return nil
}

func (a *mysqlSessionAdapter) MarkSessionRotated(ctx context.Context, id int64) error {
	query := `UPDATE sessions SET rotated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := a.db.ExecContext(ctx, query, id)
	if err != nil {
		return status.Errorf(codes.Internal, "error rotating session")
	}
	return nil
}

func (a *mysqlSessionAdapter) RevokeSessionsByUserID(ctx context.Context, userID int64) error {
	query := `UPDATE sessions SET revoked = true WHERE user_id = ? AND revoked = false`
	_, err := a.db.ExecContext(ctx, query, userID)
	if err != nil {
		return status.Errorf(codes.Internal, "error revoking user sessions")
	}
	return nil
}

func (a *mysqlSessionAdapter) UpdateSessionRotation(ctx context.Context, id int64, rotationCounter int, lastRefreshAt time.Time) error {
	query := `UPDATE sessions SET rotation_counter = ?, last_refresh_at = ?, rotated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := a.db.ExecContext(ctx, query, rotationCounter, lastRefreshAt, id)
	if err != nil {
		return status.Errorf(codes.Internal, "error updating session rotation meta")
	}
	return nil
}

func (a *mysqlSessionAdapter) GetActiveSessionsByUserID(ctx context.Context, userID int64) ([]sessionmodel.SessionInterface, error) {
	query := `SELECT id, user_id, refresh_hash, token_jti, expires_at, absolute_expires_at, created_at, rotated_at, user_agent, ip, device_id, rotation_counter, last_refresh_at, revoked FROM sessions WHERE user_id = ? AND revoked = false AND expires_at > UTC_TIMESTAMP()`
	rows, err := a.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying sessions")
	}
	defer rows.Close()
	sessions := []sessionmodel.SessionInterface{}
	for rows.Next() {
		var (
			id            int64
			uid           int64
			refreshHash   string
			tokenJTI      sql.NullString
			expiresAt     time.Time
			absoluteExp   sql.NullTime
			createdAt     time.Time
			rotatedAt     sql.NullTime
			ua            sql.NullString
			ip            sql.NullString
			deviceID      sql.NullString
			rotationCtr   sql.NullInt64
			lastRefreshAt sql.NullTime
			revoked       bool
		)
		if err := rows.Scan(&id, &uid, &refreshHash, &tokenJTI, &expiresAt, &absoluteExp, &createdAt, &rotatedAt, &ua, &ip, &deviceID, &rotationCtr, &lastRefreshAt, &revoked); err != nil {
			return nil, status.Errorf(codes.Internal, "error scanning session row")
		}
		s := sessionmodel.NewSession()
		s.SetID(id)
		s.SetUserID(uid)
		s.SetRefreshHash(refreshHash)
		s.SetExpiresAt(expiresAt)
		s.SetCreatedAt(createdAt)
		if tokenJTI.Valid {
			s.SetTokenJTI(tokenJTI.String)
		}
		if absoluteExp.Valid {
			s.SetAbsoluteExpiresAt(absoluteExp.Time)
		}
		if rotatedAt.Valid {
			rt := rotatedAt.Time
			s.SetRotatedAt(&rt)
		}
		if ua.Valid {
			s.SetUserAgent(ua.String)
		}
		if ip.Valid {
			s.SetIP(ip.String)
		}
		if deviceID.Valid {
			s.SetDeviceID(deviceID.String)
		}
		if rotationCtr.Valid {
			s.SetRotationCounter(int(rotationCtr.Int64))
		}
		if lastRefreshAt.Valid {
			lr := lastRefreshAt.Time
			s.SetLastRefreshAt(&lr)
		}
		s.SetRevoked(revoked)
		sessions = append(sessions, s)
	}
	return sessions, nil
}

// TODO: replace with extended scan including new columns once higher-level logic updated
func scanSession(row *sql.Row) (sessionmodel.SessionInterface, error) {
	s := sessionmodel.NewSession()
	var (
		id            int64
		userID        int64
		refreshHash   string
		tokenJTI      sql.NullString
		expiresAt     time.Time
		absoluteExp   sql.NullTime
		createdAt     time.Time
		rotatedAt     sql.NullTime
		ua            sql.NullString
		ip            sql.NullString
		deviceID      sql.NullString
		rotationCtr   sql.NullInt64
		lastRefreshAt sql.NullTime
		revoked       bool
	)
	if err := row.Scan(&id, &userID, &refreshHash, &tokenJTI, &expiresAt, &absoluteExp, &createdAt, &rotatedAt, &ua, &ip, &deviceID, &rotationCtr, &lastRefreshAt, &revoked); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "session not found")
		}
		return nil, status.Errorf(codes.Internal, "error scanning session")
	}
	s.SetID(id)
	s.SetUserID(userID)
	s.SetRefreshHash(refreshHash)
	s.SetExpiresAt(expiresAt)
	s.SetCreatedAt(createdAt)
	if tokenJTI.Valid {
		s.SetTokenJTI(tokenJTI.String)
	}
	if absoluteExp.Valid {
		s.SetAbsoluteExpiresAt(absoluteExp.Time)
	}
	if rotatedAt.Valid {
		rt := rotatedAt.Time
		s.SetRotatedAt(&rt)
	}
	if ua.Valid {
		s.SetUserAgent(ua.String)
	}
	if ip.Valid {
		s.SetIP(ip.String)
	}
	if deviceID.Valid {
		s.SetDeviceID(deviceID.String)
	}
	if rotationCtr.Valid {
		s.SetRotationCounter(int(rotationCtr.Int64))
	}
	if lastRefreshAt.Valid {
		lr := lastRefreshAt.Time
		s.SetLastRefreshAt(&lr)
	}
	s.SetRevoked(revoked)
	return s, nil
}
