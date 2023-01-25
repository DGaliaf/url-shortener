package reduction

type CreateShortUrlDTO struct {
	LongUrl    string `json:"long_url"`
	CustomText string `json:"custom_text,omitempty"`
}

type GetLongUrlDTO struct {
	HashedLink string `json:"hashed_link"`
}
