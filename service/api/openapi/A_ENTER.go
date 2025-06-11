package openapi

import v1 "sun-panel/api/openapi/v1"

type Openapi struct {
	Apiv1 v1.Api
}

var ApiApp = new(Openapi)
