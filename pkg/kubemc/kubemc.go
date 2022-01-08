package kubemc

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	api "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/yaml"
)

type KubeMCConfig struct {
	KubeConfigPath string
	KubeMCDir      string
	KubeMCExt      map[string]bool
	MergeRateLimit time.Duration // Only one merge can be trigger in n second
}

func NewKubeMCConfig() (k KubeMCConfig) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	// Default kubeConfigPath = "~/.kube/config"
	k.KubeConfigPath = homeDir + "/.kube/config"
	// Default kubeMCDir = "~/.kube/kubemc"
	k.KubeMCDir = homeDir + "/.kube/kubemc"
	//default configfile extension -> .yaml, yml
	k.KubeMCExt = map[string]bool{
		".yaml": true,
		".yml":  true,
	}
	//default mergeRateLimit = 2 second
	k.MergeRateLimit = 2

	// Custom KubeMCConfig if Environment value is set
	var env string
	var isEnvSet bool

	env, isEnvSet = os.LookupEnv("KUBECONFIG")
	if isEnvSet {
		k.KubeConfigPath = env
	}
	env, isEnvSet = os.LookupEnv("KUBEMC_DIR")
	if isEnvSet {
		k.KubeMCDir = env
	}
	env, isEnvSet = os.LookupEnv("KUBEMC_RATELIMIT")
	if isEnvSet {
		rate, err := strconv.Atoi(env)
		if err != nil {
			log.Fatal(err)
		}
		k.MergeRateLimit = time.Duration(time.Duration(rate).Seconds())
	}

	return
}

// Get the kubeconfig file list
func GetKubeConfigList(kmcc KubeMCConfig) (kubeConfigList []string, err error) {
	err = filepath.WalkDir(kmcc.KubeMCDir, func(path string, d fs.DirEntry, err error) error {

		if kmcc.KubeMCExt[filepath.Ext(path)] {
			kubeConfigList = append(kubeConfigList, path)
		}
		if err != nil {
			log.Fatal(err)
		}
		return nil
	},
	)
	return
}

func ReadKubeConfig(filename string) (k api.Config) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &k)
	if err != nil {
		log.Fatal(err)
	}
	return k
}

func MergeKubeConfigs(kubeconfigs []api.Config) (kubeconfig api.Config) {
	for _, k := range kubeconfigs {
		kubeconfig.Clusters = append(kubeconfig.Clusters, k.Clusters...)
		kubeconfig.Contexts = append(kubeconfig.Contexts, k.Contexts...)
		kubeconfig.AuthInfos = append(kubeconfig.AuthInfos, k.AuthInfos...)
	}
	return
}

func writeKubeConfig(filepath string, out []byte) error {
	err := ioutil.WriteFile(filepath, out, 0600)
	return err
}

func MergeKubeMC(kmcc KubeMCConfig) {
	log.Printf("kubemc: File Change in %v, trigger kubeconfig merge to %v.", kmcc.KubeMCDir, kmcc.KubeConfigPath)

	kubeConfigList, err := GetKubeConfigList(kmcc)

	if err != nil {
		log.Fatal(err)
	}

	var kubeConfigs []api.Config

	for _, file := range kubeConfigList {
		kubeConfigs = append(kubeConfigs, ReadKubeConfig(file))
	}
	kubeconfig := MergeKubeConfigs(kubeConfigs)
	yamlOut, err := yaml.Marshal(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	err = writeKubeConfig(kmcc.KubeConfigPath, yamlOut)
	if err != nil {
		log.Fatal(err)
	}
}
