package items

import (
	"admigo/model"
	"bytes"
	"fmt"
	"strconv"
)

func getSelect() (res string) {
	res = fmt.Sprintf(model.GetFormat(5),
		"select id, nm,",
		"coalesce(description, '') description,",
		"price, coalesce(thumb, '') thumb,",
		"(cast(price as numeric(256)) / 1000000000000000000) additional",
		"from items",
	)
	return
}

func ItemById(id int) (item *ItemModel, err error) {
	itm := ItemModel{}
	sql := fmt.Sprintf(model.GetFormat(2),
		getSelect(),
		"where id = $1",
	)
	err = model.Db.QueryRow(sql, id).Scan(&itm.Id, &itm.Nm, &itm.Description,
		&itm.Price, &itm.Thumb, &itm.Additional,
	)
	item = &itm
	return
}

func validateEdit(item *ItemModel) (errors map[string]string) {
	ii, _ := strconv.ParseFloat(item.Price, 256)
	prc := ""
	if ii != 0 {
		prc = item.Price
	}

	errors = make(map[string]string)
	model.Required(&errors, map[string][]string{
		"nm":    []string{item.Nm, "Name"},
		"price": []string{prc, "Price"},
	})
	if len(errors) == 0 {
		return nil
	}
	return
}

// update or create an item
func ItemUpdate(data string, thumb string) (res *model.Result) {
	var item_ed ItemModel
	model.FormToJson(bytes.NewReader([]byte(data)), &item_ed)

	if uerr := validateEdit(&item_ed); uerr != nil {
		res = model.GetErrorResult(uerr)
		return
	}

	if len(thumb) > 0 {
		item_ed.Thumb = thumb
	}

	var err error
	if item_ed.Id > 0 {
		err = item_ed.update()
	} else {
		err = item_ed.create()
	}

	if err != nil {
		res = model.GetErrorResult(map[string]string{"all": err.Error()})
		return
	}

	res = model.GetOk("Item updated")
	return
}
