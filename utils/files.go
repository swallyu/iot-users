package utils

import "os"

func WriteFile(filename string,data []byte) (int,error)  {
	fl, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return 0,err
	}
	defer fl.Close()
	n, err := fl.Write(data)
	return n,err
}
