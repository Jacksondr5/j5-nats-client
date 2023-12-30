package call

import (
	"fmt"
)


func CallHassService(domain string, service string, bodyData string) error {
	const accessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJlNjkyMGRhMDc5NWU0ZThmOGUzYzYyOTAzYzgwZmE0NyIsImlhdCI6MTcwMzg4MjQ0OCwiZXhwIjoyMDE5MjQyNDQ4fQ.BWphYONMeYF2Z64N6uAhhqNIOG3D8FfE3RjSR9XgrtM"
	err := HttpPost(
		fmt.Sprintf("http://hass.j5:8123/api/services/%s/%s", domain, service),
		bodyData, 
		fmt.Sprintf("Home Assistant service %s.%s", domain, service), 
		accessToken,
	)
	return err
}