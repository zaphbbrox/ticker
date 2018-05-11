package api_test

import (
	"testing"
	"encoding/json"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"

	"git.codecoop.org/systemli/ticker/internal/api"
	"git.codecoop.org/systemli/ticker/internal/model"
	"git.codecoop.org/systemli/ticker/internal/storage"
)

func TestGetMessages(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/messages").
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, r.Code, 200)
		assert.Equal(t, r.Body.String(), `{"data":{"messages":[]},"status":"success","error":null}`)
	})
}

func TestGetMessage(t *testing.T) {
	r := setup()

	r.GET("/v1/admin/messages/1").
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, r.Code, 404)
		assert.Equal(t, r.Body.String(), `{"data":{},"status":"error","error":{"code":1001,"message":"not found"}}`)
	})
}

func TestPostMessage(t *testing.T) {
	r := setup()

	body := `{
		"text": "message",
		"ticker": 1
	}`

	r.POST("/v1/admin/messages").
		SetBody(body).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, r.Code, 200)

		type jsonResp struct {
			Data   map[string]model.Message `json:"data"`
			Status string                   `json:"status"`
			Error  interface{}              `json:"error"`
		}

		var jres jsonResp

		err := json.Unmarshal(r.Body.Bytes(), &jres)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.ResponseSuccess, jres.Status)
		assert.Equal(t, nil, jres.Error)
		assert.Equal(t, 1, len(jres.Data))

		message := jres.Data["message"]

		assert.Equal(t, "message", message.Text)
		assert.Equal(t, 1, message.Ticker)
	})
}

func TestPutMessage(t *testing.T) {
	r := setup()

	message := model.NewMessage()
	message.Text = "Text"
	message.Ticker = 1

	storage.DB.Save(&message)

	body := `{
		"text": "New Text",
		"ticker": 1
	}`

	r.PUT("/v1/admin/messages/100").
		SetBody(body).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 404, r.Code)
	})

	r.PUT("/v1/admin/messages/1").
		SetBody(body).
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)

		type jsonResp struct {
			Data   map[string]model.Message `json:"data"`
			Status string                   `json:"status"`
			Error  interface{}              `json:"error"`
		}

		var jres jsonResp

		err := json.Unmarshal(r.Body.Bytes(), &jres)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.ResponseSuccess, jres.Status)
		assert.Equal(t, nil, jres.Error)
		assert.Equal(t, 1, len(jres.Data))

		message := jres.Data["message"]

		assert.Equal(t, 1, message.ID)
		assert.Equal(t, "New Text", message.Text)
	})
}

func TestDeleteMessage(t *testing.T) {
	r := setup()

	message := model.NewMessage()
	message.Text = "Text"
	message.Ticker = 1

	storage.DB.Save(&message)

	r.DELETE("/v1/admin/messages/2").
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 404, r.Code)
	})

	r.DELETE("/v1/admin/messages/1").
		Run(api.API(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, 200, r.Code)

		type jsonResp struct {
			Data   map[string]model.Message `json:"data"`
			Status string                   `json:"status"`
			Error  interface{}              `json:"error"`
		}

		var jres jsonResp

		err := json.Unmarshal(r.Body.Bytes(), &jres)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, model.ResponseSuccess, jres.Status)
		assert.Nil(t, jres.Data)
		assert.Nil(t, jres.Error)
	})
}