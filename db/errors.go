package db

import "errors"

// Files based error

var ErrFileIsNotFound = errors.New("File not found!")
var ErrDuplicateFile = errors.New("File with the same name already exists!")
var ErrFileFetch = errors.New("Error when loading your files!")
var ErrFileNameNotUpdated = errors.New("Failed to update filename!")
var ErrFileNotCreated = errors.New("Failed to upload your file")
var ErrFileNotUpdated = errors.New("Failed to update your file!")
var ErrFilePathNotUpdated = errors.New("Failed to update your file path!")
var ErrFileOrPathNotFound = errors.New("Cannot find the file or path provided!")
var ErrDeleteFilesInsideFolder = errors.New("error deleting files that are stored inside the folder!")
var ErrDeleteFiles = errors.New("Failed to delete files")
var ErrDeleteFolder = errors.New("Failed to delete folder")

var ErrInternal = errors.New("Something went wrong, we are checking...")
