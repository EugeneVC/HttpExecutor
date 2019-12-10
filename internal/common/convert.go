package common

import "strconv"

func ConvertStringToInt(param string) (int, error) {
	if param == "" {
		return 0, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func ConvertStringToInt64(param string) (int64, error) {
	if param == "" {
		return 0, nil
	}

	val, err := strconv.ParseInt(param,10,64)
	if err != nil {
		return 0, err
	}

	return val, nil
}
