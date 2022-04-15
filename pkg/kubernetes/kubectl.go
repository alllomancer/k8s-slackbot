package kubernetes

import (
	"fmt"
	"os/exec"
)

func RunLogs(kubeconfig, pod, line string) (string, error) {

	cmdline := "/usr/local/bin/kubectl"
	out, err := exec.Command(cmdline, "--kubeconfig="+kubeconfig, "logs", "--tail="+line, pod).Output()
	return fmt.Sprintf("%s", out), err
}

func RunGet(kubeconfig string) (string, error) {

	cmdline := "/usr/local/bin/kubectl"
	out, err := exec.Command(cmdline, "--kubeconfig="+kubeconfig, "get", "pods", "--all-namespaces").Output()
	return fmt.Sprintf("%s", out), err
}

func RunExec(kubeconfig, pod string) (string, error) {

	cmdline := "/usr/local/bin/kubectl"
	out, err := exec.Command(cmdline, "--kubeconfig="+kubeconfig, "exec", "-it", "-n", "kube-system", pod, "--", "curl", "http://localhost/version").Output()
	return fmt.Sprintf("%s", out), err
}

func RunGetNodes(kubeconfig string) (string, error) {

	cmdline := "/usr/local/bin/kubectl"
	out, err := exec.Command(cmdline, "--kubeconfig="+kubeconfig, "get", "nodes").Output()
	return fmt.Sprintf("%s", out), err
}
