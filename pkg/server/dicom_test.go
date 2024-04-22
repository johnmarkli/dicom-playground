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

	"github.com/gorilla/mux"
	"github.com/johnmarkli/dime/pkg/server"
	"github.com/johnmarkli/dime/pkg/store"
	"github.com/stretchr/testify/assert"
)

const (
	testSOPInstanceUID = "1.3.12.2.1107.5.2.6.24119.30000013121716094326500000395"
	testDataPath       = "../../testdata/IM000001-mri"
)

func TestDICOMHandlerUpload(t *testing.T) {
	st := uploadDICOM(t, testDataPath)
	assert.NotNil(t, st)
}

func TestDICOMHandlerRead(t *testing.T) {
	st := uploadDICOM(t, testDataPath)
	assert.NotNil(t, st)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/dicom", nil)
	r = mux.SetURLVars(r, map[string]string{"id": testSOPInstanceUID})

	h := server.NewDICOMHandler(st)
	h.Read(w, r)
	defer w.Result().Body.Close()
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"sopInstanceUID":"1.3.12.2.1107.5.2.6.24119.30000013121716094326500000395"}`, string(body))
}

func TestDICOMHandlerAttributes(t *testing.T) {
	st := uploadDICOM(t, testDataPath)
	assert.NotNil(t, st)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", fmt.Sprintf("/dicom/%s/attributes/?tag=(0002,0000)&tag=(0008,0016)", testSOPInstanceUID), nil)
	r = mux.SetURLVars(r, map[string]string{"id": testSOPInstanceUID})

	h := server.NewDICOMHandler(st)
	h.Attributes(w, r)
	defer w.Result().Body.Close()
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	expected := `[
  {"tag":{"Group":2,"Element":0},"VR":4,"rawVR":"UL","valueLength":4,"value":[186]},
  {"tag":{"Group":8,"Element":22},"VR":0,"rawVR":"UI","valueLength":26,"value":["1.2.840.10008.5.1.4.1.1.4"]}
  ]`
	assert.JSONEq(t, expected, string(body))
}

func TestDICOMHandlerImage(t *testing.T) {
	st := uploadDICOM(t, testDataPath)
	assert.NotNil(t, st)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", fmt.Sprintf("/dicom/%s/image/", testSOPInstanceUID), nil)
	r = mux.SetURLVars(r, map[string]string{"id": testSOPInstanceUID})

	h := server.NewDICOMHandler(st)
	h.Image(w, r)
	defer w.Result().Body.Close()
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, "image/png", w.Result().Header.Get("Content-Type"))

	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.Len(t, body, 127594) // byte length of test png
}

func TestDICOMHandlerList(t *testing.T) {
	st := uploadDICOM(t, testDataPath)
	assert.NotNil(t, st)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/dicom", nil)

	h := server.NewDICOMHandler(st)
	h.List(w, r)
	defer w.Result().Body.Close()

	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `[{"sopInstanceUID":"1.3.12.2.1107.5.2.6.24119.30000013121716094326500000395"}]`, string(body))
}

// uploadDICOM uploads a DICOM file
func uploadDICOM(t *testing.T, filePath string) *store.MemStore {
	file, err := os.Open(filePath)
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

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/dicom", &b)
	r.Header.Add("Content-Type", mw.FormDataContentType())

	st, err := store.NewMemStore()
	assert.NoError(t, err)

	h := server.NewDICOMHandler(st)
	h.Upload(w, r)
	defer w.Result().Body.Close()

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)

	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"sopInstanceUID":"1.3.12.2.1107.5.2.6.24119.30000013121716094326500000395"}`, string(body))
	return st
}
