package server

// Todo: add state if needed
type ZKServer struct {
	Port    int
	KeysDir string
}

func NewZKServer(Port int, keysDir string) *ZKServer {
	return &ZKServer{
		Port:    Port,
		KeysDir: keysDir,
	}
}
