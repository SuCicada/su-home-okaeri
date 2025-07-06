package light

type LightDevice interface {
	Get() int
	Set(light int)
	Toggle()
}
