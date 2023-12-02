package SendError

import ()

// this is for mutiple values check the value is exists
func RequestDataError(data map[string]string, data_name []string ) string {
	for i := 0; i < len(data_name); i++ {
        if data[data_name[i]] == ""{
			return data_name[i]
		}
    }
	return ""
}