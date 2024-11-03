package credentials

func Configure() error {
	projectCredentials, err := surveyCredentials()
	if err != nil {
		return err
	}
	err = SaveCredentials(projectCredentials)
	if err != nil {
		return err
	}
	return nil
}
