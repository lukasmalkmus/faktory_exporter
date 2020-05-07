package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	faktory "github.com/contribsys/faktory/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	// namespace for all metrics of this exporter.
	namespace = "faktory"
)

var (
	faktoryURL       = kingpin.Flag("faktory.url", "URL of the faktory instance").Default("tcp://localhost:7419").String()
	webListenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface").Default(":9386").String()
	webMetricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics").Default("/metrics").String()
)

// landingPage contains the HTML served at '/'.
var landingPage = `<html>
	<head>
		<title>Faktory Exporter</title>
	</head>
	<body>
		<h1>Faktory Exporter</h1>
		<p>
		<a href=` + *webMetricsPath + `>Metrics</a>
		</p>
	</body>
</html>`

type Tasks struct {
	retries Retries
}

type Retries struct {
	cycles   prometheus.Gauge
	enqueued prometheus.Gauge
	size     prometheus.Gauge
}

// Exporter collects stats from a Faktory instance by issuing the "INFO" command
// and exports them using the prometheus client library.
type Exporter struct {
	mutex sync.RWMutex

	client *faktory.Client

	// Basic exporter metrics.
	up, scrapeDuration          prometheus.Gauge
	totalScrapes, failedScrapes prometheus.Counter

	// Faktory metrics.
	commandCount prometheus.Counter
	connections  prometheus.Gauge
	jobs         *prometheus.CounterVec
	totalQueues  prometheus.Gauge
	tasks        Tasks
	queues       *prometheus.GaugeVec
}

// New creates and returns a new, initialized Faktory Exporter.
func New(faktoryURL string) (*Exporter, error) {
	// Parse provided URL.
	uri, err := url.Parse(faktoryURL)
	if err != nil {
		return nil, err
	}

	//  Create server struct and initialize it with info retrieved from
	// parsed URL.
	srv := faktory.DefaultServer()
	srv.Network = uri.Scheme
	srv.Address = fmt.Sprintf("%s:%s", uri.Hostname(), uri.Port())
	pwd := ""
	if uri.User != nil {
		pwd, _ = uri.User.Password()
	}

	// Setup connection to Faktory instance.
	client, err := faktory.Dial(srv, pwd)
	if err != nil {
		return nil, err
	}

	e := &Exporter{
		// Faktory client.
		client: client,

		// Basic exporter metrics.
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Was the last scrape of the Faktory instance successful?",
		}),
		scrapeDuration: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "exporter",
			Name:      "scrape_duration_seconds",
			Help:      "Duration of the scrape of metrics from the Faktory instance.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "exporter",
			Name:      "scrapes_total",
			Help:      "Total Faktory scrapes.",
		}),
		failedScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "exporter",
			Name:      "scrape_failures_total",
			Help:      "Total amount of scrape failures.",
		}),

		// Faktory metrics.
		commandCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "server",
			Name:      "command_count",
			Help:      "Number of commands which have been issued to the server.",
		}),
		connections: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "server",
			Name:      "connections",
			Help:      "Number of currently connected clients.",
		}),
		jobs: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "jobs",
			Name:      "total",
			Help:      "Total amount of jobs.",
		}, []string{"status"}),
		totalQueues: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "queues",
			Name:      "total",
			Help:      "Total amount of queues.",
		}),
		queues: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "queue",
			Name:      "jobs",
			Help:      "Number of jobs in every queue.",
		}, []string{"queue"}),
		tasks: Tasks{retries: Retries{
			enqueued: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "tasks_retries",
				Name:      "enqueued",
				Help:      "Task retries enqueued.",
			}),
			size: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "tasks_retries",
				Name:      "size",
				Help:      "Task retries size.",
			}),
		}},
	}

	return e, nil
}

// Describe all the metrics collected by the Faktory exporter.
// Implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.up.Describe(ch)
	e.scrapeDuration.Describe(ch)
	e.failedScrapes.Describe(ch)
	e.totalScrapes.Describe(ch)
	e.commandCount.Describe(ch)
	e.connections.Describe(ch)
	e.jobs.Describe(ch)
	e.tasks.retries.enqueued.Describe(ch)
	e.tasks.retries.size.Describe(ch)
	e.totalQueues.Describe(ch)
	e.queues.Describe(ch)
}

// Collect the stats from the configured Faktory instance and deliver them as
// Prometheus metrics.
// Implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Protect metrics from concurrent collects.
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Reset metrics.
	e.reset()

	// Scrape metrics from Faktory instance.
	if err := e.scrape(); err != nil {
		log.Error(err)
	}

	// Collect metrics.
	e.up.Collect(ch)
	e.scrapeDuration.Collect(ch)
	e.failedScrapes.Collect(ch)
	e.totalScrapes.Collect(ch)
	e.commandCount.Collect(ch)
	e.connections.Collect(ch)
	e.jobs.Collect(ch)
	e.tasks.retries.enqueued.Collect(ch)
	e.tasks.retries.size.Collect(ch)
	e.totalQueues.Collect(ch)
	e.queues.Collect(ch)
}

