package content

import (
	"strings"
)

type Resolver interface {
	UnrollImages(uuid string) (UnrolledContent, bool, error)
	UnrollLeadImages(uuid string) (UnrolledLeadImagesContent,  bool, error)
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
	if (err != nil) || (art.Uuid == "") {
		return result, false, err
	}

	//mainImage
	id := extractIdfromUrl(art.MainImage.Id)
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

	//embedded images
	emImagesUUIDs, err := ir.parser.GetEmbedded(art)
	if err != nil {
		return result, true, err
	}
	embeddedImages, err := ir.getEmbeddedImages(emImagesUUIDs)
	if err != nil {
		return result, true, err
	}
	result.Embeds = embeddedImages

	//promotional image
	id = extractIdfromUrl(art.AltImages.PromotionalImage)
	if id != "" {
		promotionalImage, err := ir.getImage(id)
		if err != nil {
			return result, true, err
		}
		if len(promotionalImage) == 1 {
			result.PromotionalImage = &promotionalImage[0]
		}
	} else {
		result.PromotionalImage = nil
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

func (ir *ImageResolver) getImage(uuid string) ([]ContentOutput, error) {
	return ir.getImageSets([]string{uuid})
}

func (ir *ImageResolver) getEmbeddedImages(UUIDs []string) ([]ContentOutput, error) {
	return ir.getImageSets(UUIDs)
}

func (ir *ImageResolver) getImageSets(uuids []string) ([]ContentOutput, error) {
	contentOutputs := []ContentOutput{}
	for _, uuid := range uuids {
		imageSet, err := ir.reader.Get(uuid)
		if err != nil {
			return contentOutputs, err
		}
		if imageSet.Uuid != "" {
			membersSet := imageSet.Members
			membersIDs := []string{}
			for _, b := range membersSet {
				membersIDs = append(membersIDs, extractIdfromUrl(b.Id))
			}
			members, err := ir.getImageSetMembers(membersIDs)
			if err != nil {
				return contentOutputs, err
			}
			imageOutput := ir.reader.ConvertContentIntoOutput(imageSet)
			imageOutput.Members = members
			contentOutputs = append(contentOutputs, imageOutput)
		}
	}
	return contentOutputs, nil
}

func (ir *ImageResolver) getImageSetMembers(membersUUIDs []string) ([]Content, error) {
	membersCh := []Content{}
	for _, member := range membersUUIDs {
		im, err := ir.reader.Get(member)
		if err != nil {
			return membersCh, err
		}
		if im.Uuid != "" {
			membersCh = append(membersCh, im)
		}
	}

	return membersCh, nil
}

func (ir *ImageResolver) UnrollLeadImages(uuid string) (UnrolledLeadImagesContent,  bool, error) {
	var result = UnrolledLeadImagesContent{}
	return result, true, nil
}
