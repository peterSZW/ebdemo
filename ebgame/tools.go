package ebgame

import (
	"math"
	"os"
	"path/filepath"
	"strings"
)

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err.Error()

	}
	return strings.Replace(dir, "\\", "/", -1)
}
func file_exist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetDegreeByXY(xx, yy float64) float64 {

	// 45  0  315
	// 90  0  270
	// 135 180 225

	if yy == 0 && xx == 0 {
		return 0
	}
	deg := 180 - math.Atan2(yy, xx)/math.Pi*180
	if deg < 270 {
		deg = deg + 90
	} else {
		deg = deg - 270
	}
	return deg
}
func GetDirectByXY(xx, yy float64) int {
	// 7 8 1
	// 6 0 2
	// 5 4 3
	deg := GetDegreeByXY(xx, yy)

	if 22 > deg || deg > 360-22 {
		return 8
	}
	if deg >= 0+22 && deg < 45+22 {
		return 7
	}
	if deg >= 45+22 && deg < 90+22 {
		return 6
	}
	if deg >= 90+22 && deg < 135+22 {
		return 5
	}

	if deg >= 135+22 && deg < 180+22 {
		return 4
	}

	if deg >= 180+22 && deg < 225+22 {
		return 3
	}
	if deg >= 225+22 && deg < 270+22 {
		return 2
	}
	if deg >= 270+22 && deg < 315+22 {
		return 1
	}
	if deg >= 315+22 && deg < 360+22 {
		return 8
	}

	if xx == 0 && yy == 0 {
		return 0
	}
	return 0

	// if xx > 0 && yy < 0 {
	// 	return 1
	// }
	// if xx > 0 && yy == 0 {
	// 	return 2
	// }
	// if xx > 0 && yy > 0 {
	// 	return 3
	// }
	// if xx == 0 && yy > 0 {
	// 	return 4
	// }
	// if xx < 0 && yy > 0 {
	// 	return 5
	// }

	// if xx < 0 && yy == 0 {
	// 	return 6
	// }
	// if xx < 0 && yy < 0 {
	// 	return 7
	// }
	// if xx == 0 && yy < 0 {
	// 	return 8
	// }
	// return 0
}
