// Package store provides a way to store DICOM images
package store

import "errors"

var (
	// ErrNotFound is an error for a DICOM that is not found
	ErrNotFound = errors.New("not found")
)

// Store is an interface for working with storage of DICOM images
type Store interface {

	// Create a new DICOM image
	Create(*DICOM) error

	// Read a DICOM image by SOP Instance UID
	Read(sopInstanceUID string) (*DICOM, error)

	// GetImage gets the DICOM image
	GetImage(sopInstanceUID string) ([]byte, error)

	// List DICOM images
	List() ([]*DICOM, error)
}
