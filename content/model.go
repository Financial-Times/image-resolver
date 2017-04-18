package content

import "time"

type Identifier struct {
	Authority       string `json:"authority"`
	IdentifierValue string `json:"identifierValue"`
}

type Annotation struct {
	Predicate string `json:"predicate"`
	Uri       string `json:"uri"`
	ApiUrl    string `json:"apiUrl"`
	Label     string `json:"label"`
	Type      string `json:"type"`
}
type Standout struct {
	EditorsChoice bool `json:"editorsChoice"`
	Exclusive     bool `json:"exclusive"`
	Scoop         bool `json:"scoop"`
}

type AlternativeTitles struct {
	PromotionalTitle string `json:"promotionalTitle"`
}

type AlternativeStandfirsts struct {
	PromotionalStandfirst string `json:"promotionalStandfirst"`
}

type AlternativeImages struct {
	PromotionalImage string `json:"promotionalImage"`
}

type Uri struct {
	Id string `json:"id"`
}

type Comments struct {
	Enabled bool `json:"enabled"`
}

type Copyright struct {
	Notice string `json:"notice"`
}

type Content struct {
	Uuid               string `json:"id"`
	Type               string `json:"type"`
	BodyXML            string `json:"bodyXML"`
	OpeningXML         string `json:"openingXML"`
	Title              string `json:"title"`
	Byline             string `json:"byline"`
	Description        string `json:"description"`
	PublishedDate      time.Time `json:"publishedDate"`
	Identifiers        []Identifier `json:"identifiers"`
	Members            []Uri `json:"members"`
	RequestUrl         string `json:"requestUrl"`
	BinaryUrl          string `json:"binaryUrl"`
	Brands             []string `json:"brands"`
	Annotations        []Annotation `json:"annotations"`
	MainImage          Uri `json:"MainImage"`
	Comments           Comments `json:"comments"`
	Realtime           bool `json:"realtime"`
	Copyright          Copyright `json:"copyright"`
	PublishReference   string `json:"publishReference"`
	PixelWidth         int `json:"pixelWidth"`
	PixelHeight        int `json:"pixelHeight"`
	Stdout             Standout `json:"standout"`
	LastModified       time.Time `json:"lastModified"`
	Standfirst         string `json:"standfirst"`
	AltTitles          AlternativeTitles `json:"alternativeTitles"`
	AltStandfirsts     AlternativeStandfirsts `json:"alternativeStandfirsts"`
	AltImages          AlternativeImages `json:"alternativeImages"`
	WebUrl             string `json:"webUrl"`
	CanBeSyndicated    string `json:"canBeSyndicated"`
	FirstPublishedDate time.Time `json:"firstPublishedDate"`
	AccessLevel        string `json:"accessLevel"`
	CanBeDistributed   string `json:"canBeDistributed"`
}

type ContentOutput struct {
	Uuid               string `json:"id"`
	Type               string `json:"type"`
	BodyXML            string `json:"bodyXML"`
	OpeningXML         string `json:"openingXML"`
	Title              string `json:"title"`
	Byline             string `json:"byline"`
	Description        string `json:"description"`
	PublishedDate      time.Time `json:"publishedDate"`
	Identifiers        []Identifier `json:"identifiers"`
	Members            []Content `json:"members"`
	RequestUrl         string `json:"requestUrl"`
	BinaryUrl          string `json:"binaryUrl"`
	Brands             []string `json:"brands"`
	Annotations        []Annotation `json:"annotations"`
	MainImage          Uri `json:"MainImage"`
	Comments           Comments `json:"comments"`
	Realtime           bool `json:"realtime"`
	Copyright          Copyright `json:"copyright"`
	PublishReference   string `json:"publishReference"`
	PixelWidth         int `json:"pixelWidth"`
	PixelHeight        int `json:"pixelHeight"`
	Stdout             Standout `json:"standout"`
	LastModified       time.Time `json:"lastModified"`
	Standfirst         string `json:"standfirst"`
	AltTitles          AlternativeTitles `json:"alternativeTitles"`
	AltStandfirsts     AlternativeStandfirsts `json:"alternativeStandfirsts"`
	AltImages          AlternativeImages `json:"alternativeImages"`
	WebUrl             string `json:"webUrl"`
	CanBeSyndicated    string `json:"canBeSyndicated"`
	FirstPublishedDate time.Time `json:"firstPublishedDate"`
	AccessLevel        string `json:"accessLevel"`
	CanBeDistributed   string `json:"canBeDistributed"`
}

type UnrolledContent struct {
	MainImage        *ContentOutput `json:"mainImage,omitempty"`
	Embeds           []ContentOutput `json:"embeds,omitempty"`
	PromotionalImage *ContentOutput  `json:"promotionalImage,omitempty"`
}

type UnrolledLeadImagesContent struct {
}
