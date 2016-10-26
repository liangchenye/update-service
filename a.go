package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Con struct {
	Login         string
	Repos_url     string
	Contributions int
}

var repos = []string{"docker/bender", "docker/birthdaysite", "docker/code-of-conduct", "docker/community", "docker/communitytools-image2docker-win", "docker/compose", "docker/containerd", "docker/datakit", "docker/dceu_tutorials", "docker/dctx", "docker/dcus-hol-2016", "docker/distribution", "docker/distribution-library-image", "docker/dnsserver", "docker/docker", "docker/docker-bb", "docker/docker-bench-security", "docker/docker-birthday-3", "docker/docker-credential-helpers", "docker/docker-e2e", "docker/docker-machine-driver-ci-test", "docker/docker-network", "docker/docker-py", "docker/docker-registry", "docker/docker-status",
	"docker/docker-tutorial", "docker/dockerhub.io", "docker/dockercloud-agent", "docker/dockercloud-authorizedkeys", "docker/dockercloud-cli", "docker/bender", "docker/birthdaysite", "docker/code-of-conduct", "docker/community", "docker/communitytools-image2docker-win", "docker/compose", "docker/containerd", "docker/datakit", "docker/dceu_tutorials", "docker/dctx",
	"docker/dcus-hol-2016", "docker/distribution", "docker/distribution-library-image",
	"docker/dnsserver",
	"docker/docker",
	"docker/docker-bb",
	"docker/docker-bench-security",
	"docker/docker-birthday-3",
	"docker/docker-credential-helpers",
	"docker/docker-e2e",
	"docker/docker-machine-driver-ci-test",
	"docker/docker-network",
	"docker/docker-py",
	"docker/docker-registry",
	"docker/docker-status",
	"docker/docker-tutorial",
	"docker/dockerhub.io",
	"docker/dockercloud-agent",
	"docker/dockercloud-authorizedkeys",
	"docker/dockercloud-cli",
	"docker/dockercloud-events",
	"docker/dockercloud-haproxy",
	"docker/dockercloud-hello-world",
	"docker/dockercloud-network-daemon",
	"docker/dockercloud-node",
	"docker/dockercloud-quickstart-go",
	"docker/dockercloud-quickstart-python",
	"docker/dockercraft",
	"docker/dockerlite",
	"docker/engine-api",
	"docker/etcd",
	"docker/example-voting-app",
	"docker/for-mac",
	"docker/for-win",
	"docker/global-hack-day-3",
	"docker/go",
	"docker/go-connections",
	"docker/go-dockercloud",
	"docker/go-events",
	"docker/go-healthcheck",
	"docker/go-metrics",
	"docker/go-p9p",
	"docker/go-plugins-helpers",
	"docker/go-redis-server",
	"docker/go-units",
	"docker/goamz",
	"docker/golem",
	"docker/gordon",
	"docker/gordon-bot",
	"docker/homebrew-core",
	"docker/hub-feedback",
	"docker/hugo",
	"docker/hyperkit",
	"docker/infrakit",
	"docker/infrakit.aws",
	"docker/irc-minutes",
	"docker/jenkins-pipeline-scripts",
	"docker/jira-test",
	"docker/kitematic",
	"docker/labs",
	"docker/leadership",
	"docker/leeroy",
	"docker/libchan",
	"docker/libcompose",
	"docker/libcontainer",
	"docker/libkv",
	"docker/libnetwork",
	"docker/libtrust",
	"docker/linkcheck",
	"docker/machine",
	"docker/markdownlint",
	"docker/migrator",
	"docker/notary",
	"docker/notary-official-images",
	"docker/notary-server-image",
	"docker/notary-signer-image",
	"docker/opensource",
	"docker/openstack-docker",
	"docker/openstack-heat-docker",
	"docker/orchestration-workshop",
	"docker/pulpo",
	"docker/python-dockercloud",
	"docker/runc",
	"docker/spdystream",
	"docker/swarm",
	"docker/swarm-frontends",
	"docker/swarm-library-image",
	"docker/swarm-microservice-demo-v1",
	"docker/swarmkit",
	"docker/toolbox",
	"docker/ucp_lab",
	"docker/v1.10-migrator",
	"docker/vpnkit",
}

func pullRepo(repo string) {
	for i := 0; ; i++ {
		rawurl := fmt.Sprintf("https://api.github.com/repos/%s/contributors?page=%d", repo, i)
		fmt.Println(rawurl)
		res, err := SendHttpRequest("GET", rawurl, nil, nil)
		if err != nil {
			return
		}
		var cons []Con
		body, err := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &cons)
		if len(cons) == 0 {
			break
		}
		pubCons = append(pubCons, cons...)
	}
}

var pubCons []Con

func main() {
	repo := "docker/docker"
	pullRepo(repo)

	fmt.Println(pubCons)
}

func SendHttpRequest(method, rawurl string, body io.Reader, header map[string]string) (*http.Response, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return &http.Response{}, err
	}

	var client *http.Client
	switch url.Scheme {
	case "":
		fallthrough
	case "https":
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	case "http":
		client = &http.Client{}
	default:
		return &http.Response{}, fmt.Errorf("bad url schema: %v", url.Scheme)
	}

	fmt.Println("aa")
	req, err := http.NewRequest(method, url.String(), body)
	fmt.Println("aa1")
	if err != nil {
		return &http.Response{}, err
	}
	fmt.Println("aa2")
	req.URL.RawQuery = req.URL.Query().Encode()
	for k, v := range header {
		req.Header.Set(k, v)
	}
	return client.Do(req)
}
