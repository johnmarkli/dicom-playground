# dime

`dime` (__DI__com__M__angement__E__ndpoint) is a web service designed to work with DICOM files. It accepts and stores DICOM files, extracts and returns DICOM header attributes, and converts DICOM files in to a PNG for web-based viewing.

A RESTful API exposes the following functionality:
- POST /dicoms - upload DICOM file to be saved
- GET /dicoms - list DICOMs saved
- GET /dicoms/:id/tags - get dicom tags by ID
- GET /dicoms/:id/image - get dicom image by ID

TODO
- should this be 1 endpoint instead of broken down?
- testing
- swagger API
- error handling
- logging / log levels
- configurable port and other options
- dockerfile to build image

Nice to have
- simple Angular front end to use API
