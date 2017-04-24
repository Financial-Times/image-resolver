package content

import (
	"strings"
	log "github.com/Sirupsen/logrus"
)

type Resolver interface {
	UnrollImages(uuid string) (UnrolledContent, bool, error)
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

func (ir *ImageResolver) UnrollImages(uuid string) (UnrolledContent, bool, error) {
	var result UnrolledContent

	art, err := ir.reader.Get(uuid)
	if (err != nil) || (art.UUID == "") {
		return result, false, err
	}

	//mainImage
	mi := art.MainImage
	if mi != nil {
		id := extractIdfromUrl(mi.Id)
		if id != "" {
			mainImage, err := ir.getImage(id)
			if err != nil {
				return result, true, err
			}

			if len(mainImage) == 1 {
				result.MainImage = &mainImage[0]
			}
		} else {
			result.MainImage = nil
		}
	}

	//embedded images
	emImagesUUIDs, err := ir.parser.GetEmbedded(art)
	if err == nil {
		embeddedImages, err := ir.getEmbeddedImages(emImagesUUIDs)
		if err != nil {
			return result, true, err
		}
		result.Embeds = embeddedImages
	} else {
		log.Infof("Error parsing body for article uuid=%s, err=%v", art.UUID, err)
	}

	//promotional image
	id := extractIdfromUrl(art.AltImages.PromotionalImage)
	if id != "" {
		promotionalImage, err := ir.getImage(id)
		if err != nil {
			return result, true, err
		}
		if len(promotionalImage) == 1 {
			result.AlternativeImages = &PromotionalImage{&promotionalImage[0]}
		}
	} else {
		result.AlternativeImages = nil
	}

	return result, true, nil
}

func extractIdfromUrl(url string) string {
	id := ""
	splitUrl := strings.Split(url, "/")
	length := len(splitUrl)
	if (length) >= 1 {
		id = splitUrl[length-1]
	}
	return id
}

func (ir *ImageResolver) getImage(uuid string) ([]Content, error) {
	return ir.getImageSets([]string{uuid})
}

func (ir *ImageResolver) getEmbeddedImages(UUIDs []string) ([]Content, error) {
	return ir.getImageSets(UUIDs)
}

func (ir *ImageResolver) getImageSets(uuids []string) ([]Content, error) {
	outputs := []Content{}
	for _, uuid := range uuids {
		imageSet, err := ir.reader.Get(uuid)
		if err != nil {
			return outputs, err
		}
		if imageSet.UUID != "" {
			membersIDs := []string{}
			for _, b := range imageSet.Members {
				id := b.UUID
				membersIDs = append(membersIDs, extractIdfromUrl(id))
			}
			members, err := ir.getImageSetMembers(membersIDs)
			if err != nil {
				return outputs, err
			}
			imageSet.Members = members
			outputs = append(outputs, imageSet)
		}
	}
	return outputs, nil
}

func (ir *ImageResolver) getImageSetMembers(membersUUIDs []string) ([]Content, error) {
	membersCh := []Content{}
	for _, member := range membersUUIDs {
		im, err := ir.reader.Get(member)
		if err != nil {
			return membersCh, err
		}
		membersCh = append(membersCh, im)
	}
	return membersCh, nil
}

