package credentials

func Configure() {
	projectName := surveyName()
	projectCredentials := surveyCredentials()
	SaveCredentials(projectName, projectCredentials)
}
