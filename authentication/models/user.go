package models

import (
	"assetra/pb"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID  `db:"id"`
	Username            string     `db:"username"`
	Email               string     `db:"email"`
	Password            string     `db:"password"`
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at"`
	Roles               []string   `db:"roles"`
	FailedLoginAttempts int        `db:"failed_login_attempts"`
	LockedUntil         *time.Time `db:"locked_until"`
}

func (u *User) ToProtoBuffer() *pb.User {
	return &pb.User{
		Id:       u.ID.String(),
		Name:     u.Username,
		Email:    u.Email,
		Role:     u.Roles,
		Created:  u.CreatedAt.Unix(),
		Updated:  u.CreatedAt.Unix(),
	}
}

func (u *User) FromProtoBuffer(user *pb.User) {
	u.Username = user.GetName()
	u.Email = user.GetEmail()
	u.Password = user.GetPassword()
	u.CreatedAt = time.Unix(user.Created, 0)
	u.UpdatedAt = time.Unix(user.Updated, 0)
}
