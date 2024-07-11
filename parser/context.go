package parser

type context struct {
	data map[string]any
}

func newcontext() *context {
	return &context{
		data: make(map[string]any),
	}
}
func (c *context) set(key string, val any) {
	c.data[key] = val
}

func (c *context) get(key string) any {
	return c.data[key]
}
