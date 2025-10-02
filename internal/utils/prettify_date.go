package utils

func PrettifyDate(dateStr string) string {
	if len(dateStr) != 13 {
		return dateStr
	}
	return dateStr[:4] + "-" + dateStr[4:6] + "-" + dateStr[6:8] + " " + dateStr[9:11] + ":" + dateStr[11:13]
}
