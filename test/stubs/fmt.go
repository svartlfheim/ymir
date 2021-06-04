package ymirstubs

type PrintFCall struct {
	Format string
	Args   []interface{}
}

type Fmt struct {
	InfolnCalls [][]interface{}
	InfofCalls  []PrintFCall
	InfoCalls   [][]interface{}

	WarnlnCalls [][]interface{}
	WarnfCalls  []PrintFCall
	WarnCalls   [][]interface{}

	SuccesslnCalls [][]interface{}
	SuccessfCalls  []PrintFCall
	SuccessCalls   [][]interface{}

	ErrorlnCalls [][]interface{}
	ErrorfCalls  []PrintFCall
	ErrorCalls   [][]interface{}
}

func (f *Fmt) Infoln(a ...interface{}) (int, error) {
	f.InfolnCalls = append(f.InfolnCalls, a)
	return 0, nil
}
func (f *Fmt) Infof(format string, a ...interface{}) (int, error) {
	f.InfofCalls = append(f.InfofCalls, PrintFCall{
		Format: format,
		Args:   a,
	})

	return 0, nil
}
func (f *Fmt) Info(a ...interface{}) (int, error) {
	f.InfoCalls = append(f.InfoCalls, a)
	return 0, nil
}

func (f *Fmt) Warnln(a ...interface{}) (int, error) {
	f.WarnlnCalls = append(f.WarnlnCalls, a)
	return 0, nil
}
func (f *Fmt) Warnf(format string, a ...interface{}) (int, error) {
	f.WarnfCalls = append(f.WarnfCalls, PrintFCall{
		Format: format,
		Args:   a,
	})

	return 0, nil
}
func (f *Fmt) Warn(a ...interface{}) (int, error) {
	f.WarnCalls = append(f.WarnCalls, a)
	return 0, nil
}

func (f *Fmt) Errorln(a ...interface{}) (int, error) {
	f.ErrorlnCalls = append(f.ErrorlnCalls, a)
	return 0, nil
}
func (f *Fmt) Errorf(format string, a ...interface{}) (int, error) {
	f.ErrorfCalls = append(f.ErrorfCalls, PrintFCall{
		Format: format,
		Args:   a,
	})

	return 0, nil
}
func (f *Fmt) Error(a ...interface{}) (int, error) {
	f.ErrorCalls = append(f.ErrorCalls, a)
	return 0, nil
}

func (f *Fmt) Successln(a ...interface{}) (int, error) {
	f.SuccesslnCalls = append(f.SuccesslnCalls, a)
	return 0, nil
}
func (f *Fmt) Successf(format string, a ...interface{}) (int, error) {
	f.SuccessfCalls = append(f.SuccessfCalls, PrintFCall{
		Format: format,
		Args:   a,
	})

	return 0, nil
}
func (f *Fmt) Success(a ...interface{}) (int, error) {
	f.SuccessCalls = append(f.SuccessCalls, a)
	return 0, nil
}

func BuildFmtStub() *Fmt {
	return &Fmt{
		InfolnCalls: [][]interface{}{},
		InfofCalls:  []PrintFCall{},
		InfoCalls:   [][]interface{}{},

		WarnlnCalls: [][]interface{}{},
		WarnfCalls:  []PrintFCall{},
		WarnCalls:   [][]interface{}{},

		SuccesslnCalls: [][]interface{}{},
		SuccessfCalls:  []PrintFCall{},
		SuccessCalls:   [][]interface{}{},

		ErrorlnCalls: [][]interface{}{},
		ErrorfCalls:  []PrintFCall{},
		ErrorCalls:   [][]interface{}{},
	}
}
