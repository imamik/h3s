package kubeconfig

import (
	"errors"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"gopkg.in/yaml.v3"
	"h3s/internal/cluster"
	"h3s/internal/config/path"
	"h3s/internal/k8s/token"
	"h3s/internal/utils/file"
	"h3s/internal/utils/ssh"
	"strings"
)

func Download(
	ctx *cluster.Cluster,
	proxy *hcloud.Server,
	remote *hcloud.Server,
) error {
	kubeConfig, err := get(ctx, proxy, remote)
	if err != nil {
		return err
	}
	if kubeConfig == nil {
		return errors.New("kubeconfig is nil")
	}

	err = saveKubeConfig(*kubeConfig)
	if err != nil {
		return err
	}
	return nil
}

func saveKubeConfig(kubeConfig KubeConfig) error {
	p := string(path.KubeConfigFileName)
	_, err := file.New(p).SetYaml(kubeConfig).Save()
	return err
}

func getUser(ctx *cluster.Cluster) (*User, error) {
	userName := "admin-user"
	userToken, err := token.Create(ctx, "kubernetes-dashboard", userName, 365*24)
	if err != nil {
		return nil, err
	}
	// Ensure the token is a single line string
	userToken = strings.ReplaceAll(userToken, "\n", "")
	user := User{
		Name: "default",
		User: UserDetails{
			Token: userToken,
		},
	}
	return &user, nil
}

func get(ctx *cluster.Cluster, proxy *hcloud.Server, remote *hcloud.Server) (*KubeConfig, error) {
	cmd := "sudo cat /etc/rancher/k3s/k3s.yaml"
	kubeConfigStr, err := ssh.ExecuteViaProxy(ctx, proxy, remote, cmd)
	if err != nil {
		return nil, err
	}

	var config KubeConfig
	err = yaml.Unmarshal([]byte(kubeConfigStr), &config)
	if err != nil {
		return nil, err
	}

	for i := range config.Clusters {
		config.Clusters[i].Cluster = ClusterDetails{
			InsecureSkipTLSVerify: ctx.Config.CertManager.Production,
			Server:                "https://k3s." + ctx.Config.Domain,
		}
	}

	user, err := getUser(ctx)
	if err != nil {
		return nil, err
	}
	config.Users = []User{*user}

	return &config, nil
}
