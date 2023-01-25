package reduction

type CreateShortUrlDTO struct {
	LongUrl string `json:"long_url"`
}

type GetLongUrlDTO struct {
	HashedLink string `json:"hashed_link"`
}
