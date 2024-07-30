package credentials

func Config() {
	projectName := surveyName()
	projectCredentials := surveyCredentials()
	Save(projectName, projectCredentials)
}
