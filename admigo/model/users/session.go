package users

import (
	"admigo/model"
	"time"
)

type SessionModel struct {
	Id        int
	Uuid      string
	UserId    int
	CreatedAt time.Time
}

func sessionByUuid(uuid string) *SessionModel {
	s := SessionModel{}
	err := model.Db.QueryRow("select id, uuid, user_id, created_at from sessions where uuid = $1", uuid).
		Scan(&s.Id, &s.Uuid, &s.UserId, &s.CreatedAt)
	if err != nil {
		return nil
	}
	return &s
}

func (s *SessionModel) user() *UserModel {
	u, err := getUser("where id = $1 and confirmed = 1", s.UserId, true, false)
	if err != nil {
		return nil
	}
	return u
}

func SessionUser(uuid string) *UserModel {
	session := sessionByUuid(uuid)
	if session == nil {
		return nil
	}
	return session.user()
}
