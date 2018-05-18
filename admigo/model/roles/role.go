package roles

const (
	ADMIN  = 1
	WRITER = 2
	READER = 3
)

type RoleModel struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"nm"`
}
