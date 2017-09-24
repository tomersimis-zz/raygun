package raygun

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Collector struct {
	AppName   string
	ApiKey    string
	Workers   int
	QueueSize int
	Logger    *log.Logger

	queue  chan Ray
	client *http.Client
	wg     sync.WaitGroup
}

type CollectorConfig func(*Collector)

func Workers(n int) CollectorConfig {
	return func(c *Collector) {
		c.Workers = n
	}
}

func QueueSize(n int) CollectorConfig {
	return func(c *Collector) {
		c.QueueSize = n
	}
}

func Logger(logger *log.Logger) CollectorConfig {
	return func(c *Collector) {
		c.Logger = logger
	}
}

func NewCollector(appName, apiKey string, options ...CollectorConfig) *Collector {

	collector := &Collector{
		AppName:   appName,
		ApiKey:    apiKey,
		Workers:   1,
		QueueSize: 10000,
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true,
			},
			Timeout: 5 * time.Second,
		},
	}

	for _, f := range options {
		f(collector)
	}

	if collector.Logger == nil {
		collector.Logger = log.New(ioutil.Discard, "raygun", 0)
	}

	collector.queue = make(chan Ray, collector.QueueSize)
	collector.start()
	return collector
}

func (c *Collector) CaptureMessage(msg string) {
	c.queue <- NewRay(msg)
	c.wg.Add(1)
}

func (c *Collector) CaptureError(err error) {
	c.CaptureMessage(err.Error())
}

func (c *Collector) CapturePanic() func() {
	return func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				c.CaptureError(err)
			} else {
				c.CaptureMessage(rec.(string))
			}
		}
	}
}

func (c *Collector) Capture(ray Ray) {
	c.queue <- ray
	c.wg.Add(1)
}

func (c *Collector) start() {
	for i := 0; i < c.Workers; i++ {
		go func() {
			for {
				ray := <-c.queue

				json, err := json.Marshal(ray)
				if err != nil {
					c.Logger.Printf("raygun: failed to marshal raygun error: %s", err.Error())
				}

				req, err := http.NewRequest("POST", RaygunEndpoint+"/entries", bytes.NewBuffer(json))
				if err != nil {
					c.Logger.Printf("raygun: failed to create error request: %s", err.Error())
				}
				req.Header.Add("X-ApiKey", c.ApiKey)
				res, err := c.client.Do(req)
				if err != nil {
					c.Logger.Printf("raygun: request failed: %s", err.Error())
				}

				res.Body.Close()
				c.wg.Done()
			}
		}()
	}
}

func (c *Collector) Wait() {
	c.wg.Wait()
}
