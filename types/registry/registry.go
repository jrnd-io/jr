package registry

var typeMap map[string]interface{}

func RegisterType(name string, t interface{}) {
	typeMap[name] = t
}

func GetType(name string) interface{} {
	return typeMap[name]
}
