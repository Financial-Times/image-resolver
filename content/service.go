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

type Unroller interface {
	UnrollContent(req UnrollEvent) UnrollResult
	UnrollInternalContent(req UnrollEvent) UnrollResult
}

type ContentUnroller struct {
	reader    Reader
	whitelist string
	apiHost   string
}

type Content map[string]interface{}

type ContentSchema map[string][]string

func NewContentUnroller(r Reader, whitelist string, apiHost string) *ContentUnroller {
	return &ContentUnroller{
		reader:    r,
		whitelist: whitelist,
		apiHost:   apiHost,
	}
}

func (u *ContentUnroller) UnrollContent(req UnrollEvent) UnrollResult {
	//make a copy of the content
	cc := req.c.clone()

	//mainImage
	schema := make(ContentSchema)
	mi, foundMainImg := cc[mainImage].(map[string]interface{})
	if foundMainImg {
		u, err := extractUUIDFromString(mi[id].(string))
		if err != nil {
			logger.Infof(req.tid, req.uuid, "Cannot find main image: %v. Skipping expanding main image", err.Error())
			foundMainImg = false
		} else {
			schema.put(mainImage, u)
		}
	} else {
		logger.Info(req.tid, req.uuid, "Cannot find main image. Skipping expanding main image")
	}

	//embedded - images and dynamic content
	emContentUUIDs, foundEmbedded := u.extractEmbeddedContentByType(cc, u.whitelist, req.tid, req.uuid)
	if foundEmbedded {
		schema.putAll(embeds, emContentUUIDs)
	}

	//promotional image
	var foundPromImg bool
	altImg, found := cc[altImages].(map[string]interface{})
	if found {
		var promImg map[string]interface{}
		promImg, foundPromImg = altImg[promotionalImage].(map[string]interface{})
		if foundPromImg {
			if id, ok := promImg[id].(string); ok {
				u, err := extractUUIDFromString(id)
				if err != nil {
					logger.Infof(req.tid, req.uuid, "Cannot find promotional image: %v. Skipping expanding promotional image", err.Error())
					foundPromImg = false
				} else {
					schema.put(promotionalImage, u)
				}
			} else {
				logger.Info(req.tid, req.uuid, "Promotional image is missing the id field. Skipping expanding promotional image")
				foundPromImg = false
			}
		} else {
			logger.Info(req.tid, req.uuid, "Cannot find promotional image. Skipping expanding promotional image")
		}
	}

	if !foundMainImg && !foundEmbedded && !foundPromImg {
		logger.Infof(req.tid, req.uuid, "No main image or body images or promotional image to expand for supplied content %s", req.uuid)
		return UnrollResult{req.c, nil}
	}

	contentMap, err := u.reader.Get(schema.toArray(), req.tid)
	if err != nil {
		return UnrollResult{req.c, errors.Wrapf(err, "Error while getting expanded images for uuid:%v", req.uuid)}
	}
	u.resolveModelsForSetsMembers(schema, contentMap, req.tid, req.tid)

	if foundMainImg {
		cc[mainImage] = contentMap[schema.get(mainImage)]
	}

	embeddedContent := schema.getAll(embeds)
	if foundEmbedded && len(embeddedContent) > 0 {
		embedded := []Content{}
		for _, eis := range embeddedContent {
			embedded = append(embedded, contentMap[eis])
		}
		cc[embeds] = embedded
	}

	if foundPromImg {
		pi, found := contentMap[schema.get(promotionalImage)]
		if found {
			cc[altImages].(map[string]interface{})[promotionalImage] = pi
		}
	}

	return UnrollResult{cc, nil}
}

func (u *ContentUnroller) UnrollInternalContent(req UnrollEvent) UnrollResult {
	cc := req.c.clone()
	expLeadImages, foundImages := u.unrollLeadImages(cc, req.tid, req.uuid)
	if foundImages {
		cc[leadImages] = expLeadImages
	}

	embedded, foundEmbedded := u.unrollEmbeddedDynamicContent(cc, req.tid, req.uuid)
	if foundEmbedded {
		cc[embeds] = embedded
	}

	return UnrollResult{cc, nil}
}

