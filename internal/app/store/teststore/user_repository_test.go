package teststore_test

import (
	"testing"

	"github.com/artemiyKew/http-rest-api/internal/app/model"
	"github.com/artemiyKew/http-rest-api/internal/app/store"
	"github.com/artemiyKew/http-rest-api/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s := teststore.New()
	err := s.User().Create(model.TestUser(t))
	assert.NoError(t, err)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststore.New()
	email := "user@example.org"
	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.User().Create(model.TestUser(t))
	u, err := s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
