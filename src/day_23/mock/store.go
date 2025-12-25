package mock

type Store interface {
	Save(name string) error
}

func DoSave(st Store) error {
	return st.Save(`Codee君`)
}

type mockStore struct {
	called bool
}

func (m *mockStore) Save(name string) error {
	m.called = true
	return nil
}

