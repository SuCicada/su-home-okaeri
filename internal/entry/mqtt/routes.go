package mqttentry

import (
	"sucicada/home/internal/mqttpkg"
)

func RegisterRoutes(r *mqttpkg.Router) {
	RegisterLightRoutes(r)
	RegisterMediaRoutes(r)
}
