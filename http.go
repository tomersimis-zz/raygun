package raygun

import (
	"os"
	"time"
)

type Ray struct {
	OccuredOn string  `json:"occurredOn"`
	Details   Details `json:"details"`
}

type Details struct {
	MachineName    string      `json:"machineName"`
	Version        string      `json:"version"`
	Error          Error       `json:"error"`
	Tags           []string    `json:"tags"`
	UserCustomData interface{} `json:"userCustomData"`
	Request        Request     `json:"request"`
	User           User        `json:"user"`
	Context        Context     `json:"context"`
	Client         Client      `json:"client"`
}

type Error struct {
	Message    string     `json:"message"`
	StackTrace StackTrace `json:"stackTrace"`
}

type Request struct {
	HostName    string            `json:"hostName"`
	URL         string            `json:"url"`
	HTTPMethod  string            `json:"httpMethod"`
	IPAddress   string            `json:"ipAddress"`
	QueryString map[string]string `json:"queryString"`
	Form        map[string]string `json:"form"`
	Headers     map[string]string `json:"headers"`
}

type Client struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	ClientURL string `json:"clientUrl"`
}

type User struct {
	Identifier string `json:"identifier"`
}

type Context struct {
	Identifier string `json:"identifier"`
}

type StackTraceElement struct {
	LineNumber  int    `json:"lineNumber"`
	PackageName string `json:"className"`
	FileName    string `json:"fileName"`
	MethodName  string `json:"methodName"`
}

type StackTrace []StackTraceElement

func (s *StackTrace) AddEntry(lineNumber int, packageName, fileName, methodName string) {
	*s = append(*s, StackTraceElement{lineNumber, packageName, fileName, methodName})
}

func NewRay(msg string) Ray {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "not available"
	}

	return Ray{
		OccuredOn: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		Details: Details{
			MachineName: hostname,
			Error: Error{
				Message:    msg,
				StackTrace: GetCurrentStack(),
			},
		},
	}
}
