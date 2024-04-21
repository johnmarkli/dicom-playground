package store

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/suyashkumar/dicom"
)

const (
	dicomDir = "dicom"
	pngDir   = "png"
)

// FileStore stores DICOM images on a file system
type FileStore struct {
	dir string
}

// NewFileStore creates a FileStore
func NewFileStore(dir string) *FileStore {
	return &FileStore{dir}
}

// Create a DICOM image in the file system
func (fs *FileStore) Create(d *DICOM) error {

	// save DICOM to file system
	dcmFile, err := os.Create(filepath.Join(fs.dir, dicomDir, fmt.Sprintf("%s.dcm", d.SOPInstanceUID)))
	if err != nil {
		panic(err)
	}
	defer dcmFile.Close()
	dicom.Write(dcmFile, *d.dataset)

	// save PNG to file system
	pngFile, err := os.Create(filepath.Join(fs.dir, pngDir, fmt.Sprintf("%s.png", d.SOPInstanceUID)))
	if err != nil {
		panic(err)
	}
	defer pngFile.Close()
	pngImg := d.Image()
	if pngImg != nil {
		err = png.Encode(pngFile, *d.Image())
		if err != nil {
			panic(err)
		}
	}

	return nil
}

// Read a DICOM image from the file system by SOP Instance UID
func (fs *FileStore) Read(sopInstanceUID string) (*DICOM, error) {
	fi, err := os.Stat(filepath.Join(fs.dir, dicomDir, fmt.Sprintf("%s.dcm", sopInstanceUID)))
	if err != nil {
		return nil, err
	}
	dataset, err := dicom.ParseFile(filepath.Join(fs.dir, dicomDir, fi.Name()), nil)
	if err != nil {
		return nil, err
	}
	dcm := NewDICOM(&dataset)
	return dcm, nil
}

// GetImage gets DICOM image as a reader
func (fs *FileStore) GetImage(sopInstanceUID string) ([]byte, error) {
	b, err := os.ReadFile(filepath.Join(fs.dir, pngDir, fmt.Sprintf("%s.png", sopInstanceUID)))
	if err != nil {
		return nil, err
	}
	return b, nil
}

// List DICOM images from the file system by SOP Instance UID
func (fs *FileStore) List() ([]*DICOM, error) {
	var dicoms []*DICOM
	files, _ := os.ReadDir(filepath.Join(fs.dir, dicomDir))
	for _, file := range files {
		dataset, err := dicom.ParseFile(filepath.Join(fs.dir, dicomDir, file.Name()), nil)
		if err != nil {
			return dicoms, err
		}
		dcm := NewDICOM(&dataset)
		dicoms = append(dicoms, dcm)
	}
	return dicoms, nil
}
