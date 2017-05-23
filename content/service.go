package content

import (
	"github.com/pkg/errors"
)

const (
	mainImage  = "mainImage"
	id         = "id"
	embeds     = "embeds"
	altImages  = "alternativeImages"
	leadImages = "leadImages"
	members    = "members"
	bodyXML    = "bodyXML"
	promoImage = "promotionalImage"
	image      = "image"
)

type Resolver interface {
	UnrollImages(c Content) (Content, error)
	UnrollLeadImages(Content) (Content, error)
}

type ImageResolver struct {
	reader  Reader
	parser  Parser
	apiHost string
}
type Content map[string]interface{}

type imageSetUUID struct {
	uuid           string
	imageModelUUID string
}

type UUIDBatch struct {
	mainImageSet     imageSetUUID
	embeddedImages   []imageSetUUID
	promotionalImage string
	leadImages       []string
}

func NewImageResolver(r Reader, p Parser, apiHost string) *ImageResolver {
	return &ImageResolver{
		reader:  r,
		parser:  p,
		apiHost: apiHost,
	}
}

func (ir *ImageResolver) UnrollImages(c Content) (Content, error) {
	//make a copy of the content
	uuid := extractUUIDFromURL(c[id].(string))
	cc := c.clone()
	b := UUIDBatch{}
	var err error

	//mainImage
	mi, foundMainImg := cc[mainImage].(map[string]interface{})
	if foundMainImg {
		id := mi[id].(string)
		var mis imageSetUUID
		mis.uuid = extractUUIDFromURL(id)
		imgModelUUID, err := getImageModelUUID(mis.uuid)
		if err != nil {
			logger.Infof(uuid, "Cannot convert image set UUID %v to a image model UUID. Skipping item.", mis.uuid)
		} else {
			mis.imageModelUUID = imgModelUUID
		}
		b.mainImageSet = mis
	} else {
		logger.Infof(uuid, "Cannot find main image for %v. Skipping expanding main image", uuid)
	}

	//embedded images
	var emImagesUUIDs []imageSetUUID
	body, foundBody := cc[bodyXML]
	if foundBody {
		bodyXML := body.(string)
		emImagesUUIDs, err = ir.parser.GetEmbedded(bodyXML)
		if err != nil {
			logger.Infof(uuid, errors.Wrapf(err, "Cannot parse body for uuid=%s", uuid).Error())
		} else {
			b.embeddedImages = emImagesUUIDs
		}
	} else {
		logger.Infof(uuid, "Missing body for %v.Skipping expanding embedded images.", uuid)
	}

	//promotional image
	var altImgMap map[string]interface{}
	var foundPromImg bool
	altImg, found := cc[altImages]
	if found {
		var promImg interface{}
		altImgMap = altImg.(map[string]interface{})
		promImg, foundPromImg = altImgMap[promoImage]
		if foundPromImg {
			promImgID := promImg.(string)
			b.promotionalImage = extractUUIDFromURL(promImgID)
		} else {
			logger.Infof(uuid, "Cannot find promotional image for %v. Skipping expanding promotional image", uuid)
		}
	}

	if !foundMainImg && !foundBody && !foundPromImg {
		return c, errors.Errorf("Cannot read supplied content %v. Nothing to expand.", uuid)
	}

	imgMap, err := ir.reader.Get(b)
	if err != nil {
		return c, errors.Wrapf(err, "Error while getting expanded images for uuid:%v", uuid)
	}

	if foundMainImg {
		cc[mainImage] = ir.resolveImageSet(b.mainImageSet, imgMap)
	}
	if foundBody && len(emImagesUUIDs) > 0 {
		cc[embeds] = ir.resolveImageSetArray(b.embeddedImages, imgMap)
	}
	if foundPromImg {
		promImgContent := imgMap[b.promotionalImage]
		altImgMap[promoImage] = promImgContent
	}

	return cc, nil
}

func (ir *ImageResolver) UnrollLeadImages(c Content) (Content, error) {
	uuid := extractUUIDFromURL(c[id].(string))
	cc := c.clone()
	b := UUIDBatch{}
	var err error

	images, foundLeadImages := cc[leadImages].([]interface{})
	if !foundLeadImages {
		return c, errors.Errorf("Nothing to expand for the supplied content %s", uuid)
	}

	var lis []map[string]interface{}
	for _, item := range images {
		li := item.(map[string]interface{})
		lis = append(lis, li)
		uuid := extractUUIDFromURL(li[id].(string))
		li[image] = uuid
		b.leadImages = append(b.leadImages, uuid)
	}

	if foundLeadImages && len(b.leadImages) <= 0 {
		return c, errors.Errorf("Cannot read UUIDs for lead images of %v. Returning same content.", uuid)
	}

	imgMap, err := ir.reader.Get(b)
	if err != nil {
		return c, errors.Wrapf(err, "Error while getting content for expanded images uuid:%v", uuid)
	}

	for _, li := range lis {
		uuid := li[image].(string)
		imageData, found := ir.resolveContent(uuid, imgMap)
		if !found {
			logger.Infof(uuid, "Missing image model %s. Returning only de id.", uuid)
		}
		li[image] = imageData
	}

	cc[leadImages] = lis
	return cc, nil
}

func (ir *ImageResolver) resolveImageSetArray(imgArray []imageSetUUID, imgMap map[string]Content) []interface{} {
	var imgSets []interface{}
	for _, imgUUID := range imgArray {
		imgSet := ir.resolveImageSet(imgUUID, imgMap)
		imgSets = append(imgSets, imgSet)
	}
	return imgSets
}

func (ir *ImageResolver) resolveImageSet(img imageSetUUID, imgMap map[string]Content) Content {
	imageSet, found := ir.resolveContent(img.uuid, imgMap)
	var imageModel Content
	if found {
		imageModel, _ = ir.resolveContent(img.imageModelUUID, imgMap)
		imageSet[members] = []Content{imageModel}
	}
	return imageSet
}

func (ir *ImageResolver) resolveContent(uuid string, imgMap map[string]Content) (Content, bool) {
	c, found := imgMap[uuid]
	if !found {
		return Content{id: createID(ir.apiHost, "content", uuid)}, false
	}
	return c, true
}

func (c Content) clone() Content {
	clone := make(Content)
	for k, v := range c {
		clone[k] = v
	}
	return clone
}

func (u UUIDBatch) toArray() (UUIDs []string) {
	UUIDs = append(UUIDs, u.mainImageSet.uuid)
	UUIDs = append(UUIDs, u.mainImageSet.imageModelUUID)
	for _, e := range u.embeddedImages {
		UUIDs = append(UUIDs, e.uuid)
		UUIDs = append(UUIDs, e.imageModelUUID)
	}
	UUIDs = append(UUIDs, u.promotionalImage)
	for _, l := range u.leadImages {
		UUIDs = append(UUIDs, l)
	}
	return UUIDs
}
