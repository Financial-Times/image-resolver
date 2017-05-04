package content

import (
	"strings"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

const (
	MainImage = "mainImage"
	ID        = "id"
	Embeds 	  = "embeds"
	AltImages = "alternativeImages"
	LeadImages= "leadImages"
	Members   = "members"
	Type      = "type"
	BodyXml   = "bodyXML"
        PromoImage= "promotionalImage"
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
	mi, found := art[MainImage].(map[string]interface{})
	if found {
		id, ok := mi[ID].(string)
		if ok && id != "" {
			mainImage := ir.getImage(id)
			if len(mainImage) == 1 {
				result[MainImage] = mainImage[0]
			}
		}
	}

	//embedded images
	bodyXml, found := art[BodyXml].(string)
	if found && bodyXml != "" {
		emImagesUUIDs, err := ir.parser.GetEmbedded(bodyXml)
		if err == nil {
			embeddedImages := ir.getEmbeddedImages(emImagesUUIDs)
			result[Embeds] = embeddedImages
		} else {
			result[Embeds] = nil
			log.Infof("Error parsing body for article uuid=%s, err=%v", art[ID].(string), err)
		}
	}

	//promotional image
	uuid, found := art[AltImages].(map[string]interface{})
	if found {
		id, ok := uuid[PromoImage].(string)
		if ok && id != ""{
			promotionalImage := ir.getImage(id)
			if len(promotionalImage) == 1 {
				promo := PromotionalImage{promotionalImage[0]}
				result[AltImages] = promo
			}
		}
	}

	//lead images
	images, found := art[LeadImages].([]interface{})
	if found {
		result[LeadImages] = ir.getLeadImages(images)
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
			imageSet[ID] = uuid
			outputs = append(outputs, imageSet)
			log.Infof("Error unrolling image uuid =%s, err=%v", id, err)
		} else {
			if imageSet[ID] != "" {
				membersIDs := []string{}
				membs, ok := imageSet[Members].([]interface{})
				if ok {
					for _, b := range membs {
						b, ok := b.(map[string]interface{})
						if !ok {
							continue
						}
						membersIDs = append(membersIDs, b[ID].(string))
					}
					members := ir.getImageSetMembers(membersIDs)
					imageSet[Members] = members
					outputs = append(outputs, imageSet)
				}
			} else {
				imageSet[ID] = uuid
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
			im[ID] = member
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
		id := img[ID].(string)
		im, err, code := ir.reader.Get(ir.ExtractIdfromUrl(id))
		if code != http.StatusOK {
			result[index].Content = map[string]interface{}{"id": id}
			result[index].Type = img[Type].(string)
			log.Infof("Error unrolling leadimage uuid =%s, err=%v", id, err)
		} else {
			result[index].Content = im
			result[index].Type = img[Type].(string)
		}
	}
	return result
}
