package content

import (
	"regexp"

	uuidutils "github.com/Financial-Times/uuid-utils-go"
)

const uuidRegex = "([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})"

func extractUUIDFromURL(url string) string {
	re, _ := regexp.Compile(uuidRegex)
	values := re.FindStringSubmatch(url)
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func getImageModelUUID(imageSetUUID string) (imageModelUUID string, err error) {
	uuidDeriver := uuidutils.NewUUIDDeriverWith(uuidutils.IMAGE_SET)
	uuid, err := uuidutils.NewUUIDFromString(imageSetUUID)
	if err != nil {
		return imageModelUUID, err
	}
	uuid, err = uuidDeriver.From(uuid)
	imageModelUUID = uuid.String()
	return
}

func createID(APIHost string, handlerPath string, uuid string) string {
	return "http://" + APIHost + "/" + handlerPath + "/" + uuid
}
