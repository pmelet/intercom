package main

func (ctx *shellContext) Set(key string, value bool) {
	ctx.mutex.Lock()
	ctx.keys[key] = value
	ctx.mutex.Unlock()
}
func (ctx *shellContext) Get(key string) bool {
	v, found := ctx.keys[key]
	return found && v
}
