//software: GoLand
//file: main.go
//time: 2020-12-17 17:36
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

// go build -o app .
// ./app -kubeconfig=/root/.kube/config -dirPth=/home/file/87lENOv

func main() {

	var kubeconfig *string
	var dirPth *string
	var namespace *string
	var name *string
	var files = []string{"ca.cert", "client.yaml", "servers.conf"}
	var data = make(map[string][]byte)
	var secretType = v1.SecretTypeOpaque

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	dirPth = flag.String("dirPth", "", "")
	namespace = flag.String("namespace", "kube-system", "")
	name = flag.String("name", "bitfusion-sercet", "")

	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// read the token file
	for _, fi := range files {
		fiPth := filepath.Join(*dirPth, fi)
		if content, err := ioutil.ReadFile(fiPth); err == nil {
			data[fi] = content
		} else {
			fmt.Println(err)
			panic(err)
		}

	}

	// assembly the secret
	secret := &v1.Secret{
		Data: data,
		Type: secretType,
		ObjectMeta: metav1.ObjectMeta{
			Name:      *name,
			Namespace: *namespace,
		},
	}

	// create the secret
	_, err = clientset.CoreV1().Secrets(*namespace).Create(context.TODO(), secret, metav1.CreateOptions{})

	if err != nil {
		panic(err.Error())
	}

}
