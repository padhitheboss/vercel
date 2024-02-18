package kubejob

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	queuemodel "github.com/padhitheboss/kube-job-scheduler/pkg/queueHelper/model"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func (k *Kubejob) Initialize() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	// Create Kubernetes client
	k.kubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	k.maxExecutions, _ = strconv.Atoi(os.Getenv("MAX_EXECUTIONS"))
	k.namespace = os.Getenv("KUBE_NAMESPACE")
}
func (k *Kubejob) CreateJobTemplate(req queuemodel.Request) {
	envVars := []corev1.EnvVar{
		{Name: "REQUEST_ID", Value: req.RequestId},
		{Name: "BLOB_ACCOUNT_NAME", Value: os.Getenv("BLOB_ACCOUNT_NAME")},
		{Name: "BLOB_ACCOUNT_KEY", Value: os.Getenv("BLOB_ACCOUNT_KEY")},
		{Name: "SOURCE_BLOB_CONTAINER_NAME", Value: os.Getenv("SOURCE_BLOB_CONTAINER_NAME")},
		{Name: "DESTINATION_BLOB_CONTAINER_NAME", Value: os.Getenv("DESTINATION_BLOB_CONTAINER_NAME")},
		{Name: "DESTINATION_BLOB_FOLDER_PATH", Value: os.Getenv("DESTINATION_BLOB_FOLDER_PATH")},
		{Name: "DOWNLOAD_FOLDER_PATH", Value: os.Getenv("DOWNLOAD_FOLDER_PATH")},
		{Name: "BLOB_INPUT_FOLDER_PATH", Value: req.BlobUrl},
		{Name: "BUILD_FOLDER_PATH", Value: os.Getenv("BUILD_FOLDER_PATH")},
		{Name: "PROJECT_TYPE", Value: req.ProjectType},
	}

	// Define the container spec for the Kubernetes Job
	container := corev1.Container{
		Name:  "build-server",
		Image: os.Getenv("BUILD_IMAGE_NAME"),
		Env:   envVars,
		// Add any other necessary container configurations here
	}

	// Define the pod template spec for the Kubernetes Job
	podTemplate := corev1.PodTemplateSpec{
		Spec: corev1.PodSpec{
			Containers:    []corev1.Container{container},
			RestartPolicy: corev1.RestartPolicyNever, // Set the restart policy
		},
	}

	// Define the Job object metadata
	jobName := fmt.Sprintf("%s-%s", req.RequestId, time.Now().Format("20060102-150405"))
	jobMetadata := metav1.ObjectMeta{Name: jobName}

	// Define the Job object spec
	jobSpec := batchv1.JobSpec{
		Template: podTemplate,
	}

	// Define the Job object
	k.jobDefination = &batchv1.Job{
		ObjectMeta: jobMetadata,
		Spec:       jobSpec,
	}

	fmt.Println(k.jobDefination)
}

// func (k *Kubejob) CreateJobTemplate(req queuemodel.Request) {
// 	envVars := []corev1.EnvVar{
// 		{Name: "REQUEST_ID", Value: req.RequestId},
// 		{Name: "BLOB_ACCOUNT_NAME", Value: os.Getenv("BLOB_ACCOUNT_NAME")},
// 		{Name: "BLOB_ACCOUNT_KEY", Value: os.Getenv("BLOB_ACCOUNT_KEY")},
// 		{Name: "SOURCE_BLOB_CONTAINER_NAME", Value: os.Getenv("SOURCE_BLOB_CONTAINER_NAME")},
// 		{Name: "DESTINATION_BLOB_CONTAINER_NAME", Value: os.Getenv("DESTINATION_BLOB_CONTAINER_NAME")},
// 		{Name: "DESTINATION_BLOB_FOLDER_PATH", Value: os.Getenv("DESTINATION_BLOB_FOLDER_PATH")},
// 		{Name: "DOWNLOAD_FOLDER_PATH", Value: os.Getenv("DOWNLOAD_FOLDER_PATH")},
// 		{Name: "BLOB_INPUT_FOLDER_PATH", Value: req.BlobUrl},
// 		{Name: "BUILD_FOLDER_PATH", Value: os.Getenv("BUILD_FOLDER_PATH")},
// 		{Name: "PROJECT_TYPE", Value: req.ProjectType},
// 	}

// 	// Define the container spec for the Kubernetes Job
// 	container := corev1.Container{
// 		Name:  "build-server",
// 		Image: os.Getenv("BUILD_IMAGE_NAME"),
// 		Env:   envVars,
// 		// Add any other necessary container configurations here
// 	}

// 	// Define the pod template spec for the Kubernetes Job
// 	podTemplate := corev1.PodTemplateSpec{
// 		Spec: corev1.PodSpec{
// 			Containers: []corev1.Container{container},
// 		},
// 	}

//		// Define the Job object
//		k.jobDefination = &batchv1.Job{
//			ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("%s-%s", k.jobName, time.Now().GoString())},
//			Spec: batchv1.JobSpec{
//				Template: podTemplate,
//			},
//		}
//		fmt.Println(k.jobDefination)
//	}
func (k *Kubejob) RunJob() {
	_, err := k.kubeClient.BatchV1().Jobs(k.namespace).Create(context.Background(), k.jobDefination, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error triggering Kubernetes Job: %v\n", err)
		return
	}
	fmt.Println("Kubernetes Job triggered successfully.")
}
