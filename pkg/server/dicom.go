package server

import (
	"encoding/json"
	"errors"
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
//
//	@Summary		Upload a DICOM image
//	@Description	Uploads a DICOM image to the server
//	@Tags			dicoms
//	@Accept			mpfd
//	@Produce		json
//	@Success		201	{object}	store.DICOM
//	@Failure		500	{object}	string
//	@Router			/dicoms [post]
func (d *DICOMHandler) Upload(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			handleError(rec, w)
		}
	}()

	// Get file upload
	err := r.ParseMultipartForm(10 << 20) // limit of 10MB files
	if err != nil {
		panic(fmt.Errorf("failed to parse form: %w", err))
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	slog.Info("Uploaded file",
		slog.String("filename", header.Filename),
		slog.Int64("size", header.Size),
		slog.String("header", fmt.Sprintf("%+v", header.Header)))

	// Parse dicom file
	dataset, err := dicom.Parse(file, header.Size, nil)
	if err != nil {
		panic(err)
	}

	// Create and store DICOM
	dcm, err := store.NewDICOM(&dataset)
	if err != nil {
		panic(err)
	}
	err = d.store.Create(dcm)
	if err != nil {
		panic(err)
	}
	slog.Info("Saved DICOM", slog.String("id", dcm.ID))

	// Return DICOM info
	var jsonBytes []byte
	jsonBytes, err = json.Marshal(dcm)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(jsonBytes)
}

// Read a DICOM image
//
//	@Summary		Read a DICOM image
//	@Description	Read a DICOM image from the server by SOP Instance UID
//	@Tags			dicoms
//	@Produce		json
//	@Param			id	path		string	true	"DICOM SOP Instance UID"
//	@Success		200	{object}	store.DICOM
//	@Failure		404	{object}	string
//	@Failure		500	{object}	string
//	@Router			/dicoms/{id} [get]
func (d *DICOMHandler) Read(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			handleError(rec, w)
		}
	}()

	// Get DICOM
	id := mux.Vars(r)["id"]
	dcm, err := d.store.Read(id)
	if err != nil {
		panic(err)
	}

	// Return DICOM info
	var jsonBytes []byte
	jsonBytes, err = json.Marshal(dcm)
	if err != nil {
		panic(err)
	}
	_, _ = w.Write(jsonBytes)
}

// Attributes from a DICOM image
//
//	@Summary		Get attributes from DICOM image
//	@Description	Get attributes from a DICOM image by tag
//	@Tags			dicoms
//	@Produce		json
//	@Param			id	path		string	true	"DICOM SOP Instance UID"
//	@Success		200	{array}		dicom.Element
//	@Failure		404	{object}	string
//	@Failure		500	{object}	string
//	@Router			/dicoms/{id}/attributes [get]
func (d *DICOMHandler) Attributes(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			handleError(rec, w)
		}
	}()

	// Get DICOM
	id := mux.Vars(r)["id"]
	dcm, err := d.store.Read(id)
	if err != nil {
		panic(err)
	}

	// Find DICOM attributes by tags
	var elements []*dicom.Element
	tags := r.URL.Query()["tag"]
	for _, strTag := range tags {
		trimmedTag := strings.Trim(strTag, "()")
		split := strings.Split(trimmedTag, ",")
		if len(split) != 2 {
			panic(fmt.Errorf("invalid tag format"))
		}
		tagGroup, err := strconv.ParseUint(split[0], 16, 16)
		if err != nil {
			panic(fmt.Errorf("tag could not be parsed: %w", err))
		}
		tagElement, err := strconv.ParseUint(split[1], 16, 16)
		if err != nil {
			panic(fmt.Errorf("tag could not be parsed: %w", err))
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

	// Return attribute data
	var jsonBytes []byte
	jsonBytes, err = json.Marshal(elements)
	if err != nil {
		panic(err)
	}
	_, _ = w.Write(jsonBytes)
}

// Image returns the DICOM image as a PNG
//
//	@Summary		Get DICOM image as a PNG
//	@Description	Get DICOM imange as a PNG
//	@Tags			dicoms
//	@Produce		png
//	@Param			id	path		string	true	"DICOM SOP Instance UID"
//	@Failure		404	{object}	string
//	@Failure		500	{object}	string
//	@Router			/dicoms/{id}/image [get]
func (d *DICOMHandler) Image(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			handleError(rec, w)
		}
	}()

	// Get DICOM Image
	id := mux.Vars(r)["id"]
	b, err := d.store.GetImage(id)
	if err != nil {
		panic(err)
	}

	// Return DICOM Image
	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(b)
}

// List DICOMS
//
//	@Summary		List DICOMs
//	@Description	List DICOMs on the server
//	@Tags			dicoms
//	@Produce		json
//	@Success		200	{array}		store.DICOM
//	@Failure		500	{object}	string
//	@Router			/dicoms [get]
func (d *DICOMHandler) List(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			handleError(rec, w)
		}
	}()

	// Get DICOMS
	dicoms, err := d.store.List()
	if err != nil {
		panic(err)
	}

	// Return DICOMS
	var jsonBytes []byte
	jsonBytes, err = json.Marshal(dicoms)
	if err != nil {
		panic(err)
	}
	_, _ = w.Write(jsonBytes)
}
func handleError(rec any, w http.ResponseWriter) {
	errVal, ok := rec.(error)
	if !ok {
		errVal = fmt.Errorf("%v", rec)
	}
	slog.Error(errVal.Error())
	if errors.Is(errVal, store.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("404 Not Found"))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(errVal.Error()))
	}
}
