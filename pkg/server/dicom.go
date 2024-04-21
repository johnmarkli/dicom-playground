package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/johnmarkli/dime/pkg/store"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

// DICOMHandler handles requests for DICOM management
type DICOMHandler struct {
	store store.Store
}

// NewDICOMHandler returns a new DICOMHandler
func NewDICOMHandler(store store.Store) *DICOMHandler {
	return &DICOMHandler{store}
}

// Upload a DICOM image
func (d *DICOMHandler) Upload(w http.ResponseWriter, r *http.Request) {
	slog.Info("Uploading DICOM...")

	// Get file upload
	r.ParseMultipartForm(10 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer file.Close()

	slog.Info("Uploaded DICOM",
		slog.String("filename", header.Filename),
		slog.Int64("size", header.Size),
		slog.String("header", fmt.Sprintf("%+v", header.Header)))

	// Parse dicom file
	dataset, err := dicom.Parse(file, header.Size, nil)
	if err != nil {
		panic(err)
	}

	dcm := store.NewDICOM(&dataset)
	err = d.store.Create(dcm)
	if err != nil {
		panic(err)
	}

	var jsonBytes []byte
	jsonBytes, err = json.Marshal(dcm)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

// Read a DICOM image
func (d *DICOMHandler) Read(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	dcm, err := d.store.Read(id)
	if err != nil {
		panic(err)
	}
	var jsonBytes []byte
	jsonBytes, err = json.Marshal(dcm)
	w.Write(jsonBytes)
}

// Read a DICOM image
func (d *DICOMHandler) Image(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	b, err := d.store.GetImage(id)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(b)
}

// Attributes from a DICOM image
func (d *DICOMHandler) Attributes(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	dcm, err := d.store.Read(id)
	if err != nil {
		panic(err)
	}
	var elements []*dicom.Element
	tags := r.URL.Query()["tag"]
	for _, strTag := range tags {
		trimmedTag := strings.Trim(strTag, "()")
		split := strings.Split(trimmedTag, ",")
		if len(split) != 2 {
			panic(err)
		}
		tagGroup, err := strconv.ParseUint(split[0], 16, 16)
		if err != nil {
			panic(err)
		}
		tagElement, err := strconv.ParseUint(split[1], 16, 16)
		if err != nil {
			panic(err)
		}
		t := tag.Tag{Group: uint16(tagGroup), Element: uint16(tagElement)}
		el, err := dcm.Dataset().FindElementByTag(t)
		if err != nil {
			panic(err)
		}
		if el != nil {
			elements = append(elements, el)
		}
	}

	var jsonBytes []byte
	jsonBytes, err = json.Marshal(elements)
	w.Write(jsonBytes)
}

// List DICOMS
func (d *DICOMHandler) List(w http.ResponseWriter, r *http.Request) {
	dicoms, err := d.store.List()
	if err != nil {
		panic(err)
	}
	var jsonBytes []byte
	jsonBytes, err = json.Marshal(dicoms)
	if err != nil {
		panic(err)
		// InternalServerErrorHandler(w, r)
		// return
	}
	w.Write(jsonBytes)
}
