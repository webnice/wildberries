package communication // import "git.webdesk.ru/wd/kit/modules/communication"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// Все ошибки определены как константы
const (
	cUnauthorized = `Unauthorized`
	cForbidden    = `Forbidden`
	cNotFound     = `Not found`
)

// Константы указаны в объектах, адрес которых фиксирован всё время работы приложения
// Ошибку с ошибкой можно сравнивать по телу, по адресу и т.п.
var (
	errSingleton    = &Error{}
	errUnauthorized = err(cUnauthorized)
	errForbidden    = err(cForbidden)
	errNotFound     = err(cNotFound)
)

type (
	// Error object of package
	Error struct{}
	err   string
)

// Error The error built-in interface implementation
func (e err) Error() string { return string(e) }

// Errors Все ошибки известного состояния, которые могут вернуть функции пакета
func Errors() *Error { return errSingleton }

// ERRORS:

// Unauthorized Unauthorized error
func (e *Error) Unauthorized() error { return &errUnauthorized }

// Forbidden Forbidden error
func (e *Error) Forbidden() error { return &errForbidden }

// NotFound Not found
func (e *Error) NotFound() error { return &errNotFound }
