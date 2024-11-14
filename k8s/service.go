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
	return s.Client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
}

func (s *K8sService) ListPodsInNamespace(nms string) (*v1.PodList, error) {
	return s.Client.CoreV1().Pods(nms).List(context.TODO(), metav1.ListOptions{})
}

func (s *K8sService) GetPod(nms, podName string) (*v1.Pod, error) {
	pod, err := s.Client.CoreV1().Pods(nms).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return pod, err
}

func (s *K8sService) GetContainerFromPods(podName, nms string) ([]v1.Container, error) {
	pod, err := s.Client.CoreV1().Pods(nms).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return []v1.Container{}, err
	}
	return pod.Spec.Containers, err
}

func (s *K8sService) Exec(nms, podName string) (string, error) {
	req := s.Client.CoreV1().
		RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(nms).
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
		return "", err
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
		return "", err
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
	// fmt.Println(strings.Join(flattened, "\n"))

	return strings.Join(flattened, "\n"), nil
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
