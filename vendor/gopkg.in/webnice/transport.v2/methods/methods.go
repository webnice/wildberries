package methods

import (
	"strings"
)

// New Function create new implementation of interface
func New() Interface {
	return new(impl)
}

// Options Return HTTP method OPTIONS
func (m *impl) Options() Value {
	return &methodType{optionsMethod}
}

// Get Return HTTP method GET
func (m *impl) Get() Value {
	return &methodType{getMethod}
}

// Head Return HTTP method HEAD
func (m *impl) Head() Value {
	return &methodType{headMethod}
}

// Post Return HTTP method POST
func (m *impl) Post() Value {
	return &methodType{postMethod}
}

// Put Return HTTP method PUT
func (m *impl) Put() Value {
	return &methodType{putMethod}
}

// Patch Return HTTP method PATCH
func (m *impl) Patch() Value {
	return &methodType{patchMethod}
}

// Delete Return HTTP method DELETE
func (m *impl) Delete() Value {
	return &methodType{deleteMethod}
}

// Trace Return HTTP method TRACE
func (m *impl) Trace() Value {
	return &methodType{traceMethod}
}

// Connect Return HTTP method CONNECT
func (m *impl) Connect() Value {
	return &methodType{connectMethod}
}

// Parse string and return method interface
func (m *impl) Parse(inp string) Value {
	var (
		tmp string
		key Type
		ret *methodType
	)

	tmp = strings.ToUpper(inp)
	for key = range maps {
		if maps[key] == tmp {
			ret = &methodType{key}
			break
		}
	}

	return ret
}

// Parse string and return interface
func Parse(inp string) Value {
	return New().Parse(inp)
}
