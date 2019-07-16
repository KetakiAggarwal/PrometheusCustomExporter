package main

import (
	"net/http"
        "github.com/prometheus/common/log"
    	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
        "gopkg.in/alecthomas/kingpin.v2"
   
)

func main() {
        var (
            listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":3874").String()
            configFile5min = kingpin.Flag("config-file-5min", "Config File. Use the .yml extension to use yaml format").Default("services5min.yml").String()
	    configFile10min = kingpin.Flag("config-file-10min", "Config File. Use the .yml extension to use yaml format").Default("services10min.yml").String()
	    configFile15min = kingpin.Flag("config-file-15min", "Config File. Use the .yml extension to use yaml format").Default("services15min.yml").String()
	    configFile30min = kingpin.Flag("config-file-30min", "Config File. Use the .yml extension to use yaml format").Default("services30min.yml").String()
	    configFile1hr = kingpin.Flag("config-file-1hr", "Config File. Use the .yml extension to use yaml format").Default("services1hr.yml").String()
	    configFile12hr = kingpin.Flag("config-file-12hr", "Config File. Use the .yml extension to use yaml format").Default("services12hr.yml").String()
	    configFile24hr = kingpin.Flag("config-file-24hr", "Config File. Use the .yml extension to use yaml format").Default("services24hr.yml").String()


        )

        log.AddFlags(kingpin.CommandLine)
        kingpin.Version("0.1")
        kingpin.HelpFlag.Short('h')
        kingpin.Parse()

	c5min := readConfig(configFile5min)
	c10min := readConfig10min(configFile10min)
	c24hr := readConfig24hr(configFile24hr)
	c15min := readConfig15min(configFile15min)
	c30min := readConfig30min(configFile30min)
	c1hr := readConfig1hr(configFile1hr)
	c12hr := readConfig12hr(configFile12hr)

	//Create a new instance of the custom collector and
	//register it with the prometheus client.
	r5min := prometheus.NewRegistry()
	foo5min := newSvcCollector(c5min)
	r5min.MustRegister(foo5min)
	handler5min := promhttp.HandlerFor(r5min, promhttp.HandlerOpts{})

	r10min := prometheus.NewRegistry()
	foo10min := newSvcCollector10min(c10min)
        r10min.MustRegister(foo10min)
        handler10min := promhttp.HandlerFor(r10min, promhttp.HandlerOpts{})
	
	r15min := prometheus.NewRegistry()
        foo15min := newSvcCollector15min(c15min)
        r15min.MustRegister(foo15min)
        handler15min := promhttp.HandlerFor(r15min, promhttp.HandlerOpts{})

	r30min := prometheus.NewRegistry()
        foo30min := newSvcCollector30min(c30min)
        r30min.MustRegister(foo30min)
        handler30min := promhttp.HandlerFor(r30min, promhttp.HandlerOpts{})

	r1hr := prometheus.NewRegistry()
        foo1hr := newSvcCollector1hr(c1hr)
        r1hr.MustRegister(foo1hr)
        handler1hr := promhttp.HandlerFor(r1hr, promhttp.HandlerOpts{})

	r12hr := prometheus.NewRegistry()
        foo12hr := newSvcCollector12hr(c12hr)
        r12hr.MustRegister(foo12hr)
        handler12hr := promhttp.HandlerFor(r12hr, promhttp.HandlerOpts{})
	
	r24hr := prometheus.NewRegistry()
        foo24hr := newSvcCollector24hr(c24hr)
        r24hr.MustRegister(foo24hr)
        handler24hr := promhttp.HandlerFor(r24hr, promhttp.HandlerOpts{})

	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.
	http.Handle("/metrics5min", handler5min)
	http.Handle("/metrics10min", handler10min)
	http.Handle("/metrics15min", handler15min)
	http.Handle("/metrics30min", handler30min)
	http.Handle("/metrics1hr", handler1hr)
	http.Handle("/metrics12hr", handler12hr)
	http.Handle("/metrics24hr", handler24hr)
	log.Info("Beginning to serve on port " + *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}
