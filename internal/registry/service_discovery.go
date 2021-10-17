package registry

type ServiceDiscoveryCommand struct{}

type HandleServiceDiscoveryResponseBody struct {
	ModuleVersion string `json:"modules.v1"`
}

type HandleServiceDiscoveryResponse struct {
	Status RegistryHandlerStatus
	Body   HandleServiceDiscoveryResponseBody
}

func (_ ServiceDiscoveryCommand) Handle() HandleServiceDiscoveryResponse {
	return HandleServiceDiscoveryResponse{
		Status: STATUS_OKAY,
		Body: HandleServiceDiscoveryResponseBody{
			ModuleVersion: "/v1/modules/",
		},
	}
}
