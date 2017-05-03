package content

import (
	"strings"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

type Resolver interface {
	UnrollImages(Content) (Content)
	ExtractIdfromUrl(string) (string)
}

type ImageResolver struct {
	reader Reader
	parser Parser
}

func NewImageResolver(r *Reader, p *Parser) *ImageResolver {
	return &ImageResolver{
		reader: *r,
		parser: *p,
	}
}

func (ir *ImageResolver) UnrollImages(art Content) (Content) {
	result := art

	//mainImage
	mi := art.MainImage
	if mi != nil {
		mainImage := ir.getImage(mi.UUID)
		if len(mainImage) == 1 {
			result.MainImage = &mainImage[0]
		}
	}

	//embedded images
	if art.BodyXML != nil {
		emImagesUUIDs, err := ir.parser.GetEmbedded(art)
		if err == nil {
			embeddedImages := ir.getEmbeddedImages(emImagesUUIDs)
			result.Embeds = embeddedImages
		} else {
			result.Embeds = nil
			log.Infof("Error parsing body for article uuid=%s, err=%v", art.UUID, err)
		}
	}

	//promotional image
	if art.AltImages != nil {
		uuid := art.AltImages.(map[string]interface{})
		if uuid != nil {
			id := uuid["promotionalImage"]
			if id != nil && id != "" {
				promotionalImage := ir.getImage(id.(string))
				if len(promotionalImage) == 1 {
					result.AltImages = &PromotionalImage{&promotionalImage[0]}
				}
			}
		}
	}

	//lead images
	if art.LeadImages != nil {
		images := art.LeadImages.([]interface{})
		result.LeadImages = ir.getLeadImages(images)
	}

	return result

}

func (ir *ImageResolver) ExtractIdfromUrl(url string) string {
	id := ""
	splitUrl := strings.Split(url, "/")
	length := len(splitUrl)
	if (length) >= 1 {
		id = splitUrl[length-1]
	}
	return id
}

func (ir *ImageResolver) getImage(uuid string) []Content {
	return ir.getImageSets([]string{uuid})
}

func (ir *ImageResolver) getEmbeddedImages(UUIDs []string) []Content {
	return ir.getImageSets(UUIDs)
}

func (ir *ImageResolver) getImageSets(uuids []string) []Content {
	outputs := []Content{}
	for _, uuid := range uuids {
		id := ir.ExtractIdfromUrl(uuid)
		imageSet, err, code := ir.reader.Get(id)
		if code != http.StatusOK {
			imageSet.UUID = uuid
			outputs = append(outputs, imageSet)
			log.Infof("Error unrolling image uuid =%s, err=%v", id, err)
		} else {
			if imageSet.UUID != "" {
				membersIDs := []string{}
				for _, b := range imageSet.Members {
					membersIDs = append(membersIDs, b.UUID)
				}
				members := ir.getImageSetMembers(membersIDs)
				imageSet.Members = members
				outputs = append(outputs, imageSet)
			} else {
				imageSet.UUID = uuid
				outputs = append(outputs, imageSet)
			}
		}
	}
	return outputs
}

func (ir *ImageResolver) getImageSetMembers(membersUUIDs []string) []Content {
	members := []Content{}
	for _, member := range membersUUIDs {
		id := ir.ExtractIdfromUrl(member)
		im, err, code := ir.reader.Get(id)
		if code != http.StatusOK {
			im.UUID = member
			log.Infof("Error unrolling image uuid =%s, err=%v", member, err)
		}
		members = append(members, im)
	}
	return members
}

func (ir *ImageResolver) getLeadImages(leadImages []interface{}) ([]ImageOutput) {
	result := make([]ImageOutput, len(leadImages))
	for index, leadImg := range leadImages {
		img := leadImg.(map[string]interface{})
		id := img["id"].(string)
		im, err, code := ir.reader.Get(ir.ExtractIdfromUrl(id))
		if code != http.StatusOK {
			result[index].Content = id
			result[index].Type = img["type"].(string)
			log.Infof("Error unrolling leadimage uuid =%s, err=%v", id, err)

		} else {
			result[index].Content = im
			result[index].Type = img["type"].(string)
		}
	}

	return result
}
