package store

import (
	"image"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

// DICOM is a model that represents a DICOM image
type DICOM struct {
	SOPInstanceUID string `json:"sopInstanceUID"`
	dataset        *dicom.Dataset
}

// NewDICOM returns a new DICOM instance
func NewDICOM(dataset *dicom.Dataset) *DICOM {
	sopInstanceUIDElement, err := dataset.FindElementByTag(tag.SOPInstanceUID)
	if err != nil {
		return nil
	}
	uids := sopInstanceUIDElement.Value.GetValue().([]string)
	var uid string
	if len(uids) > 0 {
		uid = uids[0]
	}
	return &DICOM{
		SOPInstanceUID: uid,
		dataset:        dataset,
	}
}

// Dataset returns the DICOM as a *dicom.Dataset
func (d *DICOM) Dataset() *dicom.Dataset {
	return d.dataset
}

// Image returns the DICOM as an image.Image
func (d *DICOM) Image() *image.Image {
	pixelDataElement, err := d.dataset.FindElementByTag(tag.PixelData)
	if err != nil {
		panic(err)
	}
	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	for _, fr := range pixelDataInfo.Frames {
		// assumes one frame
		img, err := fr.GetImage()
		if err != nil {
			panic(err)
		}
		return &img
	}
	return nil
}
