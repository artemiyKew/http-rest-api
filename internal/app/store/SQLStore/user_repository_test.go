package sqlstore_test

import (
	"testing"

	"github.com/artemiyKew/http-rest-api/internal/app/model"
	"github.com/artemiyKew/http-rest-api/internal/app/store"
	"github.com/artemiyKew/http-rest-api/internal/app/store/SQLStore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	err := s.User().Create(model.TestUser(t))
	assert.NoError(t, err)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	email := "user@example.org"
	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.User().Create(model.TestUser(t))
	u, err := s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
