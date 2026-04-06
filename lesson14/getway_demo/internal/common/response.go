package common

type InfoResponse struct {
	Service   string            `json:"service"`
	Instance  string            `json:"instance"`
	Method    string            `json:"method,omitempty"`
	Path      string            `json:"path"`
	Host      string            `json:"host"`
	ClientIP  string            `json:"client_ip,omitempty"`
	Forwarded map[string]string `json:"forwarded"`
}

func BuildInfoResponse(service, instance, path, host string, forwarded map[string]string) InfoResponse {
	return InfoResponse{
		Service:   service,
		Instance:  instance,
		Path:      path,
		Host:      host,
		Forwarded: forwarded,
	}
}

func BuildRequestInfoResponse(service, instance, method, path, host, clientIP string, forwarded map[string]string) InfoResponse {
	resp := BuildInfoResponse(service, instance, path, host, forwarded)
	resp.Method = method
	resp.ClientIP = clientIP
	return resp
}
