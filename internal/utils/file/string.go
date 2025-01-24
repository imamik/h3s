package file

// SetString sets the file content to the given string.
func (f *File) SetString(str string) *File {
	f.SetBytes([]byte(str))
	return f
}

// GetString returns the file content as a string.
func (f *File) GetString() (string, error) {
	err := f.Error()
	if err != nil {
		return "", err
	}
	return string(f.data), f.Error()
}
