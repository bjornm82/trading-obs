package broker

type Mock struct {
	ReturnByte   []byte
	ReturnInt    int
	ReturnError  error
	ReturnString string
}

func (m *Mock) GetToken() string {
	return m.ReturnString
}

func (m *Mock) GetAccountID() (string, error) {
	return m.ReturnString, m.ReturnError
}

func (m *Mock) Request(endpoint string) ([]byte, int, error) {
	return m.ReturnByte, m.ReturnInt, m.ReturnError
}

func (m *Mock) Read(endpoint string) ([]byte, error) {
	return m.ReturnByte, m.ReturnError
}

func (m *Mock) ReadStream(endpoint string) ([]byte, error) {
	return m.ReturnByte, m.ReturnError
}

func (m *Mock) Update(endpoint string, data []byte) ([]byte, error) {
	return m.ReturnByte, m.ReturnError
}

func (m *Mock) Create(endpoint string, data []byte) ([]byte, error) {
	return m.ReturnByte, m.ReturnError
}
