package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func sendResponse(enc *json.Encoder, data interface{}, code int, w http.ResponseWriter) {
	w.WriteHeader(code)

	if err := enc.Encode(data); err != nil {
		fmt.Println("error ")
		fmt.Println(err)
	}

}