func (u *ContentUnroller) unrollLeadImages(cc Content, tid string, uuid string) ([]Content, bool) {
	images, foundLeadImages := cc[leadImages].([]interface{})
	if !foundLeadImages {
		logger.Info(tid, uuid, "No lead images to expand for supplied content")
		return nil, false
	}

	if len(images) == 0 {
		logger.Info(tid, uuid, "No lead images to expand for supplied content")
		return nil, false
	}
	schema := make(ContentSchema)
	for _, item := range images {
		li := item.(map[string]interface{})
		uuid, err := extractUUIDFromString(li[id].(string))
		if err != nil {
			logger.Infof(tid, uuid, "Error while getting UUID for %s: %v", li[id].(string), err.Error())
			continue
		}
		li[image] = uuid
		schema.put(leadImages, uuid)
	}

	imgMap, err := u.reader.Get(schema.toArray(), tid)
	if err != nil {
		logger.Errorf(tid, uuid, errors.Wrapf(err, "Error while getting content for expanded images uuid"))
		return nil, false
	}

	expLeadImages := []Content{}
	for _, li := range images {
		rawLi := li.(map[string]interface{})
		rawLiUUID := rawLi[image].(string)
		liContent := fromMap(rawLi)
		imageData, found := u.resolveContent(rawLiUUID, imgMap)
		if !found {
			logger.Infof(tid, uuid, "Missing image model %s. Returning only de id.", rawLiUUID)
			delete(liContent, image)
			expLeadImages = append(expLeadImages, liContent)
			continue
		}
		liContent[image] = imageData
		expLeadImages = append(expLeadImages, liContent)
	}

	cc[leadImages] = expLeadImages
	return expLeadImages, true
}

func (u *ContentUnroller) unrollEmbeddedDynamicContent(cc Content, tid string, uuid string) ([]Content, bool) {
	emContentUUIDs, foundEmbedded := u.extractEmbeddedContentByType(cc, "^http://www.ft.com/ontology/content/DynamicContent", tid, uuid)
	if !foundEmbedded {
		return nil, false
	}

	contentMap, err := u.reader.GetInternal(emContentUUIDs, tid)
	if err != nil {
		logger.Errorf(tid, uuid, errors.Wrapf(err, "Error while getting embedded dynamic content"))
		return nil, false
	}

	embedded := []Content{}
	for _, ec := range emContentUUIDs {
		embedded = append(embedded, contentMap[ec])
	}

	return embedded, true
}

func (u *ContentUnroller) resolveModelsForSetsMembers(b ContentSchema, imgMap map[string]Content, tid string, uuid string) {
	mainImageUUID := b.get(mainImage)
	u.resolveImageSet(mainImageUUID, imgMap, tid, uuid)
	for _, embeddedImgSet := range b.getAll(embeds) {
		u.resolveImageSet(embeddedImgSet, imgMap, tid, uuid)
	}
}

func (u *ContentUnroller) resolveImageSet(imageSetUUID string, imgMap map[string]Content, tid string, uuid string) {
	imageSet, found := u.resolveContent(imageSetUUID, imgMap)
	if !found {
		imgMap[imageSetUUID] = Content{id: createID(u.apiHost, "content", imageSetUUID)}
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
			mUUID, err := extractUUIDFromString(mID)
			if err != nil {
				logger.Infof(tid, uuid, "Error while extracting UUID from %s: %v", mID, err.Error())
				continue
			}
			mContent, found := u.resolveContent(mUUID, imgMap)
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

func (u *ContentUnroller) resolveContent(uuid string, imgMap map[string]Content) (Content, bool) {
	c, found := imgMap[uuid]
	if !found {
		return Content{}, false
	}
	return c, true
}

func (u *ContentUnroller) extractEmbeddedContentByType(cc Content, acceptedType string, tid string, uuid string) ([]string, bool) {
	body, foundBody := cc[bodyXML]
	if !foundBody {
		logger.Info(tid, uuid, "Missing body. Skipping expanding embedded content and images.")
		return nil, false
	}

	bodyXML := body.(string)
	emContentUUIDs, err := getEmbedded(bodyXML, acceptedType, tid, uuid)
	if err != nil {
		logger.Errorf(tid, uuid, errors.Wrapf(err, "Cannot parse body for uuid=%s", uuid))
		return nil, false
	}

	if len(emContentUUIDs) == 0 {
		return nil, false
	}

	return emContentUUIDs, true
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
		u, err := extractUUIDFromString(url)
		if err != nil {
			continue
		}
		uuids = append(uuids, u)
	}
	return uuids
}

func (c Content) merge(src Content) {
	for k, v := range src {
		c[k] = v
	}
}

func (u ContentSchema) put(key string, value string) {
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

func (u ContentSchema) get(key string) string {
	if _, found := u[key]; key != mainImage && key != promotionalImage || !found {
		return ""
	}
	return u[key][0]
}

func (u ContentSchema) putAll(key string, values []string) {
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

func (u ContentSchema) getAll(key string) []string {
	if key != embeds && key != leadImages {
		return []string{}
	}
	return u[key]
}

func (u ContentSchema) toArray() (UUIDs []string) {
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
