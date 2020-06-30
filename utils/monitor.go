package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"strings"
	"time"
)

var (
	//HTTPReqDuration metric:http_request_duration_seconds
	HTTPReqDuration *prometheus.HistogramVec
	//HTTPReqTotal metric:http_request_total
	HTTPReqTotal *prometheus.CounterVec
	//HTTPReqErrTotal metric:http_request_error_total
	HTTPReqErrTotal *prometheus.CounterVec
)

// MonitorInit :called by main.go only once when system init
func MonitorInit() {
	fmt.Println("monitor utils start to init...")
	// Use Histogram to monitor the interface request duration.
	HTTPReqDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "fts_http_request_duration_seconds",
		Help:    "The HTTP request latencies in seconds.",
		Buckets: nil,
	}, []string{"method", "path"}) // labels

	// Use Counter to monitor the interface request total number
	HTTPReqTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "fts_http_requests_total",
		Help: "Total number of HTTP requests made.",
	}, []string{"method", "path", "status"}) // labels

	// Use Counter to monitor the interface request error total number
	HTTPReqErrTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "fts_http_requests_error_total",
		Help: "Total number of HTTP error status code made.",
	}, []string{"method", "path", "status"}) // labels

	prometheus.MustRegister(
		HTTPReqDuration,
		HTTPReqTotal,
		HTTPReqErrTotal,
	)
}

func parsePath(path string) string {
	fmt.Println("The original path is:")
	fmt.Println(path)
	itemList := strings.Split(path, "/")
	if len(itemList) >= 6 {
		return strings.Join(itemList[0:6], "/")
	}
	return path
}

//Metric metric middleware
func Metric() gin.HandlerFunc {
	return func(c *gin.Context) {
		tBegin := time.Now()
		c.Next()

		duration := float64(time.Since(tBegin)) / float64(time.Second)

		path := parsePath(c.Request.URL.Path)
		fmt.Println("After parsePath():")
		fmt.Println(path)
		// Counter Inc()
		HTTPReqTotal.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   path,
			"status": strconv.Itoa(c.Writer.Status()),
		}).Inc()

		// Record the current req duration
		HTTPReqDuration.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   path,
		}).Observe(duration)
	}
}
