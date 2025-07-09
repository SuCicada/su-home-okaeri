package light

type LightDevice interface {
	Get() (int, error)
	Set(light int) (string, error)
	Toggle() (string, error)
}
