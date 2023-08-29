package dto

type GetListItemsRequest struct {
	UserUUID string
}

type GetListItemsResponse []struct {
	ListItemUUID string `json:"list_item_uuid,omitempty"`
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
	Active       bool   `json:"active"`
}

type CreateListItemRequest struct {
	UserUUID    string
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateListItemResponse struct {
	ListItemUUID string `json:"list_item_uuid,omitempty"`
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
	Active       bool   `json:"active"`
}

type UpdateListItemRequest struct {
	UserUUID     string `json:"user_uuid"`
	ListItemUUID string `json:"list_item_uuid"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Active       bool   `json:"active"`
}

type UpdateListItemResponse struct {
	ListItemUUID string `json:"list_item_uuid,omitempty"`
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
	Active       bool   `json:"active"`
}

type DeleteListItemRequest struct {
	UserUUID     string `json:"user_uuid"`
	ListItemUUID string `json:"list_item_uuid"`
}
