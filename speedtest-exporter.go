package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type SpeedtestResult struct {
	Ping struct {
		Jitter  float64 `json:"jitter"`
		Latency float64 `json:"latency"`
	} `json:"ping"`
	Download struct {
		Bandwidth int64 `json:"bandwidth"`
	} `json:"download"`
	Upload struct {
		Bandwidth int64 `json:"bandwidth"`
	} `json:"upload"`
}

type metrics struct {
	Download prometheus.Gauge
	Upload   prometheus.Gauge
	Jitter   prometheus.Gauge
	Ping     prometheus.Gauge
}

func registerMetrics(registry *prometheus.Registry) *metrics {
	m := &metrics{
		Download: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "speedtest_download",
			Help: "Speedtest Download in Bytes per second (B/s)",
		}),
		Upload: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "speedtest_upload",
			Help: "Speedtest Upload in Bytes per second (B/s)",
		}),
		Jitter: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "speedtest_jitter",
			Help: "Speedtest Jitter in milliseconds (ms)",
		}),
		Ping: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "speedtest_ping",
			Help: "Speedtest Ping in millisecond (ms)",
		}),
	}
	registry.MustRegister(m.Download)
	registry.MustRegister(m.Upload)
	registry.MustRegister(m.Jitter)
	registry.MustRegister(m.Ping)
	return m
}

func runSpeedtest() (SpeedtestResult, error) {
	var result SpeedtestResult
	cmd := exec.Command("speedtest", "--accept-license", "--accept-gdpr", "-f", "json")
	log.Println("Starting speedtest")
	start := time.Now()
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing speedtest command:", err)
		return result, err
	}
	err = json.Unmarshal(output, &result)
	if err != nil {
		fmt.Println("Error parsing speedtest JSON:", err)
		return result, err
	}
	log.Println("Finished speedtest in", time.Since(start))
	return result, nil
}

func handleMetrics(m *metrics, reg *prometheus.Registry) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := runSpeedtest()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		m.Upload.Set(float64(result.Upload.Bandwidth))
		m.Download.Set(float64(result.Download.Bandwidth))
		m.Jitter.Set(result.Ping.Jitter)
		m.Ping.Set(result.Ping.Latency)
		w.Header().Set("Content-Type", "text/plain")
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9798"
	}

	reg := prometheus.NewRegistry()
	m := registerMetrics(reg)
	http.HandleFunc("/metrics", handleMetrics(m, reg))
	log.Println("Starting server on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
