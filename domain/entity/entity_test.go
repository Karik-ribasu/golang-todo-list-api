package entity

import (
	"encoding/json"
	"testing"
)

func TestListItemJSON(t *testing.T) {
	li := ListItem{ListItemUUID: "u", Title: "t", Description: "d", Active: true}
	b, err := json.Marshal(li)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out["list_item_uuid"] != "u" {
		t.Fatalf("%s", b)
	}
}

func TestUserFields(t *testing.T) {
	u := User{UserID: 1, UserUUID: "x", NickName: "n", Password: []byte("p")}
	if len(u.Password) != 1 {
		t.Fatal()
	}
}
