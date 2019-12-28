package methods

// Int constants as named type
const (
	optionsMethod Type = iota + 1 // 1 OPTIONS
	getMethod                     // 2 GET
	headMethod                    // 3 HEAD
	postMethod                    // 4 POST
	putMethod                     // 5 PUT
	patchMethod                   // 6 PATCH
	deleteMethod                  // 7 DELETE
	traceMethod                   // 8 TRACE
	connectMethod                 // 9 CONNECT
)

// maps String constants
var maps = map[Type]string{
	optionsMethod: `OPTIONS`,
	getMethod:     `GET`,
	headMethod:    `HEAD`,
	postMethod:    `POST`,
	putMethod:     `PUT`,
	patchMethod:   `PATCH`,
	deleteMethod:  `DELETE`,
	traceMethod:   `TRACE`,
	connectMethod: `CONNECT`,
}

// Type Type of methods
type Type int

// Value Value is an interface of method
type Value interface {
	// Int Return method as int constant
	Int() int

	// String Return method as string constant
	String() string

	// Type Return method as Type constant
	Type() Type

	// EqualFold Reports whether s, are equal value of method with case-folding
	EqualFold(s string) bool
}

// methodType is an implementation of Value
type methodType struct {
	value Type
}

// Interface is an methods interface
type Interface interface {
	// Options Return HTTP method OPTIONS
	Options() Value

	// Get Return HTTP method GET
	Get() Value

	// Head Return HTTP method GET
	Head() Value

	// Post Return HTTP method POST
	Post() Value

	// Put Return HTTP method PUT
	Put() Value

	// Patch Return HTTP method PATCH
	Patch() Value

	// Delete Return HTTP method DELETE
	Delete() Value

	// Trace Return HTTP method TRACE
	Trace() Value

	// Connect Return HTTP method CONNECT
	Connect() Value

	// Parse string and return value interface
	Parse(inp string) Value
}

// is an methods implementation
type impl struct {
}
