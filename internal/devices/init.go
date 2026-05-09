package devices

import (
	_ "sucicada/home/internal/devices/linux"
	_ "sucicada/home/internal/devices/redmi"
	_ "sucicada/home/internal/devices/windows"
)

func Init() {
}
