// Package credentials provides the functionality to configure project credentials.
package credentials

// Configure prompts the user for various project credentials and saves them to the configuration file.
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
