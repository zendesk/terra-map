package service

type Service interface {
	Process(serv *[]Service, resources []string, b []byte)
	Conditions() []Condition
}

type Condition struct {
	MonitorType     string
	MonitorMessage  string
	MonitorDuration int
}
