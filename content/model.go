package content

type Content map[string]interface{}

type PromotionalImage struct {
	PromotionalImage  interface {}`json:"promotionalImage,omitempty"`
}

type ImageOutput struct {
	Content  interface {} `json:"image,omitemptyt"`
	Type string      `json:"type,omitempty"`
}