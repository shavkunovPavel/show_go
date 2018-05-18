package items

import (
	"admigo/model"
	"fmt"
)

type ItemModel struct {
	Id          int    `json:"id,omitempty"`
	Nm          string `json:"nm"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Additional  string `json:"additional"`
	Thumb       string `json:"thumb"`
}

func (item_ed *ItemModel) update() (err error) {
	sql := "update items set"
	sql += fmt.Sprint("\nnm = $1,")
	sql += fmt.Sprint("\nprice = $2,")
	sql += fmt.Sprint("\ndescription = $3")

	if len(item_ed.Thumb) > 0 {
		sql += fmt.Sprintf(",\nthumb = '%s'", item_ed.Thumb)
	}

	sql += fmt.Sprint("\nwhere id=$4")

	_, err = model.Db.Exec(sql, item_ed.Nm, item_ed.Price, item_ed.Description, item_ed.Id)
	return
}

// Create a new item
func (item_ed *ItemModel) create() (err error) {
	sql := fmt.Sprintf(model.GetFormat(3),
		"insert into items (nm, description, price, thumb)",
		"values ($1, $2, $3, $4)",
		"returning id",
	)
	stmt, err := model.Db.Prepare(sql)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(item_ed.Nm, item_ed.Description,
		item_ed.Price, item_ed.Thumb,
	).Scan(&item_ed.Id)
	if err != nil {
		return
	}
	return
}

func (item_ed *ItemModel) Delete() (err error) {
	sql := "delete from items where id = $1"
	_, err = model.Db.Exec(sql, item_ed.Id)
	return
}
