package server_test

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/johnmarkli/dime/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestServerNew(t *testing.T) {
	s, err := server.New()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	defer s.Shutdown()
}

func TestServerRoutes(t *testing.T) {
	s, err := server.New()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	defer s.Shutdown()

	router := s.Server().Handler

	// GET /health
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, r)
	res := w.Result()
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "OK", string(body))
	res.Body.Close()

	// POST /dicoms
	w = httptest.NewRecorder()
	file, err := os.Open(testDataPath)
	assert.NoError(t, err)
	assert.NotNil(t, file)
	defer file.Close()

	var b bytes.Buffer
	mw := multipart.NewWriter(&b)

	var fw io.Writer
	fw, err = mw.CreateFormFile("file", file.Name())
	assert.NoError(t, err)
	assert.NotNil(t, fw)

	_, err = io.Copy(fw, file)
	assert.NoError(t, err)
	mw.Close()

	r = httptest.NewRequest(http.MethodPost, "/dicoms", &b)
	r.Header.Add("Content-Type", mw.FormDataContentType())

	router.ServeHTTP(w, r)
	res = w.Result()
	body, err = io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, testDICOMjson, string(body))
	res.Body.Close()

	// GET /dicoms
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/dicoms", nil)
	router.ServeHTTP(w, r)
	res = w.Result()
	body, err = io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, fmt.Sprintf("[%s]", testDICOMjson), string(body))
	res.Body.Close()

	// GET /dicoms/:id
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/dicoms/%s", testID), nil)
	router.ServeHTTP(w, r)
	res = w.Result()
	body, err = io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, testDICOMjson, string(body))
	res.Body.Close()

	// GET /dicoms/:id/attributes
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/dicoms/%s/attributes?tag=(0002,0000)&tag=(0008,0016)", testID), nil)
	router.ServeHTTP(w, r)
	res = w.Result()
	body, err = io.ReadAll(res.Body)
	assert.NoError(t, err)
	expected := `[
  {"tag":{"Group":2,"Element":0},"VR":4,"rawVR":"UL","valueLength":4,"value":[186]},
  {"tag":{"Group":8,"Element":22},"VR":0,"rawVR":"UI","valueLength":26,"value":["1.2.840.10008.5.1.4.1.1.4"]}
  ]`
	assert.JSONEq(t, expected, string(body))
	res.Body.Close()

	// GET /dicoms/:id/image
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/dicoms/%s/image", testID), nil)
	router.ServeHTTP(w, r)
	res = w.Result()
	body, err = io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Len(t, body, 127594) // byte length of test png
	res.Body.Close()
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, "image/png", w.Result().Header.Get("Content-Type"))
}
