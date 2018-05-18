package users

import (
	"admigo/common"
	"admigo/model"
	"admigo/model/roles"
	"errors"
	"fmt"
	"time"
)

type userRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Cpassword string `json:"cpassword"`
}

type UserModel struct {
	Id        int              `json:"id,omitempty"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	Password  string           `json:"password,omitempty"`
	CreatedAt *time.Time       `json:"created_at,omitempty"`
	Confirmed int              `json:"confirmed"`
	Thumb     string           `json:"thumb"`
	Phones    []UserAttr       `json:"phones,omitempty"`
	Emails    []UserAttr       `json:"emails,omitempty"`
	Role      *roles.RoleModel `json:"role,omitempty"`
}

type UserAttr struct {
	Id  int    `json:"id"`
	Val string `json:"val"`
}

const (
	img_path = "/static/images/users/"
)

func (user *UserModel) GetThumb() (pth string) {
	pth = img_path + user.Thumb
	return
}

func (user *UserModel) CanWrite() bool {
	if user.Role.Id == roles.ADMIN || user.Role.Id == roles.WRITER {
		return true
	}
	return false
}

func (user *UserModel) IsAdmin() bool {
	if user.Role.Id == roles.ADMIN {
		return true
	}
	return false
}

func (user *UserModel) insertAttributes(ar *[]UserAttr, tp int) (err error) {
	for _, at := range *ar {
		if len(at.Val) == 0 {
			continue
		}
		sql := "insert into users_attr(user_id, attr_id, val) values($1, $2, $3)"
		if _, err = model.Db.Exec(sql, user.Id, tp, at.Val); err != nil {
			return
		}
	}
	return
}

func (user *UserModel) updateAttributes() (err error) {
	if err = user.deleteAttributes(); err != nil {
		return
	}
	if err = user.insertAttributes(&user.Phones, 1); err != nil {
		return
	}
	if err = user.insertAttributes(&user.Emails, 2); err != nil {
		return
	}
	return
}

// Create a new user, save user info into the database
func (user *UserModel) create() (err error) {
	statement := fmt.Sprintf(model.GetFormat(3),
		"insert into users (name, email, password, created_at, confirmed, thumb)",
		"values ($1, $2, $3, $4, $5, $6)",
		"returning id, created_at",
	)
	stmt, err := model.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(user.Name,
		user.Email, common.Encrypt(user.Password),
		time.Now(), user.Confirmed, user.Thumb,
	).Scan(&user.Id, &user.CreatedAt)
	if err != nil {
		return
	}

	err = user.updateAttributes()
	if err != nil {
		return
	}

	err = user.updateRole()
	return
}

func (user *UserModel) update(logged *UserModel) (err error) {
	sql := "update users set name = $1, email = $2, thumb = $3%s%s where id = %d"
	var setpas string
	var conf string

	if len(user.Password) > 0 {
		user.Password = common.Encrypt(user.Password)
		setpas = fmt.Sprintf(", password = '%s'", user.Password)
	}

	if logged.IsAdmin() {
		conf = fmt.Sprintf(", confirmed = %d", user.Confirmed)
	}

	sql = fmt.Sprintf(sql, conf, setpas, user.Id)

	_, err = model.Db.Exec(sql, user.Name, user.Email, user.Thumb)
	if err != nil {
		return
	}

	err = user.updateAttributes()
	if err != nil {
		return
	}

	if logged.IsAdmin() {
		err = user.updateRole()
	}

	return
}

func (user *UserModel) setConfirmed() (err error) {
	sql := "update users set confirmed = 1 where id = $1"
	_, err = model.Db.Exec(sql, user.Id)
	return
}

// Get user's session
func (user *UserModel) userSession() (session *SessionModel, err error) {
	se := SessionModel{}
	sql := fmt.Sprintf(model.GetFormat(3),
		"select id, uuid, user_id, created_at",
		"from sessions",
		"where user_id = $1",
	)
	err = model.Db.QueryRow(sql, user.Id).Scan(&se.Id, &se.Uuid, &se.UserId, &se.CreatedAt)
	session = &se
	return
}

func (user *UserModel) fill(ar *[]UserAttr, tp int) (err error) {
	sql := fmt.Sprintf(model.GetFormat(4),
		"select id, val",
		"from users_attr",
		"where user_id = $1 and attr_id = $2",
		"order by id",
	)
	rows, err := model.Db.Query(sql, user.Id, tp)
	if err != nil {
		return
	}
	for rows.Next() {
		at := UserAttr{}
		if err = rows.Scan(&at.Id, &at.Val); err != nil {
			return
		}
		*ar = append(*ar, at)
	}
	rows.Close()
	return
}

// Get user's phones
func (user *UserModel) fillAttr() (err error) {
	user.fill(&user.Phones, 1)
	user.fill(&user.Emails, 2)
	return
}

func (user *UserModel) fillRole() (err error) {
	sql := fmt.Sprintf(model.GetFormat(4),
		"select r.id, r.nm",
		"from users_roles u",
		"join roles r on r.id = u.role_id",
		"where u.user_id = $1",
	)
	rows, err := model.Db.Query(sql, user.Id)
	if err != nil {
		return
	}
	if rows.Next() {
		role := roles.RoleModel{}
		if err = rows.Scan(&role.Id, &role.Name); err != nil {
			return
		}
		user.Role = &role
	}
	rows.Close()
	return
}

// Create a new session for an existing user
func (user *UserModel) createSession() (session *SessionModel, err error) {
	if session, err = user.userSession(); err == nil {
		return
	}
	statement := fmt.Sprintf(model.GetFormat(3),
		"insert into sessions (uuid, user_id, created_at)",
		"values ($1, $2, $3)",
		"returning id, uuid, user_id, created_at",
	)
	stmt, err := model.Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	se := SessionModel{}
	err = stmt.QueryRow(common.CreateUUID(),
		user.Id, time.Now(),
	).Scan(&se.Id, &se.Uuid, &se.UserId, &se.CreatedAt)
	session = &se
	return
}

func (user *UserModel) DeleteSessions() (err error) {
	sql := "delete from sessions where user_id = $1"
	_, err = model.Db.Exec(sql, user.Id)
	return
}

func (user *UserModel) deleteAttributes() (err error) {
	sql := "delete from users_attr where user_id = $1"
	_, err = model.Db.Exec(sql, user.Id)
	return
}

func (user *UserModel) insertRole(new_role int) (err error) {
	sql := fmt.Sprintf(model.GetFormat(7),
		"insert into users_roles(user_id, role_id)",
		"select $1, $2",
		"where not exists(",
		"	select 1",
		"	from users_roles uu",
		"	where uu.user_id = $1 and uu.role_id = $2",
		")",
	)
	_, err = model.Db.Exec(sql, user.Id, new_role)
	return
}

func (user *UserModel) updateRole() error {
	if user.Role == nil {
		s_err := user.insertRole(roles.READER)
		return s_err
	}

	sql := fmt.Sprintf(model.GetFormat(3),
		"update users_roles set",
		"role_id = $1",
		"where user_id = $2",
	)
	res, err := model.Db.Exec(sql, user.Role.Id, user.Id)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if ra > 0 || err != nil {
		return err
	}

	err = user.insertRole(user.Role.Id)
	return err
}

func (user *UserModel) deleteRole() (err error) {
	sql := "delete from users_roles where user_id = $1"
	_, err = model.Db.Exec(sql, user.Id)
	return
}

func (user *UserModel) Delete(logged *UserModel) (err error) {
	if !(logged.IsAdmin() || logged.Id == user.Id) {
		err = errors.New(common.Mess().InsufRights)
		return
	}

	if err = user.DeleteSessions(); err != nil {
		return
	}
	if err = user.deleteAttributes(); err != nil {
		return
	}
	if err = user.deleteRole(); err != nil {
		return
	}
	sql := "delete from users where id = $1"
	_, err = model.Db.Exec(sql, user.Id)
	return
}
