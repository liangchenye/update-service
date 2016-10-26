package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func main() {
	pubC := 0
	for _, repo := range repos {
		tmpFile := strings.Replace(repo, "/", "#", 2)
		repoData, _ := ioutil.ReadFile("githubdata/" + tmpFile)

		var cons []Con
		json.Unmarshal(repoData, &cons)
		for _, c := range cons {
			if c.Login == "crosbymichael" {
				fmt.Println(c.Contributions)
				pubC += c.Contributions
				break
			}
		}
	}
	fmt.Println(pubC)
}
