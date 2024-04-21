package server

import (
	"fmt"
	"image/png"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

type dicomHandler struct {
	server *Server
}

func (d *dicomHandler) UploadDICOM(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
	slog.Info("Uploading DICOM...")

	// Get file upload
	r.ParseMultipartForm(10 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// Parse dicom file
	dataset, err := dicom.Parse(file, header.Size, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(dataset)

	// Save dicom as png
	slog.Info("Saving DICOM as PNG")
	pixelDataElement, err := dataset.FindElementByTag(tag.PixelData)
	if err != nil {
		panic(err)
	}
	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	for i, fr := range pixelDataInfo.Frames {
		// fmt.Println(i, fr)
		img, err := fr.GetImage()
		if err != nil {
			panic(err)
		}

		f, err := os.Create(fmt.Sprintf("image_%d.png", i))
		if err != nil {
			panic(err)
		}
		err = png.Encode(f, img)
		if err != nil {
			panic(err)
		}
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}

	sopInstanceUID, err := dataset.FindElementByTag(tag.SOPInstanceUID)
	fmt.Println("sop instance UID", sopInstanceUID)

	// Save dicom to file
	defer file.Close()
	slog.Info("Uploaded DICOM",
		slog.String("filename", header.Filename),
		slog.Int64("size", header.Size),
		slog.String("header", fmt.Sprintf("%+v", header.Header)))
	tempFile, err := os.CreateTemp("dicoms", "upload-*.dicom")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	slog.Info("Saved DICOM")
}

func (d *dicomHandler) ListDICOMs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("here are some dicoms"))
}
