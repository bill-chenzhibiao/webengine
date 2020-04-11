package webengine

type methodPathMapping struct{
	method string
	pathMapping map[string]HandlersChain
}
type methodPathMappings []methodPathMapping

func (methodPathMappings methodPathMappings) get(method string) map[string]HandlersChain{
	for _, pathMappingItem := range methodPathMappings{
		if pathMappingItem.method == method{
			return pathMappingItem.pathMapping
		}
	}
	return nil
}
