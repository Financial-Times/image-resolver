package content


type Identifier struct {
	Authority       string `json:"authority,omitempty"`
	IdentifierValue string `json:"identifierValue,omitempty"`
}

type Standout struct {
	EditorsChoice bool `json:"editorsChoice"`
	Exclusive     bool `json:"exclusive"`
	Scoop         bool `json:"scoop"`
}

type AlternativeTitles struct {
	PromotionalTitle string `json:"promotionalTitle,omitempty"`
}

type AlternativeStandfirsts struct {
	PromotionalStandfirst string `json:"promotionalStandfirst,omitempty"`
}

type AlternativeImages struct {
	PromotionalImage string `json:"promotionalImage,omitempty"`
}

type Uri struct {
	Id string `json:"id,omitempty"`
}

type Comments struct {
	Enabled bool `json:"enabled,omitempty"`
}

type Copyright struct {
	Notice string `json:"notice,omitempty"`
}

type Content struct {
	UUID               string `json:"id,omitempty"`
	Type               *string `json:"type,omitempty"`
	BodyXML            *string `json:"bodyXML,omitempty"`
	OpeningXML         *string `json:"openingXML,omitempty"`
	Title              *string `json:"title,omitempty"`
	AltTitles          *AlternativeTitles `json:"alternativeTitles,omitempty"`
	Standfirst         *string `json:"standfirst,omitempty"`
	AltStandfirsts     *AlternativeStandfirsts `json:"alternativeStandfirsts,omitempty"`
	Byline             *string `json:"byline,omitempty"`
	Description        *string `json:"description,omitempty"`
	FirstPublishedDate interface{} `json:"firstPublishedDate,omitempty"`
	PublishedDate      interface{}  `json:"publishedDate,omitempty"`
	WebURL             *string `json:"webUrl,omitempty"`
	Identifiers        []Identifier `json:"identifiers,omitempty"`
	Members            []Content `json:"members,omitempty"`
	RequestURL         *string `json:"requestUrl,omitempty"`
	BinaryURL          *string `json:"binaryUrl,omitempty"`
	PixelWidth         int `json:"pixelWidth,omitempty"`
	PixelHeight        int `json:"pixelHeight,omitempty"`
	Brands             []string `json:"brands,omitempty"`
	MainImage          *Content `json:"mainImage,omitempty"`
	AltImages          interface{} `json:"alternativeImages,omitempty"`
	Comments           *Comments `json:"comments,omitempty"`
	Copyright          *Copyright `json:"copyright,omitempty"`
	Realtime           bool `json:"realtime,omitempty"`
	Stdout             *Standout `json:"standout,omitempty"`
	PublishReference   *string `json:"publishReference,omitempty"`
	LastModified       interface{}  `json:"lastModified,omitempty"`
	CanBeSyndicated    string `json:"canBeSyndicated,omitempty"`
	AccessLevel        string `json:"accessLevel,omitempty"`
	CanBeDistributed   string `json:"canBeDistributed,omitempty"`
	Embeds            []Content `json:"embeds,omitempty"`
	LeadImages        interface{} `json:"leadImages,omitempty"`
}

type UnrolledContent struct {
	MainImage         *Content `json:"mainImage,omitempty"`
	Embeds            []Content `json:"embeds,omitempty"`
	AlternativeImages *PromotionalImage  `json:"alternativeImages,omitempty"`
}

type PromotionalImage struct {
	PromotionalImage *Content `json:"promotionalImage,omitempty"`
}

type Image struct {
	Id  string   `json:"id,omitempty"`
	Type string  `json:"type,omitempty"`
}

type ImageOutput struct {
	Content  interface {} `json:"image,omitemptyt"`
	Type string      `json:"type,omitempty"`
}