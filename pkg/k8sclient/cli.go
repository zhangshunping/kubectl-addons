package k8sclient

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"kubectl-addons/pkg/utils"
	"os"
	"path/filepath"
)

var (
	kubeconfig *string
	Printer    = utils.PrintColor{}
)

type Cli struct {
	ClientSet *kubernetes.Clientset
}

// init k8sClient
func Initcli() (Cli, error) {

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, "kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	fmt.Println(*kubeconfig)
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	c := Cli{ClientSet: clientset}
	if err != nil {
		panic(err.Error())
	}
	return c, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// return --> node slice
func (client *Cli) ListNode(Labelselector string) (*v1.NodeList, error) {
	nodeList, err := client.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{LabelSelector: Labelselector})
	if err != nil {
		return nil, err
	}
	return nodeList, nil
}

//return --> pod slice
func (client *Cli) ListPod() (*v1.PodList, error) {
	podList, err := client.ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return podList, nil
}

// function --> check give annotation node and return
// return --> Fit annotation node  ` []*v1.Node{}  `And  `number of annotationed node`
func (cli *Cli) ReturnAnnoNode(ctx context.Context, GivenAnnotations map[string]string, Labelselector string, f string) ([]*v1.Node, int) {
	givenAnnoNodeSlice := []*v1.Node{}
	nodelist, _ := cli.ListNode(Labelselector)

	switch f {

	case "all":
		for i := 0; i < len(nodelist.Items); i++ {
			givenAnnoNodeSlice = append(givenAnnoNodeSlice, &nodelist.Items[i])
		}
	case "select":
		for i := 0; i < len(nodelist.Items); i++ {
			nodeAnnoMap := nodelist.Items[i].Annotations
			for GivenAk, GivenAv := range GivenAnnotations {
				v, ok := nodeAnnoMap[GivenAk]
				if ok {
					if v == GivenAv {
						givenAnnoNodeSlice = append(givenAnnoNodeSlice, &nodelist.Items[i])
					}
				}
			}
		}
	}
	return givenAnnoNodeSlice, len(givenAnnoNodeSlice)

}

// Print annotation Node
// list all annotation and list what has given annotation  in table.

func (cli *Cli) AnnoNodePrint(nodelist []*v1.Node, ctx context.Context, annotationMap map[string]string, nodelabel string) {
	var s string
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	if len(annotationMap) == 0 {
		// 考虑map
		t.AppendHeader(table.Row{"id", "Nodename", "Nodeip", "ALl", "Annotations"})
		for i := 0; i < len(nodelist); i++ {
			items := nodelist[i]
			annotationMap = items.Annotations
			annotaionstring, _ := json.Marshal(annotationMap)
			t.AppendRow([]interface{}{i + 1, items.Name, items.Status.Addresses[0].Address, "ALl", string(annotaionstring)})
		}

		s = "k8s nodes has annotations:"
	} else {
		t.AppendHeader(table.Row{"id", "Nodename", "Nodeip", "exit_or_no", "Given_Annotation"})
		annotaionstring, _ := json.Marshal(annotationMap)
		for i := 0; i < len(nodelist); i++ {
			items := nodelist[i]
			t.AppendRow([]interface{}{i + 1, items.Name, items.Status.Addresses[0].Address, "Yes", string(annotaionstring)})
		}

		s = fmt.Sprintf("k8s nodes containing %s is", string(annotaionstring))
	}

	Printer.Normal().Println(color.BlueString(s))
	t.AppendFooter(table.Row{"", "Total NOde", len(nodelist), ""})
	t.Render()

}
