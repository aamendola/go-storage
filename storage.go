package storage

// ObjectStorage ...
type ObjectStorage interface {
	Post(filenameSource, filenameDestination string)
}
