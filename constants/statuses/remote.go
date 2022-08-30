package statuses

const (
	Waiting    = "waiting"
	RoeRunning = "running"
	Stopping   = "stopping"
	Erring     = "error"
	Finished   = "finished"
	Lining     = "lining"
)

var RoeStatues = [6]string{Waiting, RoeRunning, Stopping, Erring, Finished, Lining}