// reset the vector metrics.
func (e *Exporter) reset() {
	e.jobs.Reset()
}

// scrape issues the "INFO" command to the Faktory instance, fetches the
// appropriate metrics and meassures the scrapes duration.
func (e *Exporter) scrape() (err error) {
	// Meassure scrape duration, increase total scrapes and evaluate if the
	// scrape was successful and set up/down and failes scrapes accordingly.
	defer func(begun time.Time) {
		e.scrapeDuration.Set(time.Since(begun).Seconds())
		e.totalScrapes.Inc()
		if err != nil {
			e.up.Set(0)
			e.failedScrapes.Inc()
		} else {
			e.up.Set(1)
		}
	}(time.Now())

	// Fetch info from configured Faktory instance.
	info, err := e.client.Info()
	if err != nil {
		return err
	}
	b, _ := json.MarshalIndent(info, "", "  ")
	println(string(b))

	// Get server stats.
	server := info["server"].(map[string]interface{})
	commandCount, ok := server["command_count"].(float64)
	if !ok {
		return fmt.Errorf("error getting command count")
	}
	connections, ok := server["connections"].(float64)
	if !ok {
		return fmt.Errorf("error getting connections")
	}

	// Get job and queue stats.
	faktory := info["faktory"].(map[string]interface{})
	enqueued, ok := faktory["total_enqueued"].(float64)
	if !ok {
		return fmt.Errorf("error getting enqueued jobs")
	}
	failures, ok := faktory["total_failures"].(float64)
	if !ok {
		return fmt.Errorf("error getting failed jobs")
	}
	processed, ok := faktory["total_processed"].(float64)
	if !ok {
		return fmt.Errorf("error getting processed jobs")
	}
	totalQueues, ok := faktory["total_queues"].(float64)
	if !ok {
		return fmt.Errorf("error getting queues")
	}
	tasks, ok := faktory["tasks"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("error getting tasks")
	}
	retries, ok := tasks["Retries"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("error getting retries")
	}

	// Set all metrics.
	e.commandCount.Set(commandCount)
	e.connections.Set(connections)
	e.jobs.WithLabelValues("enqueued").Set(enqueued)
	e.jobs.WithLabelValues("failure").Set(failures)
	e.jobs.WithLabelValues("processed").Set(processed)
	e.tasks.retries.enqueued.Set(retries["enqueued"].(float64))
	e.tasks.retries.size.Set(retries["size"].(float64))
	e.totalQueues.Set(totalQueues)

	queues, ok := faktory["queues"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("error getting queue counts")
	}

	for queue, jobs := range queues {
		e.queues.WithLabelValues(queue).Set(jobs.(float64))
	}
	return nil
}

func main() {
	os.Exit(Main())
}

// Main manages the complete application lifecycle, from startup to shutdown.
func Main() int {
	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("faktory_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	// Print build context and version.
	log.Info("Starting faktory_exporter", version.Info())
	log.Info("Build context", version.BuildContext())

	// Create a new Faktory exporter. Exit if an error is returned.
	exporter, err := New(*faktoryURL)
	if err != nil {
		log.Error(err)
		return 0
	}

	// Register Faktory and the collector for version information.
	// Unregister Go and Process collector which are registered by default.
	prometheus.MustRegister(exporter)
	prometheus.MustRegister(version.NewCollector("faktory_exporter"))
	prometheus.Unregister(prometheus.NewGoCollector())
	prometheus.Unregister(prometheus.NewProcessCollector(os.Getpid(), ""))

	// Setup router and handlers.
	r := http.NewServeMux()
	metricsHandler := promhttp.HandlerFor(prometheus.DefaultGatherer,
		promhttp.HandlerOpts{ErrorLog: log.NewErrorLogger()})
	// TODO: InstrumentHandler is depracted. Additional tools will be available
	// soon in the promhttp package.
	//r.Handle(*webMetricsPath, prometheus.InstrumentHandler("faktory_exporter", metricsHandler))
	r.Handle(*webMetricsPath, metricsHandler)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(landingPage))
	})

	// Setup webserver.
	srv := &http.Server{
		Addr:           *webListenAddress,
		Handler:        r,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ErrorLog:       log.NewErrorLogger(),
	}

	// Listen for termination signals.
	term := make(chan os.Signal, 1)
	webErr := make(chan error)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	// Run webserver in a separate goroutine.
	log.Infoln("Listening on", *webListenAddress)
	go func() { webErr <- srv.ListenAndServe() }()

	// Wait for a termination signal and shut down gracefully, but wait no
	// longer than 5 seconds before halting.
	select {
	case <-term:
		log.Warn("Received SIGTERM, exiting gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Errorln("Error shutting down http server:", err)
		}
	case err := <-webErr:
		log.Errorln("Error starting http server, exiting gracefully:", err)
	}

	log.Info("See you next time!")

	return 0
}
