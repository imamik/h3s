package components

import (
	"fmt"
	"strings"
)

func WaitForCRDs(component string, resources []string) string {
	waitCmd := "kubectl wait --for=condition=established --timeout=30s " + strings.Join(resources, " ") + " >/dev/null 2>&1"
	return fmt.Sprintf(`
echo "Waiting for CRDs of %s to be established"
for i in {1..5}; do
	echo "Attempt $i"
	if %s; then
		if [ "$i" -gt 1 ]; then
			sleep 10
		fi
		echo "Established successfully"
		exit 0
	fi
	sleep 10
done
echo "Timed out"
exit 1
`, component, waitCmd)
}

func WaitForNamespace(namespace string) string {
	return fmt.Sprintf(`
echo "Waiting for namespace %s to be established"
for i in {1..5}; do
	echo "Attempt $i"
	if kubectl get namespace %s >/dev/null 2>&1; then
		if [ "$i" -gt 1 ]; then
			sleep 10
		fi
		echo "Established successfully"
		exit 0
	fi
	sleep 10
done
echo "Timed out"
exit 1
`, namespace, namespace)
}
