package jr

type Context struct {
	HowMany      int
	Frequency    int
	Localization []string
	Seed         int64
}

func NewContext(howMany int, frequency int, locales []string, seed int64) *Context {
	context := Context{HowMany: howMany, Frequency: frequency, Localization: locales, Seed: seed}
	return &context
}
