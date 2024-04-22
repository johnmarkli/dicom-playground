# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.5]

### Updated

- Fix lint issues

## [0.0.4]

### Added

- Server routes integration tests

## [0.0.3]

### Added

- mem store for testing instead of file store
- Error handling and logging
- Refactored handler tests

## [0.0.2]

### Added

- store interface with file store to save DICOM data to filesystem
- Handlers for reading DICOM by SOP Instance UID with routes for reading attributes and getting the image
- Tests for handlers

## [0.0.1]

### Added

- Inital server pkg with healh and basic dicom handler
- dime cmd binary to run server
