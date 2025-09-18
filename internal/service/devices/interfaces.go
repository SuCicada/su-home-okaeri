package devices

import (
	"SuCicada/home/internal/cfg"
	"SuCicada/home/internal/structs/appconfig"
	"reflect"
	"strings"
)

type DeviceBase struct {
	Name          string            `yaml:"name"`
	DeviceControl DeviceControlUnit `yaml:"control"`
}

type DeviceControlUnit struct {
	Light  *Control `yaml:"light"`
	Volume *Control `yaml:"volume"`
}

// func (d *DeviceControlUnit) GetControl(name string) *Control {
// 	control := map[string]*Control{
// 		"light":  d.Light,
// 		"volume": d.Volume,
// 	}[name]
// 	return control
// }

func (d *DeviceControlUnit) GetControl(name string) *Control {
	v := reflect.ValueOf(d).Elem() // 获取结构体的值
	t := reflect.TypeOf(d).Elem()  // 获取结构体的类型

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		yamlTag := field.Tag.Get("yaml")

		// 如果yaml tag匹配，或者字段名匹配（忽略大小写）
		if yamlTag == name || strings.EqualFold(field.Name, name) {
			fieldValue := v.Field(i)
			if fieldValue.Kind() == reflect.Ptr && fieldValue.Type().Elem() == reflect.TypeOf(Control{}) {
				if !fieldValue.IsNil() {
					return fieldValue.Interface().(*Control)
				}
			}
		}
	}
	return nil
}

type Control struct {
	Name    string
	Device  *DeviceBase
	Control ControllerInterface
}

func (d *Control) Toggle() error {
	value := d.GetValue()
	v, err := d.Control.Get()
	if err != nil {
		return err
	}
	var newValue int
	if v >= value.High {
		newValue = value.Low
	} else {
		newValue = value.High
	}
	return d.Control.Set(newValue)
}
func (d *Control) GetValue() appconfig.Value {
	device := cfg.GetConfig().Devices[d.Device.Name]
	control := device.Control[d.Name]
	return control
}

func (d *Control) HIGH() int {
	return d.GetValue().High
}
func (d *Control) LOW() int {
	return d.GetValue().Low
}

type ControllerInterface interface {
	Get() (int, error)
	Set(value int) error
}
