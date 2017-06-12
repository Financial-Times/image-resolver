package content

import (
	"github.com/pkg/errors"
)

const (
	mainImage        = "mainImage"
	id               = "id"
	embeds           = "embeds"
	altImages        = "alternativeImages"
	leadImages       = "leadImages"
	members          = "members"
	bodyXML          = "bodyXML"
	promotionalImage = "promotionalImage"
	image            = "image"
)

type Resolver interface {
	UnrollImages(req UnrollEvent) UnrollResult
	UnrollLeadImages(req UnrollEvent) UnrollResult
}

type ImageResolver struct {
	reader    Reader
	whitelist string
	apiHost   string
}

type Content map[string]interface{}

type ImageSchema map[string][]string

func NewImageResolver(r Reader, whitelist string, apiHost string) *ImageResolver {
	return &ImageResolver{
		reader:    r,
		whitelist: whitelist,
		apiHost:   apiHost,
	}
}

func (ir *ImageResolver) UnrollImages(req UnrollEvent) UnrollResult {
	//make a copy of the content
	cc := req.c.clone()

	//mainImage
	is := make(ImageSchema)
	mi, foundMainImg := cc[mainImage].(map[string]interface{})
	if foundMainImg {
		is.put(mainImage, extractUUIDFromString(mi[id].(string)))
	} else {
		logger.Infof(req.tid, req.uuid, "Cannot find main image for %v. Skipping expanding main image", req.uuid)
	}

	//embedded images
	body, foundBody := cc[bodyXML]
	if foundBody {
		bodyXML := body.(string)
		emImagesUUIDs, err := getEmbedded(bodyXML, ir.whitelist, req.tid, req.uuid)
		if err != nil {
			logger.Infof(req.tid, req.uuid, errors.Wrapf(err, "Cannot parse body for uuid=%s", req.uuid).Error())
		} else {
			is.putAll(embeds, emImagesUUIDs)
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
		promImg, foundPromImg = altImgMap[promotionalImage]
		if foundPromImg {
			promImgID := promImg.(string)
			is.put(promotionalImage, extractUUIDFromString(promImgID))
		} else {
			logger.Infof(req.tid, req.uuid, "Cannot find promotional image for %v. Skipping expanding promotional image", req.uuid)
		}
	}

	if !foundMainImg && !foundBody && !foundPromImg {
		return UnrollResult{req.c, errors.Errorf("Cannot read supplied content %v. Nothing to expand.", req.uuid)}
	}

	imgMap, err := ir.reader.Get(is.toArray())
	if err != nil {
		return UnrollResult{req.c, errors.Wrapf(err, "Error while getting expanded images for uuid:%v", req.uuid)}
	}
	ir.resolveModelsForSetsMembers(is, imgMap)

	if foundMainImg {
		cc[mainImage] = imgMap[is.get(mainImage)]
	}

	embeddedImgSets := is.getAll(embeds)
	if foundBody && len(embeddedImgSets) > 0 {
		embedded := []Content{}
		for _, eis := range embeddedImgSets {
			embedded = append(embedded, imgMap[eis])
		}
		cc[embeds] = embedded
	}

	if foundPromImg {
		altImgMap[promotionalImage] = imgMap[is.get(promotionalImage)]
	}

	return UnrollResult{cc, nil}
}

func (ir *ImageResolver) UnrollLeadImages(req UnrollEvent) UnrollResult {
	cc := req.c.clone()
	images, foundLeadImages := cc[leadImages].([]interface{})
	if !foundLeadImages {
		return UnrollResult{req.c, errors.Errorf("Nothing to expand for the supplied content %s", req.uuid)}
	}

	b := make(ImageSchema)
	for _, item := range images {
		li := item.(map[string]interface{})
		uuid := extractUUIDFromString(li[id].(string))
		li[image] = uuid
		b.put(leadImages, uuid)
	}

	imgMap, err := ir.reader.Get(b.toArray())
	if err != nil {
		return UnrollResult{req.c, errors.Wrapf(err, "Error while getting content for expanded images uuid:%v", req.uuid)}
	}

	expLeadImages := []Content{}
	for _, li := range images {
		rawLi := li.(map[string]interface{})
		uuid := rawLi[image].(string)
		liContent := fromMap(rawLi)
		imageData, found := ir.resolveContent(uuid, imgMap)
		if !found {
			logger.Infof(req.tid, req.uuid, "Missing image model %s. Returning only de id.", uuid)
			delete(liContent, image)
			expLeadImages = append(expLeadImages, liContent)
			continue
		}
		liContent[image] = imageData
		expLeadImages = append(expLeadImages, liContent)
	}

	cc[leadImages] = expLeadImages
	return UnrollResult{cc, nil}
}

func (ir *ImageResolver) resolveModelsForSetsMembers(b ImageSchema, imgMap map[string]Content) {
	mainImageUUID := b.get(mainImage)
	ir.resolveImageSet(mainImageUUID, imgMap)
	for _, embeddedImgSet := range b.getAll(embeds) {
		ir.resolveImageSet(embeddedImgSet, imgMap)
	}
}

func (ir *ImageResolver) resolveImageSet(imageSetUUID string, imgMap map[string]Content) {
	imageSet, found := ir.resolveContent(imageSetUUID, imgMap)
	if !found {
		imgMap[imageSetUUID] = Content{id: createID(ir.apiHost, "content", imageSetUUID)}
		return
	}

	rawMembers, found := imageSet[members]
	if found {
		membList, ok := rawMembers.([]interface{})
		if !ok {
			return
		}

		expMembers := []Content{}
		for _, m := range membList {
			mData := fromMap(m.(map[string]interface{}))
			mID := mData[id].(string)
			mContent, found := ir.resolveContent(extractUUIDFromString(mID), imgMap)
			if !found {
				expMembers = append(expMembers, mData)
				continue
			}
			mData.merge(mContent)
			expMembers = append(expMembers, mData)
		}
		imageSet[members] = expMembers
	}

}

func (ir *ImageResolver) resolveContent(uuid string, imgMap map[string]Content) (Content, bool) {
	c, found := imgMap[uuid]
	if !found {
		return Content{}, false
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

func (c Content) getMembersUUID() []string {
	uuids := []string{}
	members, found := c[members]
	if !found {
		return uuids
	}

	memList, ok := members.([]interface{})
	if !ok {
		return uuids
	}
	for _, m := range memList {
		mData := m.(map[string]interface{})
		url, found := mData[id].(string)
		if !found {
			continue
		}
		uuids = append(uuids, extractUUIDFromString(url))
	}
	return uuids
}

func (c Content) merge(src Content) {
	for k, v := range src {
		c[k] = v
	}
}

func (u ImageSchema) put(key string, value string) {
	if key != mainImage && key != promotionalImage && key != leadImages {
		return
	}
	prev, found := u[key]
	if !found {
		u[key] = []string{value}
		return
	}
	act := append(prev, value)
	u[key] = act
}

func (u ImageSchema) get(key string) string {
	if _, found := u[key]; key != mainImage && key != promotionalImage || !found {
		return ""
	}
	return u[key][0]
}

func (u ImageSchema) putAll(key string, values []string) {
	if key != embeds && key != leadImages {
		return
	}
	prevValue, found := u[key]
	if !found {
		u[key] = values
		return
	}
	u[key] = append(prevValue, values...)
}

func (u ImageSchema) getAll(key string) []string {
	if key != embeds && key != leadImages {
		return []string{}
	}
	return u[key]
}

func (u ImageSchema) toArray() (UUIDs []string) {
	for _, v := range u {
		UUIDs = append(UUIDs, v...)
	}
	return UUIDs
}

func fromMap(src map[string]interface{}) Content {
	dest := Content{}
	for k, v := range src {
		dest[k] = v
	}
	return dest
}
