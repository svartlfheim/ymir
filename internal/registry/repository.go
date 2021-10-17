package registry

type ModuleRepository interface {
	ById(id string) (m Module, err error)
	ByFQN(ModuleFQN) (m Module, err error)
	All(chunkOpts ChunkingOptions, filters ModuleFilters) ([]Module, error)

	VersionById(id string) (m ModuleVersion, err error)
	VersionsByModule(moduleId string, chunkOpts ChunkingOptions) (m []ModuleVersion, err error)
	VersionsByModuleFQN(fqn ModuleFQN, chunkOpts ChunkingOptions) (m []ModuleVersion, err error)
	VersionByModuleAndValue(moduleId string, version string) (mv ModuleVersion, err error)
	VersionByFQN(fqn ModuleVersionFQN) (mv ModuleVersion, err error)

	AddModule(Module) (Module, error)
	DeleteModule(Module) error

	AddVersion(ModuleVersion) (m ModuleVersion, err error)
	DeleteVersionsForModule(Module) error
	DeleteModuleVersion(ModuleVersion) error
}
