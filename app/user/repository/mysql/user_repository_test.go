package mysql_test

import (
	"context"
	"github.com/fajardm/ewallet-example/app/base"
	"github.com/fajardm/ewallet-example/app/user/model"
	"github.com/fajardm/ewallet-example/app/user/repository/mysql"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v2"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockUser := model.User{
		Model: base.Model{
			ID:        uuid.NewV4(),
			CreatedBy: uuid.NewV4(),
			CreatedAt: time.Now(),
		},
		Username:       "john",
		Email:          "john@gmail.com",
		MobilePhone:    "08199999999",
		HashedPassword: []byte("secret"),
	}

	prep := mock.ExpectExec("^INSERT (.+)")
	prep.WithArgs(mockUser.ID, mockUser.Username, mockUser.Email, mockUser.MobilePhone, mockUser.HashedPassword, mockUser.CreatedBy, mockUser.CreatedAt)
	prep.WillReturnResult(sqlmock.NewResult(1, 1))

	r := mysql.NewUserRepository(db)
	err = r.Store(context.TODO(), mockUser)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockUser := []model.User{
		model.User{
			Model: base.Model{
				ID:        uuid.NewV4(),
				CreatedBy: uuid.NewV4(),
				CreatedAt: time.Now(),
			},
			Username:       "john",
			Email:          "john@gmail.com",
			MobilePhone:    "08199999999",
			HashedPassword: []byte("secret"),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "mobile_phone", "hashed_password", "created_by", "created_at", "updated_by", "updated_at"})
	rows.AddRow(mockUser[0].ID, mockUser[0].Username, mockUser[0].Email, mockUser[0].MobilePhone, mockUser[0].HashedPassword, mockUser[0].CreatedBy, mockUser[0].CreatedAt, mockUser[0].UpdatedBy, mockUser[0].UpdatedAt)

	prep := mock.ExpectQuery("^SELECT (.+)")
	prep.WillReturnRows(rows)

	r := mysql.NewUserRepository(db)
	res, err := r.GetByID(context.TODO(), mockUser[0].ID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user := uuid.NewV4()
	now := time.Now()

	mockUser := model.User{
		Model: base.Model{
			ID:        uuid.NewV4(),
			CreatedBy: user,
			CreatedAt: now,
			UpdatedBy: &user,
			UpdatedAt: &now,
		},
		Username:       "john",
		Email:          "john@gmail.com",
		MobilePhone:    "08199999999",
		HashedPassword: []byte("secret"),
	}

	prep := mock.ExpectExec("^UPDATE (.+)")
	prep = prep.WithArgs(mockUser.Email, mockUser.HashedPassword, mockUser.UpdatedBy, mockUser.UpdatedAt, mockUser.ID)
	prep.WillReturnResult(sqlmock.NewResult(1, 1))

	r := mysql.NewUserRepository(db)
	err = r.Update(context.TODO(), mockUser)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := uuid.NewV4()

	prep := mock.ExpectExec("^DELETE (.+)")
	prep.WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))

	r := mysql.NewUserRepository(db)
	err = r.Delete(context.TODO(), id)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
