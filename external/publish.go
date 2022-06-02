package external

import (
	"encoding/json"
	"fmt"
	"log"
)

func PublishViewModel(vm interface{}) error {

	jsonBytes, err := json.Marshal(vm)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("published view model: %s", string(jsonBytes)))

	return nil
}
