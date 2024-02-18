package kubejob

import (
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes"
)

type Kubejob struct {
	kubeClient    *kubernetes.Clientset
	jobName       string
	namespace     string
	jobDefination *batchv1.Job
	maxExecutions int
}
