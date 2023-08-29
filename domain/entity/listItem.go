package entity

type ListItem struct {
	ListItemID   int64
	UserID       int64
	ListItemUUID string
	Title        string
	Description  string
	Active       bool
}
