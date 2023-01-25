package entities

import (
	"errors"
	"regexp"
)

type Reduction struct {
	LongUrl    string `json:"long_url"`
	CustomText string `json:"custom_text"`
}

func (r Reduction) Validate() error {
	re, err := regexp.Compile("^(http:\\/\\/www\\.|https:\\/\\/www\\.|http:\\/\\/|https:\\/\\/|\\/|\\/\\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\\-\\.]{1}[a-z0-9]+)*\\.[a-z]{2,5}(:[0-9]{1,5})?(\\/.*)?$")
	if err != nil {
		return err
	}

	matched := re.MatchString(r.LongUrl)
	if !matched {
		return errors.New("incorrect url")
	}

	re, err = regexp.Compile("(^[a-zA-Z]+(-)?[a-zA-Z]+)?")
	if err != nil {
		return err
	}

	matched = re.MatchString(r.CustomText)
	if !matched {
		return errors.New("unsupported custom text")
	}

	return nil
}
