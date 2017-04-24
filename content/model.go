package content

import "time"

type Identifier struct {
	Authority       string `json:"authority,omitempty"`
	IdentifierValue string `json:"identifierValue,omitempty"`
}

type Standout struct {
	EditorsChoice bool `json:"editorsChoice,omitempty"`
	Exclusive     bool `json:"exclusive,omitempty"`
	Scoop         bool `json:"scoop,omitempty"`
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
	Type               string `json:"type,omitempty"`
	BodyXML            string `json:"bodyXML,omitempty"`
	OpeningXML         string `json:"openingXML,omitempty"`
	Title              string `json:"title,omitempty"`
	Byline             string `json:"byline,omitempty"`
	Description        string `json:"description,omitempty"`
	PublishedDate      *time.Time `json:"publishedDate,omitempty"`
	Identifiers        []Identifier `json:"identifiers,omitempty"`
	Members            []Content `json:"members,omitempty"`
	RequestURL         string `json:"requestUrl,omitempty"`
	BinaryURL          string `json:"binaryUrl,omitempty"`
	Brands             []string `json:"brands,omitempty"`
	MainImage          *Uri `json:"MainImage,omitempty"`
	Comments           *Comments `json:"comments,omitempty"`
	Realtime           bool `json:"realtime,omitempty"`
	Copyright          *Copyright `json:"copyright,omitempty"`
	PublishReference   string `json:"publishReference,omitempty"`
	PixelWidth         int `json:"pixelWidth,omitempty"`
	PixelHeight        int `json:"pixelHeight,omitempty"`
	Stdout             *Standout `json:"standout,omitempty"`
	LastModified       *time.Time `json:"lastModified,omitempty"`
	Standfirst         string `json:"standfirst,omitempty"`
	AltTitles          *AlternativeTitles `json:"alternativeTitles,omitempty"`
	AltStandfirsts     *AlternativeStandfirsts `json:"alternativeStandfirsts,omitempty"`
	AltImages          *AlternativeImages `json:"alternativeImages,omitempty"`
	WebURL             string `json:"webUrl,omitempty"`
	CanBeSyndicated    string `json:"canBeSyndicated,omitempty"`
	FirstPublishedDate *time.Time `json:"firstPublishedDate,omitempty"`
	AccessLevel        string `json:"accessLevel,omitempty"`
	CanBeDistributed   string `json:"canBeDistributed,omitempty"`
}

type UnrolledContent struct {
	MainImage         *Content `json:"mainImage,omitempty"`
	Embeds            []Content `json:"embeds,omitempty"`
	AlternativeImages *PromotionalImage  `json:"alternativeImages,omitempty"`
}

type PromotionalImage struct {
	PromotionalImage *Content `json:"promotionalImage,omitempty"`
}

