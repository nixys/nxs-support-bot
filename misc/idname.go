package misc

type IDName struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type IDNameLocale struct {
	ID   int64             `json:"id"`
	Name map[string]string `json:"name"`
}

const OnPageDefault = 5

const IDNameLocaleDefaultLang = "default"

func IDNamePaginate(elts []IDName, page, limit int64) (s []IDName, isBack, isNext bool) {

	len := int64(len(elts))
	if len == 0 {
		return nil, false, false
	}

	f := page * limit
	if f >= len {
		return nil, f > 0, false
	}

	t := (page + 1) * limit
	if t > len {
		t = len
	}

	return elts[f:t], f > 0, len > t
}

func (idn *IDNameLocale) ValueGet(lang string) *string {

	if idn == nil {
		return nil
	}

	n, b := idn.Name[lang]
	if b == false {
		s := idn.Name[IDNameLocaleDefaultLang]
		return &s
	}
	return &n
}
