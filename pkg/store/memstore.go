package store

import (
	"bytes"
	"fmt"
	"image/png"
)

// MemStore stores DICOM images in memory
type MemStore struct {
	dicoms map[string]*DICOM
	pngs   map[string][]byte
}

// NewMemStore returns a MemStore
func NewMemStore() (*MemStore, error) {
	return &MemStore{
		dicoms: map[string]*DICOM{},
		pngs:   map[string][]byte{},
	}, nil
}

// Create DICOM image in memory store
func (ms *MemStore) Create(dcm *DICOM) error {
	ms.dicoms[dcm.ID] = dcm

	pngImg, err := dcm.Image()
	if err != nil {
		return fmt.Errorf("failed to get dicom image: %w", err)
	}
	if pngImg != nil {
		var b bytes.Buffer
		err := png.Encode(&b, *pngImg)
		if err != nil {
			return err
		}
		ms.pngs[dcm.ID] = b.Bytes()
	}
	return nil
}

// Read a DICOM image from the memory by SOP Instance UID
func (ms *MemStore) Read(id string) (*DICOM, error) {
	if dcm, ok := ms.dicoms[id]; ok {
		return dcm, nil
	}
	return nil, ErrNotFound
}

// GetImage gets DICOM image as a byte array
func (ms *MemStore) GetImage(id string) ([]byte, error) {
	if b, ok := ms.pngs[id]; ok {
		return b, nil
	}
	return []byte{}, ErrNotFound
}

// List DICOM images from the file system
func (ms *MemStore) List() ([]*DICOM, error) {
	dcms := []*DICOM{}
	for _, dcm := range ms.dicoms {
		dcms = append(dcms, dcm)
	}
	return dcms, nil
}
