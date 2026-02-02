package user

import (
	"context"
	"github.com/yannml220/chat_agent_app/internal/app/db"
)

type UserRepo interface {
	GetUserById(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User, id string) error
	DeleteUser(ctx context.Context, id string) error
}

type userRepo struct {
	Db *db.Db
}

func NewRepo(db *db.Db) UserRepo {
	return &userRepo{
		Db: db,
	}

}

func (ur *userRepo) GetUserById(ctx context.Context, id string) (*User, error) {
	sql := "select id , username ,name ,profile_picture_url , email, confirmed , created_at from users where id = $1"
	user := User{}

	if err := ur.Db.GetDb().QueryRow(ctx, sql, id).Scan(&user.Id, &user.Username, &user.Name, &user.ProfilePictureUrl, &user.Email, &user.Confirmed, &user.CreatedAt); err != nil {

		return nil, err

	}

	return &user, nil

}

func (ur *userRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	sql := "select id , username ,name ,profile_picture_url , email, confirmed , created_at from users where email = $1"
	user := User{}

	if err := ur.Db.GetDb().QueryRow(ctx, sql, email).Scan(&user.Id, &user.Username, &user.Name, &user.ProfilePictureUrl, &user.Email, &user.Confirmed, &user.CreatedAt); err != nil {

		return nil, err

	}

	return &user, nil

}

func (ur *userRepo) UpdateUser(ctx context.Context, user *User, id string) error {
	//userUpdate := user.GetUpdatableUserStruct()

	sql := `
        UPDATE users 
        SET 
            username = COALESCE($2::text, username),
            name = COALESCE($3::text, name),
            profile_picture_url = COALESCE($4::text, profile_picture_url),
            email = COALESCE($5::text, email),
            confirmed = COALESCE($6::boolean, confirmed),
            updated_at = CURRENT_TIMESTAMP 
        WHERE id = $1`

	_, err := ur.Db.GetDb().Exec(ctx, sql,
		id,
		user.Username,
		user.Name,
		user.ProfilePictureUrl,
		user.Email,
		user.Confirmed,
	)

	return err
}

func (ur *userRepo) DeleteUser(ctx context.Context, userId string) error {
	sql := "delete from users where id = $1"

	if _, err := ur.Db.GetDb().Exec(ctx, sql, userId); err != nil {

		return err
	}
	return nil
}
