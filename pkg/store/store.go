package store

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
