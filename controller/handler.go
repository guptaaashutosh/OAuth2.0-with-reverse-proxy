package controller

import (
	hydraAdmin "github.com/ory/hydra-client-go/client/admin"
	hydraClient "github.com/ory/hydra-client-go/client/public"
)

type Handler struct {
	HydraAdmin hydraAdmin.ClientService
	HydraClient hydraClient.ClientService
}
