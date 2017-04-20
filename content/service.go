package content

import (
	"strings"
	"fmt"
)

type Resolver interface {
	UnrollImages(uuid string) (UnrolledContent, bool, error)
	UnrollLeadImages(uuid string) (UnrolledLeadImagesContent, bool, error)
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
		fmt.Println("Empty body")
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
	contentOutputs := []Content{}
	for _, uuid := range uuids {
		imageSet, err := ir.reader.Get(uuid)
		if err != nil {
			return contentOutputs, err
		}
		membersIDs := []string{}
		for _, b := range imageSet.Members {
			id := b.UUID
			membersIDs = append(membersIDs, extractIdfromUrl(id))
		}
		members, err := ir.getImageSetMembers(membersIDs)
		if err != nil {
			return contentOutputs, err
		}
		imageSet.Members = members
		contentOutputs = append(contentOutputs, imageSet)
	}
	return contentOutputs, nil
}

func (ir *ImageResolver) getImageSetMembers(membersUUIDs []string) ([]Content, error) {
	membersCh := []Content{}
	for _, member := range membersUUIDs {
		im, err := ir.reader.Get(member)
		fmt.Print(im)
		if err != nil {
			return membersCh, err
		}
		membersCh = append(membersCh, im)
	}

	return membersCh, nil
}

func (ir *ImageResolver) UnrollLeadImages(uuid string) (UnrolledLeadImagesContent, bool, error) {
	var result = UnrolledLeadImagesContent{}
	return result, true, nil
}
