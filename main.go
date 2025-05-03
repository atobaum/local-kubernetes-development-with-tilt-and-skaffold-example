package main

import (
	"context"
	"fmt"
	"net/http"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pods, err := getPodList()
		if err != nil {
			http.Error(w, "Error retrieving pods", http.StatusInternalServerError)
		}
		fmt.Fprintln(w, "안녕하쇼! 이것은 Skaffold와 함께 돌아가는 Go 서버여라~")
		fmt.Fprintln(w, "현재 Pod 목록:")
		for _, pod := range pods {
			fmt.Fprintf(w, "- %s\n", pod.Name)
		}
	})

	http.HandleFunc("/configmap", func(w http.ResponseWriter, r *http.Request) {
		cm, err := craateTestConfigMap()
		if err != nil {
			http.Error(w, "Error creating configmap", http.StatusInternalServerError)
			fmt.Printf("error: %v\n", err)
			return
		}
		fmt.Fprintln(w, "ConfigMap 생성 완료!")
		fmt.Fprintf(w, "ConfigMap 이름: %s\n", cm.Name)
		fmt.Fprintf(w, "ConfigMap 데이터: %s\n", cm.Data["test-key"])
	})

	http.HandleFunc("/increase-configmap", func(w http.ResponseWriter, r *http.Request) {
		cm, err := increseConfigmap()
		if err != nil {
			http.Error(w, "Error creating configmap", http.StatusInternalServerError)
			fmt.Printf("error: %v\n", err)
			return
		}

		fmt.Fprintln(w, "ConfigMap 업데이트 완료!")
		fmt.Fprintf(w, "ConfigMap 이름: %s\n", cm.Name)
		fmt.Fprintf(w, "ConfigMap 데이터: %s\n", cm.Data["test-key"])
	})

	port := "8080"

	fmt.Printf("Servier started. Listening http://0.0.0.0:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func getPodList() ([]v1.Pod, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	pods, err := client.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}

func craateTestConfigMap() (*v1.ConfigMap, error) {
	configMapName := "test-config"

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	cm, err := client.CoreV1().ConfigMaps("default").Get(context.TODO(), configMapName, metav1.GetOptions{})

	if err != nil && !kerrors.IsNotFound(err) {
		fmt.Println(err.Error())
		return nil, err
	}

	if kerrors.IsNotFound(err) {

		configMap := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      configMapName,
				Namespace: "default",
			},
			Data: map[string]string{
				"test-key": fmt.Sprintf("%d", 1),
			},
		}

		return client.CoreV1().ConfigMaps("default").Create(context.TODO(), configMap, metav1.CreateOptions{})
	}

	return cm, nil
}

func increseConfigmap() (*v1.ConfigMap, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	cm, err := craateTestConfigMap()

	if err != nil {
		return nil, err
	}

	var curVal int
	_, err = fmt.Sscanf(cm.Data["test-key"], "%d", &curVal)
	if err != nil {
		curVal = 0
	}

	curVal++
	cm.Data["test-key"] = fmt.Sprintf("%d", curVal)

	return client.CoreV1().ConfigMaps("default").Update(context.TODO(), cm, metav1.UpdateOptions{})
}
