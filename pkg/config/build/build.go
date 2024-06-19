package build

import (
	"github.com/AlecAivazis/survey/v2"
	"hcloud-k3s-cli/pkg/config"
)

func InitConfig() (config.Config, error) {
	var conf config.Config

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
			Name: "controlPlanePool.name",
			Prompt: &survey.Input{
				Message: "What's the control plane pool name?",
			},
		},
		{
			Name: "controlPlanePool.nodes",
			Prompt: &survey.Input{
				Message: "How many nodes in the control plane pool?",
			},
		},
		{
			Name: "controlPlanePool.location",
			Prompt: &survey.Select{
				Message: "Choose a location for the control plane pool:",
				Options: []string{"nbg1", "fsn1", "hel1", "ash", "hil"},
			},
		},
		{
			Name: "controlPlanePool.asWorker",
			Prompt: &survey.Confirm{
				Message: "Should the control plane pool be used as a worker pool?",
			},
		},
		{
			Name: "controlPlanePool.loadBalancer",
			Prompt: &survey.Confirm{
				Message: "Should the control plane pool have a load balancer?",
			},
		},
		{
			Name: "combinedLoadBalancer",
			Prompt: &survey.Confirm{
				Message: "Do you want to combine the load balancer?",
			},
		},
	}

	err := survey.Ask(qs, &conf)
	if err != nil {
		return conf, err
	}

	var numWorkerPools int
	prompt := &survey.Input{
		Message: "How many worker pools do you want to create?",
	}
	survey.AskOne(prompt, &numWorkerPools, nil)

	conf.WorkerPools.WorkerPools = make([]config.NodePool, numWorkerPools)

	for i := 0; i < numWorkerPools; i++ {
		workerPool := config.NodePool{}
		qsWorkerPool := []*survey.Question{
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
		err := survey.Ask(qsWorkerPool, &workerPool)
		if err != nil {
			return conf, err
		}
		conf.WorkerPools.WorkerPools[i] = workerPool
	}

	return conf, nil
}
