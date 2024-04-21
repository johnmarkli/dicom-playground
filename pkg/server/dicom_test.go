package server_test

import (
	"bytes"
	"encoding/json"
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

func TestDICOMHandlerUpload(t *testing.T) {
	w := uploadDICOM(t, "../../testdata/IM000001-mri")
	defer w.Result().Body.Close()

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)

	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"sopInstanceUID":"1.3.12.2.1107.5.2.6.24119.30000013121716094326500000395"}`, string(body))
}

func TestDICOMHandlerRead(t *testing.T) {
	w := uploadDICOM(t, "../../testdata/IM000001-mri")
	defer w.Result().Body.Close()

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)

	var dcm store.DICOM
	err = json.Unmarshal(body, &dcm)
	assert.NoError(t, err)
	assert.Equal(t, "1.3.12.2.1107.5.2.6.24119.30000013121716094326500000395", dcm.SOPInstanceUID)

	w = httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/dicom", nil)
	r = mux.SetURLVars(r, map[string]string{"id": dcm.SOPInstanceUID})

	store := store.NewFileStore("data")
	h := server.NewDICOMHandler(store)
	h.Read(w, r)
	defer w.Result().Body.Close()

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	body, err = io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"sopInstanceUID":"1.3.12.2.1107.5.2.6.24119.30000013121716094326500000395"}`, string(body))
}

func TestDICOMHandlerAttributes(t *testing.T) {
	w := uploadDICOM(t, "../../testdata/IM000001-mri")
	defer w.Result().Body.Close()

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)

	var dcm store.DICOM
	err = json.Unmarshal(body, &dcm)
	assert.NoError(t, err)
	assert.Equal(t, "1.3.12.2.1107.5.2.6.24119.30000013121716094326500000395", dcm.SOPInstanceUID)

	w = httptest.NewRecorder()

	r := httptest.NewRequest("GET", fmt.Sprintf("/dicom/%s/attributes/?tag=(0002,0000)&tag=(0008,0016)", "1.3.12.2.1107.5.2.6.24119.30000013121716094326500000395"), nil)
	r = mux.SetURLVars(r, map[string]string{"id": dcm.SOPInstanceUID})

	store := store.NewFileStore("data")
	h := server.NewDICOMHandler(store)
	h.Attributes(w, r)
	defer w.Result().Body.Close()

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	body, err = io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	expected := `[
  {"tag":{"Group":2,"Element":0},"VR":4,"rawVR":"UL","valueLength":4,"value":[186]},
  {"tag":{"Group":8,"Element":22},"VR":0,"rawVR":"UI","valueLength":26,"value":["1.2.840.10008.5.1.4.1.1.4"]}
  ]`
	assert.JSONEq(t, expected, string(body))
}

func TestDICOMHandlerList(t *testing.T) {
	w := uploadDICOM(t, "../../testdata/IM000001-mri")
	defer w.Result().Body.Close()

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)

	w = httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/dicom", nil)

	store := store.NewFileStore("data")
	h := server.NewDICOMHandler(store)
	h.List(w, r)
	defer w.Result().Body.Close()

	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `[{"sopInstanceUID":"1.3.12.2.1107.5.2.6.24119.30000013121716094326500000395"}]`, string(body))
}

func uploadDICOM(t *testing.T, filePath string) *httptest.ResponseRecorder {
	// get dicom test data from file
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
	store := store.NewFileStore("data")
	h := server.NewDICOMHandler(store)
	h.Upload(w, r)
	return w
}
