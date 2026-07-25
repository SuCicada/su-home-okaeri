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

func syncWithMultiTry(syncStatus func()) {
	// for i := 0; i < tryCount; i++ {
	// 操作后不要只同步一次，做几次延迟校准
	syncStatus()
	time.Sleep(100 * time.Millisecond)
	syncStatus()
	time.Sleep(500 * time.Millisecond)
	syncStatus()
	// }
}
