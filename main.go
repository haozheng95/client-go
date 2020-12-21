//software: GoLand
//file: main.go
//time: 2020-12-17 17:36
package main

import (
	"context"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
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
	//var dirPth *string
	var namespace *string
	var name *string
	//var files = []string{"ca.crt", "client.yaml", "servers.conf"}
	//var data = make(map[string][]byte)
	//var secretType = v1.SecretTypeOpaque

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	//dirPth = flag.String("dirPth", "", "")
	namespace = flag.String("namespace", "kube-system", "")
	name = flag.String("name", "bitfusion-secret", "")

	flag.Parse()

	//config, err := rest.InClusterConfig()
	//if err != nil {
	//	panic(err.Error())
	//}

	//use the current context in kubeconfig
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
	//for _, fi := range files {
	//	fiPth := filepath.Join(*dirPth, fi)
	//	if content, err := ioutil.ReadFile(fiPth); err == nil {
	//		data[fi] = content
	//	} else {
	//		fmt.Println(err)
	//		panic(err)
	//	}
	//
	//}
	//
	//// assembly the secret
	//secret := &v1.Secret{
	//	Data: data,
	//	Type: secretType,
	//	ObjectMeta: metav1.ObjectMeta{
	//		Name:      *name,
	//		Namespace: *namespace,
	//	},
	//}

	_, err = clientset.CoreV1().Secrets(*namespace).Get(context.TODO(), *name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Secrets %s  not found in default namespace  %s \n", *name, *namespace)
		secret, err := clientset.CoreV1().Secrets("kube-system").Get(context.TODO(), *name, metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			panic(err.Error())
		}
		// create the secret
		_, err = clientset.CoreV1().Secrets("kube-system").Create(context.TODO(), secret, metav1.CreateOptions{})
		if err != nil {
			fmt.Println(err)
			panic(err.Error())
		}
	}

}
