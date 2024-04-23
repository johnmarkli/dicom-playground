package store

import (
	"fmt"
	"image"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

// DICOM is a model that represents a DICOM image
type DICOM struct {
	ID      string `json:"id" example:"1.3.12.2.1107.5.2.6.24119.30000013121716094326500000436"`
	dataset *dicom.Dataset
}

// NewDICOM returns a new DICOM instance
func NewDICOM(dataset *dicom.Dataset) (*DICOM, error) {

	// Take SOP Instance UID as unique ID for DICOM
	sopInstanceUIDElement, err := dataset.FindElementByTag(tag.SOPInstanceUID)
	if err != nil {
		return nil, fmt.Errorf("failed to find sop instance uid: %w", err)
	}
	uids := sopInstanceUIDElement.Value.GetValue().([]string)
	var uid string
	if len(uids) > 0 {
		uid = uids[0] // assume first UID
	}
	return &DICOM{
		ID:      uid,
		dataset: dataset,
	}, nil
}

// Dataset returns the DICOM as a *dicom.Dataset
func (d *DICOM) Dataset() *dicom.Dataset {
	return d.dataset
}

// Image returns the DICOM as an image.Image
func (d *DICOM) Image() (*image.Image, error) {
	pixelDataElement, err := d.dataset.FindElementByTag(tag.PixelData)
	if err != nil {
		return nil, fmt.Errorf("failed to find pixel data: %w", err)
	}
	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	frames := pixelDataInfo.Frames
	if len(frames) > 0 {
		img, err := frames[0].GetImage() // assuming image is in first frame
		if err != nil {
			return nil, fmt.Errorf("failed to get image from frame: %w", err)
		}
		return &img, nil
	}
	return nil, nil
}
