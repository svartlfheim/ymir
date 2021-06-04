package registry

type RegistryHandlerStatus string

const STATUS_INTERNAL_ERROR RegistryHandlerStatus = "INTERNAL_ERROR"
const STATUS_NOT_FOUND RegistryHandlerStatus = "NOT_FOUND"
const STATUS_OKAY RegistryHandlerStatus = "OK"
const STATUS_INVALID RegistryHandlerStatus = "INVALID_PARAMS"
const STATUS_CREATED RegistryHandlerStatus = "CREATED"
const STATUS_MODIFIED RegistryHandlerStatus = "MODIFIED"
const STATUS_CONFLICT RegistryHandlerStatus = "CONFLICT"
