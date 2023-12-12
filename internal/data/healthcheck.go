package data

type Healthcheck struct {
	Status     string `json:"status"`
	Enviroment string `json:"enviroment"`
	Version    string `json:"version"`
}
