package roles

import (
	"admigo/model"
	"fmt"
	"net/http"
)

type ApiRoleRequest struct {
	Query *model.DataQuery `json:"query,omitempty"`
	Roles []*RoleModel     `json:"items,omitempty"`
}

func List(r *http.Request) (resp *ApiRoleRequest, err error) {
	var data ApiRoleRequest
	if r != nil {
		model.FormToJson(r.Body, &data)
	}
	sql := fmt.Sprintf(model.GetFormat(2),
		"select id, nm",
		"from roles",
	)
	rows, err := model.GetRows(data.Query, &sql)
	if err != nil {
		return
	}
	for rows.Next() {
		im := RoleModel{}
		if err = rows.Scan(&im.Id, &im.Name); err != nil {
			return
		}
		data.Roles = append(data.Roles, &im)
	}
	rows.Close()
	resp = &data
	return
}
