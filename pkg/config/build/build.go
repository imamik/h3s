package build

import (
	"fmt"
	"hcloud-k3s-cli/pkg/config"
)

func InitConfig() (config.Config, error) {
	var conf config.Config

	surveyConf, err := surveyConfig()
	if err != nil {
		fmt.Println(err)
		return conf, err
	}

	conf.Name = surveyConf.Name
	conf.K3sVersion = surveyConf.K3sVersion
	conf.NetworkZone = config.NetworkZone(surveyConf.NetworkZone)
	conf.ControlPlanePool = config.ControlPlanePool{
		Nodes:        surveyConf.ControlPlanePoolNodes,
		Location:     config.Location(surveyConf.ControlPlanePoolLocation),
		AsWorkerPool: surveyConf.ControlPlanePoolAsWorkerPool,
		LoadBalancer: surveyConf.ControlPlanePoolLoadBalancer,
	}

	surveyWorkerPools, err := surveyWorkerPools(surveyConf.WorkerPools)
	if err != nil {
		fmt.Println(err)
		return conf, err
	}

	conf.WorkerPools = make([]config.NodePool, 0, len(surveyWorkerPools))
	for _, pool := range surveyWorkerPools {
		conf.WorkerPools = append(conf.WorkerPools, config.NodePool{
			Name:     pool.Name,
			Nodes:    pool.Nodes,
			Location: config.Location(pool.Location),
		})
	}

	return conf, nil
}
