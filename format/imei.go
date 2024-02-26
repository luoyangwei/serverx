package format

import "strings"

func IMEI(mac string) string {
	return strings.Replace(mac, ":", "", -1)
}
