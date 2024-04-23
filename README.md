# dime

`dime` (**di**com **m**angement **e**ndpoint) is a small web service designed to work with DICOM files. It accepts and stores DICOM files, extracts and returns DICOM header attributes, and converts DICOM files in to a PNG for web-based viewing.

A RESTful API exposes the following functionality:
- `POST /dicoms` - upload dicom file to be saved
- `GET  /dicoms` - list dicoms saved
- `GET  /dicoms/:id/attributes?tag=<tag1>&tag=<tagN>` - get dicom header attributes by ID and tags
- `GET  /dicoms/:id/image` - get dicom image by ID
- `GET  /health` - server health check
- `GET  /swagger` - API docs

## Getting Started

Install

```
go install
```

Run with default port 8080 and data directory `data`
```
dime
```

Run with environment variables
```
DIME_PORT=8081 DIME_DATA_DIR=/tmp dime
```

## Testing

Unit and integration tests
```
make test
```

With the server running, upload a directory of DICOMs
```
./scripts/upload-dir.sh -d <upload directory> -u <dime url>
```
