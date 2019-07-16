package main

import (
	"bufio"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
//	"time"
	"sync"
//	"encoding/json"
)

type svcTarget30min struct {
	svcName              string
	svcCmd               string
	svcKey		     string
	svcData		     string
	svcMetric[]*prometheus.Desc
}

type svcCollector30min struct {
	svcTargets []svcTarget
}

type yml_conf30min struct {
	Monitor struct {
		Services map[string][]struct {
			Bin     string
			Key     string
			Enabled string
			Data    string
		}
	}
}

type conf30min struct {
	Monitor struct {
		Services struct {
			Uq_name []struct {
				Bin string
			}
		}
	}
}

func (c *conf30min) getConf() *conf30min {

	yamlFile, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}


func readConfig30min(configFile *string) *svcCollector30min {
	fmt.Println("*****readConfig30min")

	ext := filepath.Ext(*configFile)

	fmt.Printf("ext: %s\n", ext)

	c := svcCollector30min{}
	if ext == ".yml" || ext == ".yaml" {
		fmt.Printf("Parsing config as yaml file\n")
		yamlFile, err := ioutil.ReadFile(*configFile)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
                y := yml_conf{}

		err = yaml.Unmarshal(yamlFile, &y)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}

		targets := make([]svcTarget, 0, 100)
		for k, _ := range y.Monitor.Services {
			t := svcTarget{}
                        fmt.Printf("   Loading check: %-15s    - %s\n", k, y.Monitor.Services[k][0].Bin)
			t.svcName = k
			t.svcCmd = y.Monitor.Services[k][0].Bin
			t.svcKey=y.Monitor.Services[k][0].Key
			t.svcData=y.Monitor.Services[k][0].Data
			targets = append(targets, t)
			fmt.Printf("   Loading check: %-15s    - %s\n", k, y.Monitor.Services[k][0].Key)
		}
		c.svcTargets = targets

	} else {
		fmt.Printf("Parsing config as simple property file (for testing use only)!\n")
		f, err := os.Open(*configFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		fmt.Printf("Parsing config as txt file\n")
		scanner := bufio.NewScanner(f)
		targets := make([]svcTarget, 0, 100)
		for scanner.Scan() {
			line := scanner.Text()
			a := strings.Fields(line)
			fmt.Printf("Service check --- %s: %s\n", a[0], a[1])

			if a[0] != "" {
				t := svcTarget{}
				t.svcName = a[0]
				t.svcCmd = a[1]
				targets = append(targets, t)
			}
		}

		c.svcTargets = targets
	}
	return &c

}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func newSvcCollector30min(c *svcCollector30min) *svcCollector30min {

	fmt.Println("*****newSvcCollector30min")
	for i, t := range c.svcTargets {
	//	mname := strings.ToLower(t.svcName) + "_avail"
	//	lname := "fus_" + strings.ToLower(t.svcName) + "_avail_response_latency_second"
		keyStr := strings.Split(t.svcKey,",")
		dataStr := strings.Split(t.svcData,",")
		for _, data := range dataStr {	
		// You must not use "t" variable to update the metric as it is pass by reference
		// Use the actual index is required.
			mname := strings.ToLower(t.svcName) + "_" +data
			if (t.svcKey != ""){
                        c.svcTargets[i].svcMetric = append(c.svcTargets[i].svcMetric, prometheus.NewDesc(mname, "Show "+strings.ToUpper(t.svcName)+"_"+ data,keyStr , nil))} else{ c.svcTargets[i].svcMetric = append(c.svcTargets[i].svcMetric, prometheus.NewDesc(mname, "Show "+strings.ToUpper(t.svcName)+"_"+ data,nil , nil))
                	}
		}
		fmt.Printf("  Index: %d -> %s\n", i, c.svcTargets[i])
	}
	return c
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (c *svcCollector30min) Describe(ch chan<- *prometheus.Desc) {
	fmt.Println("*****Describe each collector")
	for i, t := range c.svcTargets {
		fmt.Printf("  Index: %d -> %s\n", i, t)
		for _, m := range t.svcMetric {
			ch <- m
		}
	//	ch <- t.svcRespLatencyMetric
	}
}
//Collect implements required collect function for all promehteus collectors
func (c *svcCollector30min) Collect(ch chan<- prometheus.Metric) {
       // channel := make(chan<- prometheus.Metric)

        fmt.Println("*****Collect each collector 12 hr")
        fmt.Println("Do collections every 12 hour")
        wg := new(sync.WaitGroup)
        wg.Add(len(c.svcTargets))

        for _, t := range c.svcTargets {
                fmt.Printf("   collecting %s: %s n", t.svcName, t.svcCmd)
                go exe_cmd_custom(t, wg,ch)
        }
        wg.Wait()
}

