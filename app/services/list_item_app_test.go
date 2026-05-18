package appServices

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Karik-ribasu/golang-todo-list-api/app/dto"
	"github.com/Karik-ribasu/golang-todo-list-api/domain/entity"
)

func TestListItem_GetListItems_EmptyUUID(t *testing.T) {
	svc := newListItemAppService(&stubListDom{}, &stubUserDom{})
	_, herr := svc.GetListItems(dto.GetListItemsRequest{})
	if herr == nil || herr.StatusCode != 404 {
		t.Fatal()
	}
}

func TestListItem_GetListItems_UserErr(t *testing.T) {
	u := &stubUserDom{err: sql.ErrNoRows}
	svc := newListItemAppService(&stubListDom{}, u)
	_, herr := svc.GetListItems(dto.GetListItemsRequest{UserUUID: "x"})
	if herr == nil || herr.StatusCode != 404 {
		t.Fatal()
	}
}

func TestListItem_GetListItems_ListErr(t *testing.T) {
	u := &stubUserDom{u: entity.User{UserID: 1}}
	l := &stubListDom{ge: sql.ErrNoRows}
	svc := newListItemAppService(l, u)
	_, herr := svc.GetListItems(dto.GetListItemsRequest{UserUUID: "x"})
	if herr == nil {
		t.Fatal()
	}
}

func TestListItem_GetListItems_OtherSQL(t *testing.T) {
	u := &stubUserDom{u: entity.User{UserID: 1}}
	l := &stubListDom{ge: errors.New("boom")}
	svc := newListItemAppService(l, u)
	_, herr := svc.GetListItems(dto.GetListItemsRequest{UserUUID: "x"})
	if herr == nil || herr.StatusCode != 503 {
		t.Fatal()
	}
}

func TestListItem_GetListItems_OK(t *testing.T) {
	u := &stubUserDom{u: entity.User{UserID: 1}}
	l := &stubListDom{items: []entity.ListItem{{ListItemUUID: "a", Title: "t", Description: "d", Active: true}}}
	svc := newListItemAppService(l, u)
	resp, herr := svc.GetListItems(dto.GetListItemsRequest{UserUUID: "x"})
	if herr != nil || len(resp) != 1 {
		t.Fatalf("%v %#v", herr, resp)
	}
}

func TestListItem_Create_EmptyUser(t *testing.T) {
	svc := newListItemAppService(&stubListDom{}, &stubUserDom{})
	_, herr := svc.CreateListItem(dto.CreateListItemRequest{})
	if herr == nil || herr.StatusCode != 404 {
		t.Fatal()
	}
}

func TestListItem_Create_UserErr(t *testing.T) {
	u := &stubUserDom{err: sql.ErrNoRows}
	svc := newListItemAppService(&stubListDom{}, u)
	_, herr := svc.CreateListItem(dto.CreateListItemRequest{UserUUID: "u", Title: "t"})
	if herr == nil {
		t.Fatal()
	}
}

func TestListItem_Create_RepoErr(t *testing.T) {
	u := &stubUserDom{u: entity.User{UserID: 1}}
	l := &stubListDom{ce: sql.ErrNoRows}
	svc := newListItemAppService(l, u)
	_, herr := svc.CreateListItem(dto.CreateListItemRequest{UserUUID: "u", Title: "t"})
	if herr == nil {
		t.Fatal()
	}
}

func TestListItem_Create_OK(t *testing.T) {
	u := &stubUserDom{u: entity.User{UserID: 1}}
	l := &stubListDom{}
	svc := newListItemAppService(l, u)
	resp, herr := svc.CreateListItem(dto.CreateListItemRequest{UserUUID: "u", Title: "t", Description: "d"})
	if herr != nil || resp.Title != "t" || resp.ListItemUUID == "" {
		t.Fatalf("%v %+v", herr, resp)
	}
}

func TestListItem_Update_NotFound(t *testing.T) {
	svc := newListItemAppService(&stubListDom{}, &stubUserDom{})
	_, herr := svc.UpdateListItem(dto.UpdateListItemRequest{})
	if herr == nil || herr.StatusCode != 404 {
		t.Fatal()
	}
}

