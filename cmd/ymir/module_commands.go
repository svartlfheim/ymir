package ymir

import (
	"github.com/svartlfheim/ymir/internal/output"
	"github.com/svartlfheim/ymir/internal/registry"
)

func module_add(c YmirCommand) error {
	l := c.GetLogger()
	o := c.GetOutput()
	f, err := c.cobra.LocalFlags().GetString("file")

	if err != nil {
		o.Error("the 'file' option was not configured for this command")
		return nil
	}

	if f == "" {
		o.Infoln("No file supplied entering interactive mode...")
	}

	cb := buildCommandBus(c)

	res, err := cb.AddModuleV1FromCLI(f)

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
		o.Successf("Id: %s\n", res.Module.Id)
		o.Successf("Name: %s\n", res.Module.Name)
		o.Successf("Namespace: %s\n", res.Module.Namespace)
		o.Successf("Provider: %s\n", res.Module.Provider)
	default:
		o.Errorln("Whoopsie, an unexpected error occurred, see logs!")
	}

	return nil
}

func module_list(c YmirCommand) error {
	o := c.GetOutput()

	_, err := c.cobra.LocalFlags().GetString("output")

	if err != nil {
		o.Error("the 'output' option was not configured for this command")
		return nil
	}

	provider, err := c.cobra.LocalFlags().GetString("provider")

	if err != nil {
		o.Error("the 'provider' option was not configured for this command")
		return nil
	}

	ns, err := c.cobra.LocalFlags().GetString("namespace")

	if err != nil {
		o.Error("the 'namespace' option was not configured for this command")
		return nil
	}

	cb := buildCommandBus(c)

	res, err := cb.ListModulesV1FromCLI(provider, ns)

	if err != nil {
		o.Errorln("Whoopsie, an unexpected error occurred, see logs!")
		return nil
	}

	tf := buildTableFactory()

	switch res.Status {
	case registry.STATUS_OKAY:
		if len(res.List) == 0 {
			o.Warnln("No modules found!")
			return nil
		}

		h, r := registry.BuildModuleTable(res.List)
		tf.CreateAndPrint(h, r, output.WithAutoMergeByIndexes([]int{0, 1}))
	default:
		o.Errorln("Whoopsie, an unexpected error occurred, see logs!")
	}

	return nil
}

func module_show(c YmirCommand) error {
	o := c.GetOutput()

	_, err := c.cobra.LocalFlags().GetString("output")

	if err != nil {
		o.Error("the 'output' option was not configured for this command")
		return nil
	}

	idOrFQN := c.GetArg(0, "")

	cb := buildCommandBus(c)

	res, err := cb.ShowModuleV1FromCLI(idOrFQN)

	if err != nil {
		o.Errorln("Whoopsie, an unexpected error occurred, see logs!")
		return nil
	}

	switch res.Status {
	case registry.STATUS_INVALID:
		o.Errorln("Data was invalid!")
		for _, err := range res.ValidationErrors {
			o.Errorf("%s: %s\n", err.Field, err.Message)
		}
	case registry.STATUS_NOT_FOUND:
		o.Errorln("Module not found!")
	case registry.STATUS_OKAY:
		o.Successf("Id: %s\n", res.Module.Id)
		o.Successf("Name: %s\n", res.Module.Name)
		o.Successf("Namespace: %s\n", res.Module.Namespace)
		o.Successf("Provider: %s\n", res.Module.Provider)
	default:
		o.Errorln("Whoopsie, an unexpected error occurred, see logs!")
	}

	return nil
}

func module_delete(c YmirCommand) error {
	o := c.GetOutput()

	forceDelete, err := c.cobra.LocalFlags().GetBool("yes")

	if err != nil {
		o.Error("the 'yes' option was not configured for this command")
		return nil
	}

	deleteVersions, err := c.cobra.LocalFlags().GetBool("versions")

	if err != nil {
		o.Error("the 'versions' option was not configured for this command")
		return nil
	}

	idOrFQN := c.GetArg(0, "")

	cb := buildCommandBus(c)

	res, err := cb.DeleteModuleV1FromCLI(idOrFQN, deleteVersions, forceDelete)

	if err != nil {
		if _, ok := err.(registry.ErrFailedToConfirmAction); ok {
			o.Errorln("Aborted due to failed confirmation!")

			return nil
		}
		o.Errorln("Whoopsie, an unexpected error occurred, see logs!")
		return nil
	}

	switch res.Status {
	case registry.STATUS_INVALID:
		o.Errorln("Data was invalid!")
		for _, err := range res.ValidationErrors {
			o.Errorf("%s: %s\n", err.Field, err.Message)
		}
	case registry.STATUS_NOT_FOUND:
		o.Errorln("Module not found!")
	case registry.STATUS_OKAY:
		o.Successf("Module %s/%s/%s (%s) deleted successfully!\n", res.Module.Provider, res.Module.Namespace, res.Module.Name, res.Module.Id)
	default:
		o.Errorln("Whoopsie, an unexpected error occurred, see logs!")
	}

	return nil
}
