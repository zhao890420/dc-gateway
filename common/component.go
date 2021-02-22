package common

type Destroyable interface {
	Destroy()
}

var Destroyables []Destroyable

func RegisterComponent(c Destroyable) {
	Destroyables = append(Destroyables, c)
}
