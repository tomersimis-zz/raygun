package raygun

import (
	"reflect"
)

const RaygunEndpoint = "https://api.raygun.io"
const ClientName = "Raygun Go"
const ClientVersion = "0.1"
const ClientURL = "http://bitbucket.org/ubeedev/engage"

var PackageName = reflect.TypeOf(Collector{}).PkgPath()

var GlobalCollector *Collector

func SetGlobalCollector(collector *Collector) {
	GlobalCollector = collector
}

func CaptureError(err error) {
	GlobalCollector.CaptureError(err)
}

func CapturePanic() func() {
	return GlobalCollector.CapturePanic()
}

func CaptureMessage(msg string) {
	GlobalCollector.CaptureMessage(msg)
}

func Capture(ray Ray) {
	GlobalCollector.Capture(ray)
}

func Wait() {
	GlobalCollector.Wait()
}
