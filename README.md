# dime

`dime` (**DI**com**M**angement**E**ndpoint) is a small web service designed to work with DICOM files. It accepts and stores DICOM files, extracts and returns DICOM header attributes, and converts DICOM files in to a PNG for web-based viewing.

A RESTful API exposes the following functionality:
- `POST /dicoms` - upload dicom file to be saved
- `GET  /dicoms` - list dicoms saved
- `GET  /dicoms/:id/attributes` - get dicom header attributes by ID and tag
- `GET  /dicoms/:id/image` - get dicom image by ID
- `GET  /health` - server health check
- `GET  /swagger` - API docs
