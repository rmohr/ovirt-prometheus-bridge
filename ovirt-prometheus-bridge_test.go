package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestParsingHosts(t *testing.T) {
	hosts := loadHosts()
	if len(hosts.Host) != 3 {
		t.Error("Expected 3, got ", len(hosts.Host))
	}
	if hosts.Host[0].Address != "host1.com" {
		t.Errorf("Expected address 'host1.com', got '%s'", hosts.Host[0].Address)
	}
	if hosts.Host[0].Cluster.Id != "cluster1_id" {
		t.Errorf("Expected cluster Id 'cluster1_id', got '%s'", hosts.Host[0].Cluster.Id)
	}
	if hosts.Host[1].Address != "host2.com" {
		t.Errorf("Expected address 'host2.com', got '%s'", hosts.Host[1].Address)
	}
}

func TestMapToTarget(t *testing.T) {
	samples := map[string][]string{
		"cluster1_id": []string{"host1.com"},
		"cluster2_id": []string{"host2.com", "host3.com"}}

	hostsChan := make(chan *Hosts, 1)
	hosts := loadHosts()
	hostsChan <- hosts
	close(hostsChan)
	result := <-MapToTarget(hostsChan)
	if len(result) != 2 {
		t.Error("Expected 2, got ", len(result))
	}
	for _, value := range result {
		for i, sample := range samples[value.Labels["cluster"]] {
			if sample != value.Targets[i] {
				t.Errorf("Expected hosts '%s', got '%s'", sample, value.Targets[i])
			}
		}
	}
}

func TestWriteJson(t *testing.T) {
	done := writeTargets("generated-targets.json", MapToTarget(loadHostsIntoChan()))
	<-done
	defer os.Remove("generated-targets.json")
	original, err := ioutil.ReadFile("targets.json")
	check(err)
	generated, err := ioutil.ReadFile("generated-targets.json")
	check(err)
	if !bytes.Equal(original, generated) {
		t.Errorf("Expected '%s', got '%s'", string(original), string(generated))
	}
}

func TestNoTargets(t *testing.T) {
	done := writeTargets("generated-targets.json", MapToTarget(noHosts()))
	<-done
	defer os.Remove("generated-targets.json")
	generated, err := ioutil.ReadFile("generated-targets.json")
	check(err)
	if !bytes.Equal([]byte("[]"), generated) {
		t.Errorf("Expected '%s', got '%s'", "[]", string(generated))
	}
}

func loadHosts() *Hosts {
	data, err := ioutil.ReadFile("hosts.json")
	check(err)
	dataChan := make(chan []byte, 1)
	dataChan <- data
	close(dataChan)
	hosts := ParseJson(dataChan)
	return <-hosts
}

func noHosts() chan *Hosts {
	hosts := make(chan *Hosts, 1)
	hosts <- &Hosts{}
	close(hosts)
	return hosts
}

func loadHostsIntoChan() chan *Hosts {
	hostsChan := make(chan *Hosts, 1)
	hostsChan <- loadHosts()
	close(hostsChan)
	return hostsChan
}
