package credentials

func Config() {
	projectName := surveyName()
	projectCredentials := surveyCredentials()
	save(projectName, projectCredentials)
}
