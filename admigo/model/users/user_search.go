package users

import (
	"admigo/model"
	"net/http"
)

type ApiUserRequest struct {
	Query *model.DataQuery `json:"query,omitempty"`
	Users []*UserModel     `json:"users,omitempty"`
}

func List(r *http.Request) (resp *ApiUserRequest, err error) {
	var data ApiUserRequest
	model.FormToJson(r.Body, &data)
	sql := getSelect()
	rows, err := model.GetRows(data.Query, &sql)
	if err != nil {
		return
	}
	for rows.Next() {
		us := UserModel{}
		if err = rows.Scan(&us.Id, &us.Name, &us.Email, &us.Password,
			&us.CreatedAt, &us.Confirmed, &us.Thumb,
		); err != nil {
			return
		}
		us.fillRole()
		data.Users = append(data.Users, &us)
	}
	rows.Close()
	resp = &data
	return
}
