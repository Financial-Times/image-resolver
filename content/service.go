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
	UnrollImages(req UnrollRequest) UnrollResponse
	UnrollLeadImages(req UnrollRequest) UnrollResponse
}

type ImageResolver struct {
	reader    Reader
	whitelist string
	apiHost   string
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

func NewImageResolver(r Reader, whitelist string, apiHost string) *ImageResolver {
	return &ImageResolver{
		reader:    r,
		whitelist: whitelist,
		apiHost:   apiHost,
	}
}

func (ir *ImageResolver) UnrollImages(req UnrollRequest) UnrollResponse {
	//make a copy of the content
	c := req.c
	cc := c.clone()
	b := UUIDBatch{}

	//mainImage
	mi, foundMainImg := cc[mainImage].(map[string]interface{})
	if foundMainImg {
		id := mi[id].(string)
		var mis imageSetUUID
		mis.uuid = extractUUIDFromURL(id)
		imgModelUUID, err := getImageModelUUID(mis.uuid)
		if err != nil {
			logger.Infof(req.tid, req.uuid, "Cannot convert image set UUID %v to a image model UUID. Skipping item.", mis.uuid)
		} else {
			mis.imageModelUUID = imgModelUUID
		}
		b.mainImageSet = mis
	} else {
		logger.Infof(req.tid, req.uuid, "Cannot find main image for %v. Skipping expanding main image", req.uuid)
	}

	//embedded images
	var emImagesUUIDs []imageSetUUID
	var err error
	body, foundBody := cc[bodyXML]
	if foundBody {
		bodyXML := body.(string)
		emImagesUUIDs, err = getEmbedded(bodyXML, ir.whitelist, req.tid, req.uuid)
		if err != nil {
			logger.Infof(req.tid, req.uuid, errors.Wrapf(err, "Cannot parse body for uuid=%s", req.uuid).Error())
		} else {
			b.embeddedImages = emImagesUUIDs
		}
	} else {
		logger.Infof(req.tid, req.uuid, "Missing body for %v.Skipping expanding embedded images.", req.uuid)
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
			logger.Infof(req.tid, req.uuid, "Cannot find promotional image for %v. Skipping expanding promotional image", req.uuid)
		}
	}

	if !foundMainImg && !foundBody && !foundPromImg {
		return UnrollResponse{c, errors.Errorf("Cannot read supplied content %v. Nothing to expand.", req.uuid)}
	}

	imgMap, err := ir.reader.Get(b)
	if err != nil {
		return UnrollResponse{c, errors.Wrapf(err, "Error while getting expanded images for uuid:%v", req.uuid)}
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

	return UnrollResponse{cc, nil}
}

func (ir *ImageResolver) UnrollLeadImages(req UnrollRequest) UnrollResponse {
	cc := req.c.clone()
	b := UUIDBatch{}
	var err error

	images, foundLeadImages := cc[leadImages].([]interface{})
	if !foundLeadImages {
		return UnrollResponse{req.c, errors.Errorf("Nothing to expand for the supplied content %s", req.uuid)}
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
		return UnrollResponse{req.c, errors.Errorf("Cannot read UUIDs for lead images of %v. Returning same content.", req.uuid)}
	}

	imgMap, err := ir.reader.Get(b)
	if err != nil {
		return UnrollResponse{req.c, errors.Wrapf(err, "Error while getting content for expanded images uuid:%v", req.uuid)}
	}

	for _, li := range lis {
		uuid := li[image].(string)
		imageData, found := ir.resolveContent(uuid, imgMap)
		if !found {
			logger.Infof(req.tid, req.uuid, "Missing image model %s. Returning only de id.", uuid)
		}
		li[image] = imageData
	}

	cc[leadImages] = lis
	return UnrollResponse{cc, nil}
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
