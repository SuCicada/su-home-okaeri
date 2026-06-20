package mqttentry

import "time"

func startStatusTicker(syncStatus func(), interval time.Duration) {
	go func() {
		syncStatus()
		ticker := time.NewTicker(interval)
		for range ticker.C {
			syncStatus()
		}
	}()
}
