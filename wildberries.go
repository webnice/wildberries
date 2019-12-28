package wildberries

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
//import ()

// New Creates an new object of package and return interface
func New(apiKey string) Interface {
	var wbs = &impl{
		apiKey: apiKey,
	}
	return wbs
}
