package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/yannml220/chat_agent_app/internal/app/api/user"
	"github.com/yannml220/chat_agent_app/internal/app/db"

	"github.com/jackc/pgx/v5"
)

type AuthRepo interface {
	IsTokenReuseDetected(ctx context.Context, hashedCookieRefreshToken string, sessionId string, userId string, deviceId string, userAgent string) (bool, error)
	RotateRefreshTokenAndExtendSession(ctx context.Context, oldHashedToken string, userId string, sessionId string, deviceId string, userAgent string, newHashedToken string, newExpiresAt time.Time) error
	GetSessionById(ctx context.Context, id string) (*Session, error)
	GetSessionByRefreshToken(ctx context.Context, hashedRefreshToken string) (*Session, error)
	GetUserIdentity(ctx context.Context, userId string, provider string) (*Identity, error)
	CreateSession_(ctx context.Context, session *Session, userId string) (string, error)
	CreateIdentity(ctx context.Context, identity *Identity, userId string) (string, error)
	UpdateUserIdentity(ctx context.Context, identity *Identity, userId string, provider string) error
	DeleteSessionById(ctx context.Context, sessionId string) error
	DeleteSessionByRefreshTokenHash(ctx context.Context, hashed_refresh_token string) error
	DeleteAllUserSessions(ctx context.Context, userId string) error
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	CreateUserAndIdentity(ctx context.Context, user *user.User, identity *Identity) (userId string, identityId string, err error)
	VerifyUserEmail(ctx context.Context, userId string) error
}

type authRepo struct {
	Db *db.Db
}

func NewRepo(db *db.Db) AuthRepo {
	return &authRepo{
		Db: db,
	}

}

func (a *authRepo) GetSessionById(ctx context.Context, id string) (*Session, error) {

	sql := "select id , user_id,device_id , user_agent , expires_at ,last_activity_at from sessions  where sessions.id = $1"

	session := Session{}

	if err := a.Db.GetDb().QueryRow(ctx, sql, id).Scan(&session.Id, &session.UserId, &session.DeviceId, &session.UserAgent, &session.ExpiresAt, &session.LastActivityAt); err != nil {

		return nil, err

	}
	return &session, nil

}

func (a *authRepo) GetSessionByRefreshToken(ctx context.Context, hashedRefreshToken string) (*Session, error) {

	sql := "select id , user_id,device_id , user_agent , expires_at ,last_activity_at from sessions  where sessions.hashed_refresh_token = $1"

	session := Session{}

	if err := a.Db.GetDb().QueryRow(ctx, sql, hashedRefreshToken).Scan(&session.Id, &session.UserId, &session.DeviceId, &session.UserAgent, &session.ExpiresAt, &session.LastActivityAt); err != nil {

		return nil, err

	}
	return &session, nil

}

func (a *authRepo) CreateSession_(ctx context.Context, session *Session, userId string) (string, error) {
	sql := `insert into sessions (user_id, device_id, user_agent, hashed_refresh_token,refresh_token_expires_at,expires_at ) values ($1, $2, $3, $4, $5, $6) returning id`

	log.Print("here is the session values : ")
	log.Print("useri id", session.UserId)
	log.Print("device id", session.DeviceId)
	log.Print("user agent", session.UserAgent)
	log.Print("hashed refresh token", session.HashedRefreshToken)
	log.Print("token expiration", session.RefreshTokenExpiresAt)
	log.Print("session expiration", session.ExpiresAt)

	var sessionId string
	err := a.Db.GetDb().QueryRow(
		ctx,
		sql,
		userId,
		session.DeviceId,
		session.UserAgent,
		session.HashedRefreshToken,
		session.RefreshTokenExpiresAt,
		session.ExpiresAt,
	).Scan(&sessionId)

	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (a *authRepo) IsTokenReuseDetected(ctx context.Context, hashedCookieRefreshToken string, sessionId string, userId string, deviceId string, userAgent string) (bool, error) {

	sql := `select user_id,device_id, user_agent, hashed_refresh_token, refresh_token_expires_at from sessions where id = $1 `

	var session Session
	err := a.Db.GetDb().QueryRow(ctx, sql, sessionId).Scan(&session.UserId, &session.DeviceId, &session.UserAgent, &session.HashedRefreshToken, &session.RefreshTokenExpiresAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return true, nil
		}
		return false, err
	}

	if session.UserId != userId ||
		session.DeviceId != deviceId ||
		session.UserAgent != userAgent ||
		session.HashedRefreshToken != hashedCookieRefreshToken {
		return true, nil
	}

	if session.RefreshTokenExpiresAt.Before(time.Now()) {
		return true, nil
	}

	return false, nil
}

