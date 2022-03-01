package cmd

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/watch"
	"k8s.io/kubernetes/pkg/client/conditions"
)

func (o *DebugOptions) generateDebugNodePod(nodeName string) *corev1.Pod {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nodeDebug",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "testDebug",
					// Env:   debugger.Env,
					Image: "nicolaka/netshoot:latest",
					// ImagePullPolicy:          core,
					Stdin:                    true,
					TTY:                      true,
					TerminationMessagePolicy: corev1.TerminationMessageReadFile,
					VolumeMounts: []corev1.VolumeMount{
						{
							MountPath: "/hostfs",
							Name:      "host-root",
						},
					},
				},
			},
			HostIPC:       true,
			HostNetwork:   true,
			HostPID:       true,
			NodeName:      nodeName,
			RestartPolicy: corev1.RestartPolicyNever,
			Volumes: []corev1.Volume{
				{
					Name: "host-root",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{Path: "/"},
					},
				},
			},
		},
	}

	return pod
}

// launchPod launch given pod until it's running
func (o *DebugOptions) launchPod(pod *corev1.Pod) (*corev1.Pod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	pod, err := o.clientset.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return pod, err
	}

	watcher, err := o.clientset.CoreV1().Pods(pod.Namespace).Watch(ctx, metav1.SingleObject(pod.ObjectMeta))
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(o.Out, "Waiting for pod %s to run...\n", pod.Name)
	event, err := watch.UntilWithoutRetry(ctx, watcher, conditions.PodRunning)
	if err != nil {
		fmt.Fprintf(o.ErrOut, "Error occurred while waiting for pod to run:  %v\n", err)
		return nil, err
	}
	pod = event.Object.(*corev1.Pod)
	return pod, nil
}
