package option

type Options struct {
	File                    string
	Output                  string
	LogFormat               string
	QueryString             bool
	QueryStringIgnoreValues bool
	MatchingGroups          []string
	Filters                 []string
	KeyOption               KeyOption
}

type KeyOption struct {
	UriKey    string
	MethodKey string
	UserKey   string
}
