package storage_test

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
)

func TestDeleteUpload(t *testing.T) {
	setup()

	upload := initialUploadTestData(t)
	err := storage.DeleteUpload(upload)
	if err != nil {
		t.Fail()
	}

	var u *model.Upload
	err = storage.DB.Find("ID", upload.ID, &u)
	if err == nil {
		t.Fail()
	}

	_, err = os.Open(upload.FullPath())
	if err == nil {
		t.Fail()
	}
}

func TestDeleteUploadNonExisting(t *testing.T) {
	setup()

	err := storage.DeleteUpload(&model.Upload{})
	if err == nil {
		t.Fail()
	}
}

func TestDeleteUploads(t *testing.T) {
	setup()

	upload := initialUploadTestData(t)
	uploads := []*model.Upload{upload}

	storage.DeleteUploads(uploads)

	var u *model.Upload
	err := storage.DB.Find("ID", upload.ID, &u)
	if err == nil {
		t.Fail()
	}
}

func TestFindUploadsByMessageNonExistingUpload(t *testing.T) {
	setup()

	message := model.NewMessage()
	attachment := model.Attachment{UUID: uuid.New().String(), Extension: "jpg", ContentType: "image/jpeg"}
	message.Attachments = []model.Attachment{attachment}
	err := storage.DB.Save(message)
	if err != nil {
		t.Fail()
	}

	u := storage.FindUploadsByMessage(message)
	assert.Equal(t, 0, len(u))
}

func TestDeleteUploadsByTicker(t *testing.T) {
	setup()

	ticker := model.NewTicker()
	_ = storage.DB.Save(ticker)

	err := storage.DeleteUploadsByTicker(ticker)

	assert.Nil(t, err)

	upload := model.NewUpload("name.jpg", "image/jpeg", ticker.ID)
	_ = storage.DB.Save(upload)

	err = storage.DeleteUploadsByTicker(ticker)

	assert.Nil(t, err)
}

func initialUploadTestData(t *testing.T) *model.Upload {
	upload := model.NewUpload("name.jpg", "image/jpeg", 1)
	err := storage.DB.Save(upload)
	if err != nil {
		t.Fail()
	}

	return upload
}
