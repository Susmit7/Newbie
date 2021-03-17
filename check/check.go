package check

import (
	controller "Newbie/controllers"
	model "Newbie/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Checkout(w http.ResponseWriter, r *http.Request) {

	controller.Check("checkout", "POST", w, r)
	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	} else {
		controller.CheckoutHandler(w, id.ID1)

	}
}
