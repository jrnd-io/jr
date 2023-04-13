package registry

var typeMap = make(map[string]interface{})

func Register(name string, t interface{}) {
	typeMap[name] = t
}

func GetType(name string) interface{} {
	return typeMap[name]
}
