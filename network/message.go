package network

type GetStatusMessage struct {
}

type StatusMessage struct {
	ID           string
	ActualHeight uint32
	Version      uint32
}
