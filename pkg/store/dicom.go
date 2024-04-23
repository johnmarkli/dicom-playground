package store

import (
	"fmt"
	"image"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

// DICOM is a model that represents a DICOM image
type DICOM struct {
	ID                string `json:"id" example:"1.3.12.2.1107.5.2.6.24119.30000013121716094326500000436"`
	StudyInstanceUID  string `json:"studyInstanceUID" example:"1.2.840.114202.4.833393677.4209323108.691055951.3610221745"`
	SeriesInstanceUID string `json:"seriesInstanceUID" example:"1.3.12.2.1107.5.2.6.24119.30000013121716094326500000394"`
	dataset           *dicom.Dataset
}

// NewDICOM returns a new DICOM instance
func NewDICOM(dataset *dicom.Dataset) (*DICOM, error) {
	var id, studyInstanceUID, seriesInstanceUID string

	// Take SOP Instance UID as unique ID for DICOM
	element, err := dataset.FindElementByTag(tag.SOPInstanceUID)
	if err != nil {
		return nil, fmt.Errorf("failed to find sop instance uid: %w", err)
	}
	uids := element.Value.GetValue().([]string)
	if len(uids) > 0 {
		id = uids[0] // assume first UID
	}

	// Get Study Instance UID and Series Instance UID
	element, err = dataset.FindElementByTag(tag.StudyInstanceUID)
	if err != nil {
		return nil, fmt.Errorf("failed to find study instance uid: %w", err)
	}
	uids = element.Value.GetValue().([]string)
	if len(uids) > 0 {
		studyInstanceUID = uids[0] // assume first UID
	}
	element, err = dataset.FindElementByTag(tag.SeriesInstanceUID)
	if err != nil {
		return nil, fmt.Errorf("failed to find series instance uid: %w", err)
	}
	uids = element.Value.GetValue().([]string)
	if len(uids) > 0 {
		seriesInstanceUID = uids[0] // assume first UID
	}

	return &DICOM{
		ID:                id,
		StudyInstanceUID:  studyInstanceUID,
		SeriesInstanceUID: seriesInstanceUID,
		dataset:           dataset,
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