func (a *authRepo) RotateRefreshTokenAndExtendSession(ctx context.Context, oldHashedToken string, userId string, sessionId string, deviceId string, userAgent string, newHashedToken string, newExpiresAt time.Time) error {

	log.Print("Session id : ", sessionId)
	log.Print("user id : ", userId)
	sql := `
        update sessions
        set 
            hashed_refresh_token = $1,
            refresh_token_expires_at = $2,
            expires_at = $2,
            last_activity_at = CURRENT_TIMESTAMP,
            user_agent = $3,
            updated_at = CURRENT_TIMESTAMP
        where 
            id = $4
            and user_id = $5
            and device_id = $6
            and hashed_refresh_token = $7
    `

	result, err := a.Db.GetDb().Exec(
		ctx,
		sql,
		newHashedToken,
		newExpiresAt,
		userAgent,
		sessionId,
		userId,
		deviceId,
		oldHashedToken,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("session not found or token mismatch")
	}

	return nil
}

func (a *authRepo) CreateIdentity(ctx context.Context, identity *Identity, userId string) (string, error) {
	sql := "insert into identities (provider_user_id,user_id,provider) values ($1, $2 , $3)  returning id"
	var ret string

	err := a.Db.GetDb().QueryRow(ctx, sql, identity.ProviderUserId, userId, identity.Provider).Scan(&ret)

	if err != nil {
		return "", err
	}

	return ret, nil

}

func (a *authRepo) UpdateUserIdentity(ctx context.Context, identity *Identity, userId string, provider string) error {

	type identityUpdateStruct struct {
		providerUserId *string
		userId         *string
		provider       *string
	}

	identityUpdate := identityUpdateStruct{
		providerUserId: &identity.ProviderUserId,
		userId:         &userId,
		provider:       &identity.Provider,
	}

	sql := "update identities set provider_user_id = coalesce($2,provider_user_id) , provider = coalesce($3,provider)  , last_signin_at = CURRENT_TIMESTAMP where user_id = $1 and provider = $3"

	if _, err := a.Db.GetDb().Exec(ctx, sql, userId, identityUpdate.providerUserId, identityUpdate.provider); err != nil {
		return err
	}

	return nil

}

func (a *authRepo) DeleteSessionById(ctx context.Context, sessionId string) error {
	sql := "delete from sessions where id = $1"

	if _, err := a.Db.GetDb().Exec(ctx, sql, sessionId); err != nil {
		return err

	}
	return nil

}

func (a *authRepo) DeleteSessionByRefreshTokenHash(ctx context.Context, hashedRefreshToken string) error {
	sql := "delete from sessions where sessions.hashed_refresh_token = $1"

	if _, err := a.Db.GetDb().Exec(ctx, sql, hashedRefreshToken); err != nil {
		return err

	}
	return nil

}

func (a *authRepo) DeleteAllUserSessions(ctx context.Context, userId string) error {

	sql := "delete from sessions where user_id = $1"

	if _, err := a.Db.GetDb().Exec(ctx, sql, userId); err != nil {
		return err

	}
	return nil

}

func (a *authRepo) GetUserIdentity(ctx context.Context, userId string, provider string) (*Identity, error) {

	sql := "select id,provider_user_id,user_id,provider from identities where user_id = $1 and provider = $2"
	identity := Identity{}

	if err := a.Db.GetDb().QueryRow(ctx, sql, userId, provider).Scan(&identity.Id, &identity.ProviderUserId, &identity.UserId, &identity.Provider); err != nil {

		if errors.Is(err, pgx.ErrNoRows) || err.Error() == "no rows in result set" {
            return nil, nil
        }

		return nil, err

	}

	return &identity, nil

}

func (a *authRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	sql := "select * from users where email = $1"
	user := user.User{}

	if err := a.Db.GetDb().QueryRow(ctx, sql, email).Scan(&user.Id, &user.Username, &user.Email, &user.Confirmed,  &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {

		if errors.Is(err, pgx.ErrNoRows) || err.Error() == "no rows in result set" {
            return nil, nil
        }

		return nil, err

	}

	return &user, nil

}



func (a *authRepo) CreateUserAndIdentity(ctx context.Context, user *user.User, identity *Identity) (string, string, error) {
	sqlUser := "insert into users (email, confirmed) values ($1, $2) returning id"
	sqlIdentity := "insert into identities (provider_user_id, user_id, provider) values ($1, $2, $3) returning id"

	var userId, identityId string

	tx, err := a.Db.GetDb().Begin(ctx)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	if err := tx.QueryRow(ctx, sqlUser, user.Email, user.Confirmed).Scan(&userId); err != nil {
		return "", "", err
	}

	if err := tx.QueryRow(ctx, sqlIdentity, identity.ProviderUserId, userId, identity.Provider).Scan(&identityId); err != nil {
		return "", "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", "", err
	}

	return userId, identityId, nil
}

func (a *authRepo) VerifyUserEmail(ctx context.Context, userId string) error {

	sql := "update users set confirmed = true where id = $1"

	if _, err := a.Db.GetDb().Exec(ctx, sql, userId); err != nil {
		return err

	}
	return nil

}
