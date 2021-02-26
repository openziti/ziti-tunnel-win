package service

import "github.com/openziti/desktop-edge-win/service/ziti-tunnel/dto"

var GET_STATUS = dto.CommandMsg{
	Function: "Status",
}

var templateIdentity = "{{printf \"%40s\" \"Name\"}} | {{printf \"%41s\" \"FingerPrint\"}} | {{printf \"%6s\" \"Active\"}} | {{printf \"%30s\" \"Config\"}} | {{\"Status\"}}\n" +
	"{{range .}}{{printf \"%40s\" .Name}} | {{printf \"%41s\" .FingerPrint}} | {{printf \"%6t\" .Active}} | {{printf \"%30s\" .Config}} | {{.Status}}\n{{end}}"

var templateService = "{{printf \"%40s\" \"Name\"}} | {{printf \"%15s\" \"Id\"}} | {{printf \"%40s\" \"InterceptHost\"}} | {{printf \"%14s\" \"InterceptPort\"}} | {{printf \"%15s\" \"AssignedIP\"}} | {{printf \"%15s\" \"AssignedHost\"}} | {{\"OwnsIntercept\"}}\n" +
	"{{range .}}{{printf \"%40s\" .Name}} | {{printf \"%15s\" .Id}} | {{printf \"%40s\" .InterceptHost}} | {{printf \"%14d\" .InterceptPort}} | {{printf \"%15s\" .AssignedIP}} | {{printf \"%15s\" .AssignedHost}} | {{.OwnsIntercept}}\n{{end}}"

type IdentityCli struct {
	Name              string
	FingerPrint       string
	Active            bool
	Config            string
	ControllerVersion string
	Status            string
}

type ServiceCli struct {
	Name          string
	AssignedIP    string
	InterceptHost string
	InterceptPort uint16
	Id            string
	AssignedHost  string
	OwnsIntercept bool
	Identity      string
}