func TestListItem_Update_UserErr(t *testing.T) {
	u := &stubUserDom{err: sql.ErrNoRows}
	svc := newListItemAppService(&stubListDom{}, u)
	_, herr := svc.UpdateListItem(dto.UpdateListItemRequest{UserUUID: "u", ListItemUUID: "i"})
	if herr == nil {
		t.Fatal()
	}
}

func TestListItem_Update_RepoErr(t *testing.T) {
	u := &stubUserDom{u: entity.User{UserID: 1}}
	l := &stubListDom{ue: sql.ErrNoRows}
	svc := newListItemAppService(l, u)
	_, herr := svc.UpdateListItem(dto.UpdateListItemRequest{UserUUID: "u", ListItemUUID: "i", Title: "t"})
	if herr == nil {
		t.Fatal()
	}
}

func TestListItem_Update_OK(t *testing.T) {
	u := &stubUserDom{u: entity.User{UserID: 1}}
	l := &stubListDom{}
	svc := newListItemAppService(l, u)
	resp, herr := svc.UpdateListItem(dto.UpdateListItemRequest{UserUUID: "u", ListItemUUID: "i", Title: "t", Description: "d", Active: false})
	if herr != nil || resp.Title != "t" || resp.Active != false {
		t.Fatalf("%v %+v", herr, resp)
	}
}

func TestListItem_Delete_NotFound(t *testing.T) {
	svc := newListItemAppService(&stubListDom{}, &stubUserDom{})
	if herr := svc.DeleteListItem(dto.DeleteListItemRequest{}); herr == nil {
		t.Fatal()
	}
}

func TestListItem_Delete_UserErr(t *testing.T) {
	u := &stubUserDom{err: sql.ErrNoRows}
	svc := newListItemAppService(&stubListDom{}, u)
	if herr := svc.DeleteListItem(dto.DeleteListItemRequest{UserUUID: "u", ListItemUUID: "i"}); herr == nil {
		t.Fatal()
	}
}

func TestListItem_Delete_RepoErr(t *testing.T) {
	u := &stubUserDom{u: entity.User{UserID: 1}}
	l := &stubListDom{de: sql.ErrNoRows}
	svc := newListItemAppService(l, u)
	if herr := svc.DeleteListItem(dto.DeleteListItemRequest{UserUUID: "u", ListItemUUID: "i"}); herr == nil {
		t.Fatal()
	}
}

func TestListItem_Delete_OK(t *testing.T) {
	u := &stubUserDom{u: entity.User{UserID: 1}}
	l := &stubListDom{}
	svc := newListItemAppService(l, u)
	if herr := svc.DeleteListItem(dto.DeleteListItemRequest{UserUUID: "u", ListItemUUID: "i"}); herr != nil {
		t.Fatal(herr)
	}
}

func TestListItem_GetListItems_JSONMarshalError(t *testing.T) {
	prev := jsonMarshalListItems
	defer func() { jsonMarshalListItems = prev }()
	jsonMarshalListItems = func(any) ([]byte, error) { return nil, errors.New("m") }
	u := &stubUserDom{u: entity.User{UserID: 1}}
	l := &stubListDom{items: []entity.ListItem{{}}}
	svc := newListItemAppService(l, u)
	_, herr := svc.GetListItems(dto.GetListItemsRequest{UserUUID: "x"})
	if herr == nil || herr.StatusCode != 500 {
		t.Fatal(herr)
	}
}

func TestListItem_GetListItems_JSONUnmarshalError(t *testing.T) {
	prevM := jsonMarshalListItems
	prevU := jsonUnmarshalListItems
	defer func() {
		jsonMarshalListItems = prevM
		jsonUnmarshalListItems = prevU
	}()
	jsonMarshalListItems = json.Marshal
	jsonUnmarshalListItems = func([]byte, any) error { return errors.New("u") }
	u := &stubUserDom{u: entity.User{UserID: 1}}
	l := &stubListDom{items: []entity.ListItem{{}}}
	svc := newListItemAppService(l, u)
	_, herr := svc.GetListItems(dto.GetListItemsRequest{UserUUID: "x"})
	if herr == nil || herr.StatusCode != 500 {
		t.Fatal(herr)
	}
}
