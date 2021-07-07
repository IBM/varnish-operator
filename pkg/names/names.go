package names

func ClusterRole(vcName, vcNamespace string) string {
	return vcName + "-varnish-clusterrole-" + vcNamespace
}

func ClusterRoleBinding(vcName, vcNamespace string) string {
	return vcName + "-varnish-clusterrolebinding-" + vcNamespace
}

func HeadlessService(vcName string) string {
	return vcName + "-headless-service"
}

func PodDisruptionBudget(vcName string) string {
	return vcName + "-varnish-pdb"
}

func Role(vcName string) string {
	return vcName + "-varnish-role"
}

func RoleBinding(vcName string) string {
	return vcName + "-varnish-rolebinding"
}

func NoCacheService(vcName string) string {
	return vcName + "-no-cache"
}

func StatefulSet(vcName string) string {
	return vcName + "-varnish"
}

func ServiceAccount(vcName string) string {
	return vcName + "-varnish-serviceaccount"
}

func VarnishSecret(vcName string) string {
	return vcName + "-varnish-secret"
}

func ServiceMonitor(vcName string) string {
	return vcName + "-prometheus-servicemonitor"
}

func GrafanaDashboard(vcName string) string {
	return vcName + "-grafana-dashboard"
}

func GrafanaDashboardFile(vcName string) string {
	return vcName + "-dashboard.json"
}

func HaproxyGrafanaDashboardFile(vcName string) string {
	return vcName + "haproxy-dashboard.json"
}
