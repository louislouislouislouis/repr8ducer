package k8s

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type K8sService struct {
	config *string
	Client *kubernetes.Clientset
}

func (s *K8sService) ListNamespace() (*v1.NamespaceList, error) {
	tests, err := s.Client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

	// for _, test := range tests.Items {
	// utils.Log.Debug().Msg(test.Name)
	// }
	return tests, err
}

func (s *K8sService) ListPodsInNamespace(nms string) (*v1.PodList, error) {
	pods, err := s.Client.CoreV1().Pods(nms).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return pods, err
}

func (s *K8sService) Exec() error {
	req := s.Client.CoreV1().
		RESTClient().
		Post().
		Resource("pods").
		Name("metering-66cfb8bc6b-frkfj").
		Namespace("kiwios-cloud-metering").
		SubResource("exec")
	req.VersionedParams(&v1.PodExecOptions{
		Container: "metering",
		Command:   []string{"cat", "config/application.yaml"},
		Stdin:     true,
		Stdout:    true,
	}, scheme.ParameterCodec)
	// find / -type f -name application.yaml 2>/dev/null | xargs cat

	kubeCfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	restCfg, err := kubeCfg.ClientConfig()

	exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", req.URL())

	fmt.Println(req.URL())
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:  nil,
		Tty:    false,
		Stdout: buf,
		Stderr: buf2,
	})
	if err != nil {
		return err
	}
	test3 := buf.String()

	// Analyse du YAML
	var data map[string]interface{}
	err = yaml.Unmarshal([]byte(test3), &data)
	if err != nil {
		log.Fatalf("Erreur lors de l'analyse YAML : %v", err)
	}

	// Transformation en slice de string
	flattened := flattenYaml("", data)

	// Impression du résultat
	fmt.Println(strings.Join(flattened, "\n"))

	return nil
}

// Fonction pour transformer une structure en slice de string
func flattenYaml(prefix string, data map[string]interface{}) []string {
	var result []string
	for key, value := range data {
		fullKey := prefix + key

		// Si la valeur est une autre map, on appelle récursivement flattenYaml
		if reflect.TypeOf(value).Kind() == reflect.Map {
			nestedMap := value.(map[string]interface{})
			result = append(result, flattenYaml(fullKey+".", nestedMap)...)
		} else {
			// Sinon, on ajoute le paramètre dans le format demandé
			formattedValue := fmt.Sprintf("--env \"%s='%v'\"", fullKey, value)
			result = append(result, formattedValue)
		}
	}
	return result
}
