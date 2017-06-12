package content

import (
	"regexp"
)

const uuidRegex = "([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})"

func extractUUIDFromString(url string) string {
	re, err := regexp.Compile(uuidRegex)
	if err != nil {
		return ""
	}

	values := re.FindStringSubmatch(url)
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func createID(APIHost string, handlerPath string, uuid string) string {
	return "http://" + APIHost + "/" + handlerPath + "/" + uuid
}
