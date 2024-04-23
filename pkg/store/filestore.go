package store

import (
	"errors"
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
func NewFileStore(dir string) (*FileStore, error) {

	// Create directories if they don't exist
	dirs := []string{dir, filepath.Join(dir, dicomDir), filepath.Join(dir, pngDir)}
	for _, dir := range dirs {
		err := createDirIfNotExist(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return &FileStore{dir}, nil
}

// Create a DICOM image in the file system along with PNG file
func (fs *FileStore) Create(dcm *DICOM) error {

	// save DICOM to file system
	dcmFile, err := os.Create(filepath.Join(fs.dir, dicomDir, fmt.Sprintf("%s.dcm", dcm.ID)))
	if err != nil {
		return fmt.Errorf("failed to create dicom file: %w", err)
	}
	defer dcmFile.Close()
	err = dicom.Write(dcmFile, *dcm.dataset)
	if err != nil {
		return fmt.Errorf("failed to write dicom file: %w", err)
	}

	// save PNG to file system
	pngFile, err := os.Create(filepath.Join(fs.dir, pngDir, fmt.Sprintf("%s.png", dcm.ID)))
	if err != nil {
		return fmt.Errorf("failed to create png file: %w", err)
	}
	defer pngFile.Close()
	pngImg, err := dcm.Image()
	if err != nil {
		return fmt.Errorf("failed to get dicom image: %w", err)
	}
	if pngImg != nil {
		err = png.Encode(pngFile, *pngImg)
		if err != nil {
			return fmt.Errorf("failed to encode png file: %w", err)
		}
	}
	return nil
}

// Read a DICOM image from the file system by SOP Instance UID
func (fs *FileStore) Read(id string) (*DICOM, error) {
	fi, err := os.Stat(filepath.Join(fs.dir, dicomDir, fmt.Sprintf("%s.dcm", id)))
	if err != nil {
		return nil, ErrNotFound
	}
	dataset, err := dicom.ParseFile(filepath.Join(fs.dir, dicomDir, fi.Name()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dicom file: %w", err)
	}
	dcm, err := NewDICOM(&dataset)
	if err != nil {
		return nil, err
	}
	return dcm, nil
}

// GetImage gets DICOM image as a byte array
func (fs *FileStore) GetImage(id string) ([]byte, error) {
	_, err := os.Stat(filepath.Join(fs.dir, pngDir, fmt.Sprintf("%s.png", id)))
	if err != nil {
		return nil, ErrNotFound
	}
	b, err := os.ReadFile(filepath.Join(fs.dir, pngDir, fmt.Sprintf("%s.png", id)))
	if err != nil {
		return nil, fmt.Errorf("failed to read png file: %w", err)
	}
	return b, nil
}

// List DICOM images from the file system by SOP Instance UID
func (fs *FileStore) List() ([]*DICOM, error) {
	dicoms := []*DICOM{}
	files, _ := os.ReadDir(filepath.Join(fs.dir, dicomDir))
	for _, file := range files {
		dataset, err := dicom.ParseFile(filepath.Join(fs.dir, dicomDir, file.Name()), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to parse dicom file: %w", err)
		}
		dcm, err := NewDICOM(&dataset)
		if err != nil {
			return nil, err
		}
		dicoms = append(dicoms, dcm)
	}
	return dicoms, nil
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
