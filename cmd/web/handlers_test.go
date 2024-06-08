package main

import (
	"net/http"
	"testing"

	"github.com/ptakpatryk/list-me/internals/assert"
	"github.com/ptakpatryk/list-me/internals/models/mocks"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
  defer ts.Close()

	statusCode, _, body := ts.get(t, "/ping")

	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, string(body), "OK")
}

func TestSnippetView(t *testing.T) {
  app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
  defer ts.Close()

  tests := []struct {
    name string
    urlPath string
    wantCode int
    wantBody string
  }{
    {
      name: "Valid ID",
      urlPath: "/list/view/1",
      wantCode: http.StatusOK,
      wantBody: mocks.MockList.Title,
    },
    {
      name: "non-existent ID",
      urlPath: "/list/view/2",
      wantCode: http.StatusNotFound,
    },
    {
      name: "negative ID",
      urlPath: "/list/view/-1",
      wantCode: http.StatusNotFound,
    },
    {
      name: "decimal ID",
      urlPath: "/list/view/1.23",
      wantCode: http.StatusNotFound,
    },
    {
      name: "string ID",
      urlPath: "/list/view/foo",
      wantCode: http.StatusNotFound,
    },
    {
      name: "empty ID",
      urlPath: "/list/view/",
      wantCode: http.StatusNotFound,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      code, _, body := ts.get(t, tt.urlPath)

      assert.Equal(t, code, tt.wantCode)
      if tt.wantBody != "" {
        assert.StringContains(t, body, tt.wantBody)
      }
    })
  }
}
