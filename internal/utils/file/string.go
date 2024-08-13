package file

func (f *File) SetString(str string) *File {
	f.SetBytes([]byte(str))
	return f
}

func (f *File) GetString() (string, error) {
	err := f.Error()
	if err != nil {
		return "", err
	}
	return string(f.data), f.Error()
}
