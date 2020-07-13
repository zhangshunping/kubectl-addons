/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"kubectl-addons/pkg/k8sclient"
	"kubectl-addons/pkg/utils"
)

var (
	annotaion string
	ctx       context.Context
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get k8s node resouces",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
	TraverseChildren: true,
}

// define node resources
var nodeannoCmd = &cobra.Command{
	Use:   "nodeanno",
	Short: "get node Annotation",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		annotationMap := make(map[string]string)
		ctx, cancle := context.WithCancel(context.TODO())
		defer cancle()
		if 0 != len(annotaion) {
			print_anno_node(cmd, ctx, cli, annotationMap, nodelabel)
		} else {
			utils.Printer.Err().Printf(color.CyanString("'-a' ca|CA|Ca, all|All|ALL"))
		}

	},
	Example: utils.Printer.Tips().Sprintf(color.BlueString(" 1.kubect-addons get nodeanno -a \"CA\"  --> to get ClusterAutoSacler node that not to clam down \n" +
		"2. kubectl-addons get nodeanno -a \"All\" --> to get all Node Annotation\n" +
		"3. kubectl-addons get nodeanno -a '{\"flannel.alpha.coreos.com/backend-type\":\"vxlan\"}'  --> to list given Annotation Node  \n" +
		"4. kubectl-addons get nodeanno -a '{\"cluster-autoscaler.kubernetes.io/scale-down-disabled\":\"true\"}' -k C:/Users/39295/kube/config ",)),
}

// kubectl-addon get nodeanno  -a "ca，CA，Ca" ===> print annotation node
func print_anno_node(cmd *cobra.Command, ctx context.Context, cli k8sclient.Cli, annotationMap map[string]string, nodelabel string) {
	var nodelist []*v1.Node

	switch annotaion {

	case "ca", "CA", "Ca":
		annotationMap = utils.AnnotationMap["CA-NoDown"]
		nodelist, _ = cli.ReturnAnnoNode(ctx, annotationMap, nodelabel, "select")
	case "all", "All", "ALL":
		_ = json.Unmarshal([]byte(annotaion), &annotationMap)
		nodelist, _ = cli.ReturnAnnoNode(ctx, annotationMap, nodelabel, "all")
	case "":
		cmd.Help()
		fmt.Println(cmd.Example)
		return
	default:
		str:=[]byte(annotaion)
		_=json.Unmarshal(str,&annotationMap)
		fmt.Println(annotationMap)
		nodelist, _ = cli.ReturnAnnoNode(ctx, annotationMap, nodelabel, "select")
	}

	cli.AnnoNodePrint(nodelist, ctx, annotationMap, nodelabel)
}

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.AddCommand(nodeannoCmd)

	nodeannoCmd.Flags().StringVarP(&annotaion, "annotation", "a", "", " get nodeanno  -a '{\"flannel.alpha.coreos.com/backend-type\":\"vxlan\"}'\n ,to list gien node")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
