package cmd

// 主服务
// 一个便捷的模块

import (
	"fmt"
	"io"
	"kube-debug/pkg/util/interrupt"
	"kube-debug/pkg/util/term"

	dockerterm "github.com/docker/docker/pkg/term"
	"k8s.io/kubectl/pkg/scheme"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"

	corev1 "k8s.io/api/core/v1"
)

// StreamOptions
type StreamOptions struct {
	// Namespace     string
	// PodName       string
	// ContainerName string
	Stdin bool // 一定为true
	TTY   bool //一定为true
	// // minimize unnecessary output
	// Quiet bool
	// InterruptParent, if set, is used to handle interrupts while attached
	InterruptParent *interrupt.Handler

	genericclioptions.IOStreams

	// for testing
	overrideStreams func() (io.ReadCloser, io.Writer, io.Writer)
	isTerminalIn    func(t term.TTY) bool
}

// DebugOptions 参数
type DebugOptions struct {
	StreamOptions

	args []string

	Config *rest.Config

	configFlags *genericclioptions.ConfigFlags
	Builder     *resource.Builder

	clientset *kubernetes.Clientset
	// restset   *restclient.RESTClient
}

var (
	example = `
		# debug a container in the running pod, the first container will be picked by default
		kubectl debug POD_NAME
	`
	errNoContext = fmt.Errorf("no context is currently set, use %q to select a new one", "kubectl config use-context <context>")
)

// NewDebugOptions provides an instance of DebugOptions with default values
func NewDebugOptions(streams genericclioptions.IOStreams) *DebugOptions {
	return &DebugOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		// IOStreams: streams,
		StreamOptions: StreamOptions{
			IOStreams: streams,
		},
	}
}

// 我这里创建一个debug的命令
// NewDebugCmd returns a cobra
func NewDebugCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewDebugOptions(streams)

	cmd := &cobra.Command{
		// Use:                   "debug POD [-c CONTAINER] -- COMMAND [args...]",
		Use: "demo POD [-c CONTAINER] -- COMMAND [args...]",
		// DisableFlagsInUseLine: true,
		Short: "Run a container in a running pod",
		// Long:                  longDesc,
		// Example:               example,
		// Version:               version.Version(),
		RunE: func(c *cobra.Command, args []string) error {

			if err := o.Complete(c, args); err != nil {
				return err
			}
			// if err := o.Validate(); err != nil {
			// 	return err
			// }
			// 前面的验证都不改， 直接看看能不能启动起一个pod
			if err := o.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	// cmd.Flags().BoolVar(&o.listNamespaces, "list", o.listNamespaces, "if true, print the list of all namespaces in the current KUBECONFIG")
	// cmd.Flags().StringVarP(&o.ContainerName, "container", "c", "",
	// 	"Target container to debug, default to the first container in pod")

	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

// Complete sets all information required for updating the current context
func (o *DebugOptions) Complete(cmd *cobra.Command, args []string) error {

	o.args = args
	fmt.Println(args)

	// var err error

	// configLoader := o.configFlags.ToRawKubeConfigLoader()
	// o.Namespace, _, err = configLoader.Namespace()
	// o.PodName = args[0]

	// 	o.Builder = resource.NewBuilder(o.configFlags)

	// 	if err != nil {
	// 		return err
	// 	}

	config, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}
	o.Config = config
	o.clientset, err = kubernetes.NewForConfig(config)

	if err != nil {
		return err
	}

	return nil
}

// // Validate ensures that all required arguments and flag values are provided
// func (o *DebugOptions) Validate() error {

// 	if len(o.PodName) == 0 {
// 		return fmt.Errorf("pod name required")
// 	}
// 	return nil
// }

// SetupTTY xx
func (o *DebugOptions) SetupTTY() term.TTY {

	t := term.TTY{
		Parent: o.InterruptParent,
		Out:    o.Out,
	}

	t.In = o.In
	// if we get to here, the user wants to attach stdin, wants a TTY, and o.In is a terminal, so we
	// can safely set t.Raw to true
	t.Raw = true

	if !t.IsTerminalIn() {
		if o.ErrOut != nil {
			fmt.Fprintln(o.ErrOut, "Unable to use a TTY - input is not a terminal or the right kind of file")
		}
		return t
	}
	stdin, stdout, _ := dockerterm.StdStreams()
	o.In = stdin
	t.In = stdin
	if o.Out != nil {
		o.Out = stdout
		t.Out = stdout
	}

	return t
}

// Run lists all available debugs on a user's KUBECONFIG or updates the
// current context based on a provided debug.
func (o *DebugOptions) Run() error {
	var err error
	// pod, err := o.clientset.CoreV1().Pods(o.Namespace).Get(o.PodName, metav1.GetOptions{})
	// if err != nil {
	// 	return err
	// }

	// containerName := pod.Spec.Containers[0].Name

	t := o.SetupTTY()
	var sizeQueue remotecommand.TerminalSizeQueue

	if t.Raw {
		// this call spawns a goroutine to monitor/update the terminal size
		sizeQueue = t.MonitorSize(t.GetSize())

		// 	// unset p.Err if it was previously set because both stdout and stderr go over p.Out when tty is
		// 	// true
		o.ErrOut = nil

	}

	pod := o.generateDebugNodePod("docker-desktop")
	pod, err = o.launchPod(pod)
	containerName := "test-debug"

	if err != nil {
		fmt.Println(2222)
		fmt.Println(err)
	}

	fn := func() error {
		req := o.clientset.CoreV1().RESTClient().Post().
			Resource("pods").
			Name(pod.Name).
			Namespace(pod.Namespace).
			SubResource("exec")
		req.VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   []string{"bash"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

		var executor remotecommand.Executor

		// fmt.Println(req)
		// 	// change to ws
		if executor, err = remotecommand.NewSPDYExecutor(o.Config, "POST", req.URL()); err != nil {
			fmt.Println(5675675)
			fmt.Println(err)
		}

		// 	// fmt.Println(req.URL())
		if err = executor.Stream(remotecommand.StreamOptions{
			Stdin:             o.In,
			Stdout:            o.Out,
			Stderr:            o.ErrOut,
			TerminalSizeQueue: sizeQueue,
			Tty:               true,
		}); err != nil {
			fmt.Println(err)
		}
		return nil
	}

	if err := t.Safe(fn); err != nil {
		fmt.Println(5675675)
		// 	fmt.Println(err)
		// 	return err
	}

	return nil
}
