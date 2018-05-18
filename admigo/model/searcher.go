package model

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
)

type DataQuery struct {
	Page    int            `json:"page,omitempty"`
	PerPg   int            `json:"perpage,omitempty"`
	CntPage int            `json:"cntpage,omitempty"`
	OrderBy []string       `json:"order_by,omitempty"`
	Filter  []*FilterQuery `json:"filter,omitempty"`
}

type FilterQuery struct {
	Field string `json:"field"`
	Val   string `json:"val"`
	Label string `json:"label"`
}

func GetRows(query *DataQuery, sql_in *string) (*sql.Rows, error) {
	sql := fmt.Sprintf(GetFormat(2), *sql_in, getWhereSQL(query))
	pages := getPageSQl(query, sql)
	sql = fmt.Sprintf(GetFormat(3), sql, getOrderBy(query), pages)
	return Db.Query(sql)
}

func getWhereSQL(query *DataQuery) (wh string) {
	if query == nil {
		return
	}
	var t string
	for _, f := range query.Filter {
		if len(f.Val) > 0 {
			if len(t) > 0 {
				t += "\nand "
			}
			t += fmt.Sprintf("lower(%s) like lower('%%%s%%')", f.Field, f.Val)
		}
	}
	if len(t) > 0 {
		wh = fmt.Sprintf("where %s", t)
		return
	}
	return
}

func orderVcToNum(s_in string) (res string) {
	res = strings.Replace(s_in, "price", "price::numeric", -1)
	return
}

func getOrderBy(query *DataQuery) (by string) {
	if query == nil {
		return
	}

	ord := orderVcToNum(strings.Join(query.OrderBy, ","))

	if len(ord) > 0 {
		by = fmt.Sprintf("order by %s", ord)
		return
	}
	return
}

func getPageSQl(query *DataQuery, sql string) string {
	if query == nil {
		return ""
	}
	count := 0
	cnt_sql := fmt.Sprintf("select count(1) from (%s) z", sql)
	Db.QueryRow(cnt_sql).Scan(&count)
	query.CntPage = int(math.Ceil(float64(count) / float64(query.PerPg)))
	if query.Page > query.CntPage {
		query.Page = 1
	}
	return fmt.Sprintf("limit %d offset %d", query.PerPg, query.PerPg*(query.Page-1))
}
