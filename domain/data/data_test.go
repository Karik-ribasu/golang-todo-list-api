package data_test

import (
	"github.com/Karik-ribasu/golang-todo-list-api/domain/entity"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/util"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database bootstrap", func() {
	It("Can open a connection to database", func() {
		Expect(errConn).To(BeNil())
	})
})

var _ = Describe("ListItem repo operations", func() {
	var userID int64 = 1
	listItemUUID := uuid.NewString()

	It("Can get list items from listItem repo", func() {
		listItems, err := dbManager.ListItemRepo().GetListItems(1)
		Expect(err).To(BeNil())
		Expect(len(listItems)).NotTo(BeZero())
	})

	It("Can insert a list item into listItemRepo", func() {
		listItem := entity.ListItem{
			UserID:       userID,
			ListItemUUID: listItemUUID,
			Title:        util.RandomString(10),
			Description:  util.RandomString(20),
		}

		err := dbManager.ListItemRepo().CreateListItem(listItem.UserID, listItem.ListItemUUID, listItem.Title, listItem.Description)
		Expect(err).Should(BeNil())
	})

	It("Can update a list item into listItemRepo", func() {
		listItem := entity.ListItem{
			UserID:       userID,
			ListItemUUID: listItemUUID,
			Title:        util.RandomString(10),
			Description:  util.RandomString(20),
			Active:       util.RandomBool(),
		}

		err := dbManager.ListItemRepo().UpdateListItem(listItem.UserID, listItem.ListItemUUID, listItem.Title, listItem.Description, listItem.Active)
		Expect(err).Should(BeNil())
	})

	It("Can delete a list item into listItemRepo", func() {
		listItem := entity.ListItem{
			UserID:       userID,
			ListItemUUID: listItemUUID,
			Title:        util.RandomString(10),
			Description:  util.RandomString(20),
			Active:       util.RandomBool(),
		}

		err := dbManager.ListItemRepo().DeleteListItem(listItem.UserID, listItem.ListItemUUID)
		Expect(err).Should(BeNil())
	})
})

var _ = Describe("User repo operations", func() {
	userUUID := uuid.NewString()
	userNickName := util.RandomString(6)
	userPassword := util.RandomBytes(10)

	It("Can insert a new user from User repo", func() {
		err := dbManager.UserRepo().CreateUser(userUUID, userNickName, userPassword)
		Expect(err).To(BeNil())
	})

	It("Can get user data by nickname from UserRepo", func() {
		user, err := dbManager.UserRepo().GetUserByNickName(userNickName)
		Expect(err).Should(BeNil())
		Expect(user.UserID).ShouldNot(BeZero())
	})

	It("Can get user data by userUUID from UserRepo", func() {
		user, err := dbManager.UserRepo().GetUserByUUID(userUUID)
		Expect(err).Should(BeNil())
		Expect(user.UserID).ShouldNot(BeZero())
	})
})
