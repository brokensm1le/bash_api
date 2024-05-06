package customTime

import "time"

func GetMoscowTime() time.Time {
	return time.Now().Local().Add(time.Hour * time.Duration(3)).UTC()
}
