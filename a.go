package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Con struct {
	Login         string
	MonitRepo     string
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

func pullRepo(repo string) []Con {
	var repoCon []Con
	for i := 1; ; i++ {
		rawurl := fmt.Sprintf("https://api.github.com/repos/%s/contributors?page=%d", repo, i)
		fmt.Println(rawurl)
		res, err := SendHttpRequest("GET", rawurl, nil, nil)
		if err != nil {
			return repoCon
		}
		var cons []Con
		body, err := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &cons)
		if len(cons) == 0 {
			break
		}
		for i, _ := range cons {
			cons[i].MonitRepo = repo
		}
		repoCon = append(repoCon, cons...)
	}

	return repoCon
}

func main() {
	var pubCons []Con
	for _, repo := range repos {
		repoCon := pullRepo(repo)
		repoData, _ := json.Marshal(repoCon)
		tmpFile := strings.Replace(repo, "/", "#", 2)
		ioutil.WriteFile("githubdata/"+tmpFile, repoData, 0644)

		pubCons = append(pubCons, repoCon...)
	}

	data, _ := json.Marshal(pubCons)
	ioutil.WriteFile("githubdata/alldata", data, 0644)

	contrib := make(map[string]int)

	for _, c := range pubCons {
		val, ok := contrib[c.Login]
		if !ok {
			contrib[c.Login] = c.Contributions
		} else {
			contrib[c.Login] = c.Contributions + val
		}
	}

	conData, _ := json.Marshal(contrib)
	ioutil.WriteFile("githubdata/mapdata", conData, 0644)
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

	req, err := http.NewRequest(method, url.String(), body)
	req.SetBasicAuth("initlove", "david840318")
	if err != nil {
		return &http.Response{}, err
	}
	req.URL.RawQuery = req.URL.Query().Encode()
	for k, v := range header {
		req.Header.Set(k, v)
	}
	return client.Do(req)
}
