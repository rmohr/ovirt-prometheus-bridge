package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
)

type Targets struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

type Hosts struct {
	Host []Host
}

type Host struct {
	Address string
	Cluster Cluster
}

type Cluster struct {
	Id string
}

func main() {
	target := flag.String("output", "engine-hosts.json", "target for the configuration file")
	engineURL := flag.String("engine-url", "https://localhost:8443", "Where to find the engine")
	engineUser := flag.String("engine-user", "admin@internal", "User")
	enginePassword := flag.String("engine-password", "engine", "Password")
	noVerify := flag.Bool("no-verify", false, "Don't verify the engine certificate")
	flag.Parse()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: *noVerify,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", *engineURL+"/ovirt-engine/api/hosts", nil)
	check(err)
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(*engineUser, *enginePassword)
	res, err := client.Do(req)
	check(err)
	hosts, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	check(err)
	writeTargets(*target, MapToTarget(ParseJson(hosts)))
}

func ParseJson(data []byte) *Hosts {
	hosts := new(Hosts)
	err := json.Unmarshal(data, hosts)
	check(err)
	return hosts
}

func MapToTarget(hosts *Hosts) []*Targets {
	targetMap := make(map[string]*Targets)
	var targets []*Targets
	for _, host := range hosts.Host {
		if value, ok := targetMap[host.Cluster.Id]; ok {
			value.Targets = append(value.Targets, host.Address)
		} else {
			targetMap[host.Cluster.Id] = &Targets{
				Labels:  map[string]string{"cluster": host.Cluster.Id},
				Targets: []string{host.Address}}
			targets = append(targets, targetMap[host.Cluster.Id])
		}
	}
	return targets
}

func writeTargets(fileName string, targets []*Targets) {
	data, _ := json.MarshalIndent(targets, "", "  ")
	data = append(data, '\n')
	err := ioutil.WriteFile(fileName, data, 0644)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
