package content

import (
	log "github.com/Sirupsen/logrus"
	"regexp"
)

const (
	MainImage  = "mainImage"
	ID         = "id"
	Embeds     = "embeds"
	AltImages  = "alternativeImages"
	LeadImages = "leadImages"
	Members    = "members"
	Type       = "type"
	BodyXml    = "bodyXML"
	PromoImage = "promotionalImage"
	UUID       = "uuid"
	AppName    = "image-resolver"
	Image 	   = "image"
)

type Content map[string]interface{}

type Resolver interface {
	UnrollImages(Content) Content
	UnrollLeadImages(Content) Content
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

func (ir *ImageResolver) UnrollImages(body Content) Content {
	result := body

	//mainImage
	mi, found := body[MainImage].(map[string]interface{})
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
	bodyXml, found := body[BodyXml].(string)
	if found && bodyXml != "" {
		emImagesUUIDs, err := ir.parser.GetEmbedded(bodyXml)
		if err == nil {
			embeddedImages := ir.getImageSets(emImagesUUIDs)
			if len(embeddedImages) > 0 {
				result[Embeds] = embeddedImages
			}
		} else {
			log.Errorf("Error parsing body for article uuid=%s, err=%v", body[ID].(string), err)
		}
	}

	//promotional image
	uuid, found := body[AltImages].(map[string]interface{})
	if found {
		id, ok := uuid[PromoImage].(string)
		if ok && id != "" {
			promotionalImage, err := ir.reader.Get(extractIdfromUrl(id))
			if err == nil {
				result[AltImages] =  map[string]interface{}{PromoImage: promotionalImage}
			} else {
				log.Errorf("Error unrolling promotional image uuid =%s, err=%v", id, err)
			}
		}
	}
	return result
}

func (ir *ImageResolver) UnrollLeadImages(body Content) Content {
	result := body
	//lead images
	images, found := body[LeadImages].([]interface{})
	if found {
		result[LeadImages] = ir.getLeadImages(images)
	}
	return result
}

func (ir *ImageResolver) getImage(uuid string) []Content {
	return ir.getImageSets([]string{uuid})
}

func (ir *ImageResolver) getImageSets(uuids []string) []Content {
	outputs := []Content{}
	for _, uuid := range uuids {
		id := extractIdfromUrl(uuid)
		imageSet, err := ir.reader.Get(id)
		if err != nil {
			imageSet = make(map[string]interface{})
			imageSet[ID] = uuid
			outputs = append(outputs, imageSet)
			log.Errorf("Error unrolling image uuid = %s, err: %v", id, err)
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
		id := extractIdfromUrl(member)
		im, err := ir.reader.Get(id)
		if err != nil {
			im[ID] = member
			log.Errorf("Error unrolling image uuid =%s, err=%v", member, err)
		}
		members = append(members, im)
	}
	return members
}

func (ir *ImageResolver) getLeadImages(leadImages []interface{}) []Content {
	result := make([]Content, len(leadImages))
	for index, leadImg := range leadImages {
		result[index] = make(map[string]interface{})
		img := leadImg.(map[string]interface{})
		id := img[ID].(string)
		im, err := ir.reader.Get(extractIdfromUrl(id))
		if err != nil{
			result[index][ID]= id
			result[index][Type] = img[Type].(string)
			log.Errorf("Error unrolling leadimage uuid =%s, err=%v", id, err)
		} else {
			result[index][Image] = im
			result[index][Type] = img[Type].(string)
			result[index][ID] = id
		}
	}
	return result
}

func extractIdfromUrl(url string) string {
	var id string
	re, _ := regexp.Compile("([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})")
	values := re.FindStringSubmatch(url)
	if len(values) > 0 {
		id = values[0]
	}
	return id
}