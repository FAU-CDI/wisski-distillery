package cancel

import "context"

// ValuesOf returns a new context that has the same deadline and cancelation behviour as parent.
// However when requesting values from the context, checks the values in context first.
func ValuesOf(parent, values context.Context) context.Context {
	return &valuesOf{
		Context: parent,
		values:  values,
	}
}

type valuesOf struct {
	context.Context
	values context.Context
}

func (vv *valuesOf) Value(key any) any {
	if value := vv.values.Value(key); value != nil {
		return value
	}
	return vv.Context.Value(key)
}
