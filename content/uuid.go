package content

import (
	"regexp"

	"github.com/pkg/errors"
)

const uuidRegex = "([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})"

func extractUUIDFromString(url string) (string, error) { // both
	re, err := regexp.Compile(uuidRegex)
	if err != nil {
		return "", errors.Wrap(err, "Error during extracting UUID")
	}

	values := re.FindStringSubmatch(url)
	if len(values) > 0 {
		return values[0], nil
	}
	return "", errors.Errorf("Cannot extract UUID from %s", url)
}

func createID(APIHost string, handlerPath string, uuid string) string { // both
	return "http://" + APIHost + "/" + handlerPath + "/" + uuid
}
