package grpc

type AntiBruteForceServer struct {
}

func NewAntiBruteForceServer() *AntiBruteForceServer {
	return &AntiBruteForceServer{}
}

func (a *AntiBruteForceServer) Serve(addr string) error {
	return nil
}
