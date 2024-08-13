package commands

func PostInstall() string {
	return "restorecon -v /usr/local/bin/k3s"
}
