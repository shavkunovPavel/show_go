package users

import (
	"admigo/common"
	"admigo/model"
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

func getSelect() (res string) {
	res = fmt.Sprintf(model.GetFormat(2),
		"select id, name, email, password, created_at, confirmed, coalesce(nullif(thumb, ''), 'noa.png') thumb",
		"from users",
	)
	return
}

func getUser(stringWhere string, val interface{}, hidePass bool, fillAttributes bool) (u *UserModel, err error) {
	if len(stringWhere) == 0 {
		return
	}
	us := UserModel{}
	sql := fmt.Sprintf(model.GetFormat(2),
		getSelect(),
		stringWhere,
	)
	err = model.Db.QueryRow(sql, val).Scan(&us.Id, &us.Name, &us.Email,
		&us.Password, &us.CreatedAt, &us.Confirmed, &us.Thumb,
	)
	if err != nil {
		return
	}
	if hidePass {
		us.Password = ""
	}
	if fillAttributes {
		us.fillAttr()
	}
	us.fillRole()
	u = &us
	return
}

// Get a single user given the email
func UserByEmail(email string) (u *UserModel, err error) {
	u, err = getUser("where email = $1", email, false, false)
	return
}

// Get a single user given the id
func UserById(id int, fillAttributes bool) (u *UserModel, err error) {
	u, err = getUser("where id = $1", id, true, fillAttributes)
	return
}

func validateEdit(u *UserModel) (errors map[string]string) {
	errors = make(map[string]string)
	model.Required(&errors, map[string][]string{
		"name":  []string{u.Name, "Name"},
		"email": []string{u.Email, "Email"},
	})
	if len(errors) == 0 {
		return nil
	}
	return
}

// Validating login parameters
func validateLogin(f *userRequest) (errors map[string]string) {
	errors = make(map[string]string)
	model.Required(&errors, map[string][]string{
		"email":    []string{f.Email, "Email"},
		"password": []string{f.Password, "Password"},
	})
	if len(errors) == 0 {
		return nil
	}
	return
}

// Validating before register a new user
func validateRegister(f *userRequest) (errors map[string]string) {
	errors = make(map[string]string)
	model.Required(&errors, map[string][]string{
		"name":      []string{f.Name, "Name"},
		"email":     []string{f.Email, "Email"},
		"password":  []string{f.Password, "Password"},
		"cpassword": []string{f.Cpassword, "Confirm Password"},
	})
	model.Confirmed(&errors, map[string][]string{
		"cpassword": []string{f.Password, f.Cpassword, "Confirm Password"},
	})
	if len(errors) == 0 {
		return nil
	}
	return
}

// Create a new user
func UserCreate(request *http.Request) (res *model.Result) {
	var form userRequest
	model.FormToJson(request.Body, &form)
	if uerr := validateRegister(&form); uerr != nil {
		res = model.GetErrorResult(uerr)
		return
	}
	new_user := UserModel{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := new_user.create(); err != nil {
		res = model.GetErrorResult(map[string]string{"all": err.Error()})
		return
	}
	common.SendEmailUserRegistration(form.Email)
	res = model.GetOk(common.Mess().Regfin)
	return
}

// Login user
func UserLogin(request *http.Request) (res *model.Result, sessUuid string) {
	var form userRequest
	model.FormToJson(request.Body, &form)
	if uerr := validateLogin(&form); uerr != nil {
		res = model.GetErrorResult(uerr)
		return
	}
	user, err := UserByEmail(form.Email)
	if err != nil {
		res = model.GetErrorResult(map[string]string{"email": common.Mess().EmailNotFound})
		return
	}
	if user.Confirmed == 0 {
		res = model.GetErrorResult(map[string]string{"all": common.Mess().UserNotConfirmed})
		return
	}
	if user.Password == common.Encrypt(form.Password) {
		session, err := user.createSession()
		if err != nil {
			res = model.GetErrorResult(map[string]string{"all": err.Error()})
			return
		}
		sessUuid = session.Uuid
		res = model.GetOkName(fmt.Sprintf("%s, %s", common.Mess().Welcome, user.Name), user.Role.Name)
		return
	}
	res = model.GetErrorResult(map[string]string{"password": common.Mess().IncorrectPassword})
	return
}

func UserForceLogin(email string) (uuid string, err error) {
	user, err := UserByEmail(email)
	if err != nil {
		err = errors.New(common.Mess().UserNotFound)
		return
	}
	if err = user.setConfirmed(); err != nil {
		return
	}
	sess, err := user.createSession()
	uuid = sess.Uuid
	return
}

// update user
func UserUpdate(data string, thumb string, logged *UserModel) (res *model.Result) {
	var user UserModel
	model.FormToJson(bytes.NewReader([]byte(data)), &user)

	if uerr := validateEdit(&user); uerr != nil {
		res = model.GetErrorResult(uerr)
		return
	}

	if len(thumb) > 0 {
		user.Thumb = thumb
	}

	var err error
	if user.Id > 0 {
		err = user.update(logged)
	} else {
		err = user.create()
	}

	if err != nil {
		res = model.GetErrorResult(map[string]string{"all": err.Error()})
		return
	}

	res = model.GetOk("user updated")
	return
}
