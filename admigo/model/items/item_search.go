package items

import (
	"admigo/model"
	"net/http"
)

type ApiItemRequest struct {
	Query *model.DataQuery `json:"query,omitempty"`
	Items []*ItemModel     `json:"items,omitempty"`
}

func List(r *http.Request) (resp *ApiItemRequest, err error) {
	var data ApiItemRequest
	model.FormToJson(r.Body, &data)
	sql := getSelect()
	rows, err := model.GetRows(data.Query, &sql)
	if err != nil {
		return
	}
	for rows.Next() {
		im := ItemModel{}
		if err = rows.Scan(&im.Id, &im.Nm, &im.Description, &im.Price, &im.Thumb, &im.Additional); err != nil {
			return
		}
		data.Items = append(data.Items, &im)
	}
	rows.Close()
	resp = &data
	return
}
