# dime

`dime` (**DI**com**M**angement**E**ndpoint) is a small web service designed to work with DICOM files. It accepts and stores DICOM files, extracts and returns DICOM header attributes, and converts DICOM files in to a PNG for web-based viewing.

A RESTful API exposes the following functionality:
- POST /dicoms - upload DICOM file to be saved
- GET  /dicoms - list DICOMs saved
- GET  /dicoms/:id/attributes - get dicom header attributes by ID and tag
- GET  /dicoms/:id/image - get dicom image by ID
- GET  /health - server health check

### TODO
- swagger API
- error handling
- logging / log levels
- configurable port and other options
- dockerfile to build image
