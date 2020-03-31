module routing-manager

go 1.12.1

require (
	gerrit.o-ran-sc.org/r/ric-plt/xapp-frame v0.4.4
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/errors v0.19.3
	github.com/go-openapi/loads v0.19.4
	github.com/go-openapi/runtime v0.19.4
	github.com/go-openapi/spec v0.19.3
	github.com/go-openapi/strfmt v0.19.4
	github.com/go-openapi/swag v0.19.7
	github.com/go-openapi/validate v0.19.6
	github.com/jessevdk/go-flags v1.4.0
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297
	nanomsg.org/go/mangos/v2 v2.0.5
)

replace gerrit.o-ran-sc.org/r/ric-plt/sdlgo => gerrit.o-ran-sc.org/r/ric-plt/sdlgo.git v0.5.2

replace gerrit.o-ran-sc.org/r/ric-plt/xapp-frame => gerrit.o-ran-sc.org/r/ric-plt/xapp-frame.git v0.4.4

replace gerrit.o-ran-sc.org/r/com/golog => gerrit.o-ran-sc.org/r/com/golog.git v0.0.1

replace nanomsg.org/go/mangos/v2 => nanomsg.org/go/mangos/v2 v2.0.5
