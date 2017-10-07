package raygun

import (
	"reflect"
)

const RaygunEndpoint = "https://api.raygun.io"
const ClientName = "Raygun Go"
const ClientVersion = "0.1"
const ClientURL = "http://github.com/tomersimis/raygun"

var PackageName = reflect.TypeOf(NoopCollector{}).PkgPath()

var GlobalCollector Collector = &NoopCollector{}

func SetGlobalCollector(collector Collector) {
	GlobalCollector = collector
}

func CaptureError(err error, opts ...CaptureOption) {
	GlobalCollector.CaptureError(err, opts...)
}

func CapturePanic() {
	if rec := recover(); rec != nil {
		if err, ok := rec.(error); ok {
			GlobalCollector.CaptureError(err)
		} else {
			GlobalCollector.CaptureMessage(rec.(string))
		}
	}
}

func CaptureMessage(msg string, opts ...CaptureOption) {
	GlobalCollector.CaptureMessage(msg, opts...)
}

func Capture(ray Ray) {
	GlobalCollector.Capture(ray)
}

func Wait() {
	switch c := GlobalCollector.(type) {
	case *RaygunCollector:
		c.Wait()
	}
}
