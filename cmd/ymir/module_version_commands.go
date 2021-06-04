package ymir

import (
	"github.com/svartlfheim/ymir/internal/output"
	"github.com/svartlfheim/ymir/internal/registry"
)

func module_version_add(c YmirCommand) error {
	o := c.GetOutput()
	l := c.GetLogger()
	f, err := c.cobra.LocalFlags().GetString("file")

	if err != nil {
		o.Error("the 'file' option was not configured for this command")
		return nil
	}

	if f == "" {
		o.Infoln("No file supplied entering interactive mode...")
	}

	cb := buildCommandBus(c)

	res, err := cb.AddModuleVersionV1FromCLI(f)

	if err != nil {
		switch err.(type) {
		case registry.ErrCouldNotReadFile, registry.ErrCouldNotUnmarshalJSONToDTO:
			o.Errorf("Could not load file at '%s', is it valid JSON?\n", f)
			return nil
		case registry.ErrQuestionFailed:
			l.Error().Err(err).Msg("failed to get questions responses")
			o.Errorln("Failed to read questions responses!")
			return nil
		default:
			o.Errorln("Whoopsie, an unexpected error occurred, see logs!")
			return nil
		}
	}

	switch res.Status {
	case registry.STATUS_INVALID:
		o.Errorln("Data was invalid!")
		for _, err := range res.ValidationErrors {
			o.Errorf("%s: %s\n", err.Field, err.Message)
		}
	case registry.STATUS_CREATED:
		o.Successln("Successfully created!")
		o.Successf("Id: %s\n", res.ModuleVersion.Id)
		o.Successf("Module ID: %s\n", res.ModuleVersion.ModuleId)
		o.Successf("Version: %s\n", res.ModuleVersion.Version)
		o.Successf("Source: %s\n", res.ModuleVersion.Source)
		o.Successf("Repository URL: %s\n", res.ModuleVersion.RepositoryURL)
		o.Successf("Status: %s\n", string(res.ModuleVersion.Status))
		o.Successln("Download URL: pending")
	default:
		o.Errorln("Whoopsie, an unexpected error occurred, see logs!")
	}

	return nil
}

func module_version_list(c YmirCommand) error {
	o := c.GetOutput()

	_, err := c.cobra.LocalFlags().GetString("output")

	if err != nil {
		o.Error("the 'output' option was not configured for this command")
		return nil
	}

	idOrFQN := c.GetArg(0, "")

	cb := buildCommandBus(c)

	res, err := cb.ListModuleVersionsV1FromCLI(idOrFQN)

	if err != nil {
		if _, ok := err.(registry.ErrCouldNotParseModuleFQN); ok {
			o.Errorln(err.Error())
			return nil
		}
		o.Errorln("Whoopsie, an unexpected error occurred, see logs!")
		return nil
	}

	tf := buildTableFactory()

	switch res.Status {
	case registry.STATUS_NOT_FOUND:
		o.Warnln("No versions found for this module!")
	case registry.STATUS_OKAY:
		if len(res.List) == 0 {
			o.Warnln("No versions found for this module!")
			return nil
		}

		h, r := registry.BuildModuleVersionsTable(res.List)
		tf.CreateAndPrint(h, r, output.WithAutoMergeByIndexes([]int{0}))
	default:
		o.Errorln("Whoopsie, an unexpected error occurred, see logs!")
	}

	return nil
}

func module_version_show(c YmirCommand) error {
	o := c.GetOutput()

	_, err := c.cobra.LocalFlags().GetString("output")

	if err != nil {
		o.Error("the 'output' option was not configured for this command")
		return nil
	}

	idOrFQN := c.GetArg(0, "")

	cb := buildCommandBus(c)

	res, err := cb.ShowModuleVersionV1FromCLI(idOrFQN)

	if err != nil {
		if _, ok := err.(registry.ErrCouldNotParseModuleVersionFQN); ok {
			o.Errorln(err.Error())
			return nil
		}

		o.Errorln("Whoopsie, an unexpected error occurred, see logs!")
		return nil
	}

	switch res.Status {
	case registry.STATUS_NOT_FOUND:
		o.Warnln("Module version not found!")
	case registry.STATUS_INVALID:
		o.Errorln("Data was invalid!")
		for _, err := range res.ValidationErrors {
			o.Errorf("%s: %s\n", err.Field, err.Message)
		}
	case registry.STATUS_OKAY:
		o.Successf("Id: %s\n", res.ModuleVersion.Id)
		o.Successf("Module Id: %s\n", res.ModuleVersion.ModuleId)
		o.Successf("Version: %s\n", res.ModuleVersion.Version)
		o.Successf("Repository URL: %s\n", res.ModuleVersion.RepositoryURL)
		o.Successf("Download URL: %s\n", res.ModuleVersion.DownloadURL)
		o.Successf("Status: %s\n", string(res.ModuleVersion.Status))
	default:
		o.Errorln("Whoopsie, an unexpected error occurred, see logs!")
	}

	return nil
}

func module_version_delete(c YmirCommand) error {
	o := c.GetOutput()

	forceDelete, err := c.cobra.LocalFlags().GetBool("yes")

	if err != nil {
		o.Error("the 'yes' option was not configured for this command")
		return nil
	}

	idOrFQN := c.GetArg(0, "")

	cb := buildCommandBus(c)

	res, err := cb.DeleteModuleVersionV1FromCLI(idOrFQN, forceDelete)

	if err != nil {
		if _, ok := err.(registry.ErrCouldNotParseModuleVersionFQN); ok {
			o.Errorln(err.Error())
			return nil
		}

		o.Errorln("Whoopsie, an unexpected error occurred, see logs!")
		return nil
	}

	switch res.Status {
	case registry.STATUS_NOT_FOUND:
		o.Warnln("Module version not found!")
	case registry.STATUS_INVALID:
		o.Errorln("Data was invalid!")
		for _, err := range res.ValidationErrors {
			o.Errorf("%s: %s\n", err.Field, err.Message)
		}
	case registry.STATUS_OKAY:
		o.Successf("Id: %s\n", res.ModuleVersion.Id)
		o.Successf("Module Id: %s\n", res.ModuleVersion.ModuleId)
		o.Successf("Version: %s\n", res.ModuleVersion.Version)
		o.Successf("Repository URL: %s\n", res.ModuleVersion.RepositoryURL)
		o.Successf("Download URL: %s\n", res.ModuleVersion.DownloadURL)
		o.Successf("Status: %s\n", string(res.ModuleVersion.Status))
	default:
		o.Errorln("Whoopsie, an unexpected error occurred, see logs!")
	}

	return nil
}
