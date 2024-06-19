package build

import "github.com/AlecAivazis/survey/v2"

type SurveyConfig struct {
	Name                         string `yaml:"name"`
	K3sVersion                   string `yaml:"k3sVersion"`
	NetworkZone                  string `yaml:"region"`
	ControlPlanePoolLocation     string `yaml:"controlPlanePoolLocation"`
	ControlPlanePoolNodes        int    `yaml:"controlPlanePoolNodes"`
	ControlPlanePoolLoadBalancer bool   `yaml:"controlPlanePoolLoadBalancer"`
	ControlPlanePoolAsWorkerPool bool   `yaml:"controlPlanePoolAsWorker"`
	WorkerPools                  int    `yaml:"workerPools"`
}

func surveyConfig() (SurveyConfig, error) {
	var conf SurveyConfig

	qs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "What's the project name?",
			},
		},
		{
			Name: "k3sVersion",
			Prompt: &survey.Input{
				Message: "What's the K3s version?",
			},
		},
		{
			Name: "networkZone",
			Prompt: &survey.Select{
				Message: "Choose a network zone:",
				Options: []string{"eu-central", "us-east", "us-west"},
			},
		},
		{
			Name: "controlPlanePoolLocation",
			Prompt: &survey.Select{
				Message: "Choose a location for the control plane pool:",
				Options: []string{"nbg1", "fsn1", "hel1", "ash", "hil"},
			},
		},
		{
			Name: "controlPlanePoolNodes",
			Prompt: &survey.Input{
				Message: "How many nodes in the control plane pool?",
			},
		},
		{
			Name: "controlPlanePoolAsWorkerPool",
			Prompt: &survey.Confirm{
				Message: "Should the control plane pool be used as a worker pool?",
			},
		},
		{
			Name: "controlPlanePoolLoadBalancer",
			Prompt: &survey.Confirm{
				Message: "Should the control plane pool have a load balancer?",
			},
		},
		{
			Name: "workerPools",
			Prompt: &survey.Input{
				Message: "How many worker pools?",
			},
		},
	}

	err := survey.Ask(qs, &conf)

	return conf, err
}

type SurveyWorkerPool struct {
	Name     string `yaml:"name"`
	Nodes    int    `yaml:"nodes"`
	Location string `yaml:"location"`
}

func surveyWorkerPools(count int) ([]SurveyWorkerPool, error) {
	var workerPools []SurveyWorkerPool

	qs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "What's the worker pool name?",
			},
		},
		{
			Name: "nodes",
			Prompt: &survey.Input{
				Message: "How many nodes in the worker pool?",
			},
		},
		{
			Name: "location",
			Prompt: &survey.Select{
				Message: "Choose a location for the worker pool:",
				Options: []string{"nbg1", "fsn1", "hel1", "ash", "hil"},
			},
		},
	}

	for i := 0; i < count; i++ {
		var workerPool SurveyWorkerPool
		err := survey.Ask(qs, &workerPool)
		if err != nil {
			return workerPools, err
		}
		workerPools = append(workerPools, workerPool)
	}

	return workerPools, nil
}
