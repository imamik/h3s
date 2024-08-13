package credentials

func Configure() error {
	projectCredentials := surveyCredentials()
	err := SaveCredentials(projectCredentials)
	if err != nil {
		return err
	}
	return nil
}
