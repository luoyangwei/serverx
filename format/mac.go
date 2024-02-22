package format

import "strings"

func SubMac(mac string) string {
	return strings.Replace(mac, ":", "", -1)
}
