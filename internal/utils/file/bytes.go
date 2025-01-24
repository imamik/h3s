// Package file contains the functionality for managing files
package file

// SetBytes sets the data of the File
func (f *File) SetBytes(bytes []byte) *File {
	f.data = bytes
	return f
}

// GetBytes returns the data of the File
func (f *File) GetBytes() ([]byte, error) {
	err := f.Error()
	if err != nil {
		return nil, err
	}
	return f.data, f.Error()
}
