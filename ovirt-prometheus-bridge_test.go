package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestParsingHosts(t *testing.T) {
	hosts := loadJson()
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

	hosts := loadJson()
	result := MapToTarget(hosts)
	if len(result) != 2 {
		t.Error("Expected 2, got ", len(hosts.Host))
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
	writeTargets("generated-targets.json", MapToTarget(loadJson()))
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
	writeTargets("generated-targets.json", MapToTarget(&Hosts{}))
	defer os.Remove("generated-targets.json")
	generated, err := ioutil.ReadFile("generated-targets.json")
	check(err)
	if !bytes.Equal([]byte("[]"), generated) {
		t.Errorf("Expected '%s', got '%s'", "[]", string(generated))
	}
}

func loadJson() *Hosts {
	data, err := ioutil.ReadFile("hosts.json")
	check(err)
	hosts := ParseJson(data)
	return hosts
}
