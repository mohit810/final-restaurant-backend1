package goroutine

import (
	"bytes"
	"encoding/json"
	"net/http"
	"restaurant_backend/models"
	"sync"
)

func FcmRequest(Wg *sync.WaitGroup,UserKey,ServerKey,Title ,Content *string){
	/*var u  models.Fcmpayload
	u.Data = models.Fcmdata{Title: *Title,Content: *Content}
	u.To = *UserKey*/
	var z models.Fcmdata
	z.Title = *Title
	z.Content = *Content
	c,err:= json.Marshal(z)
	if err != nil {
		panic(err)
	}
	url := "https://fcm.googleapis.com/fcm/send"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(c))
	if err != nil {
		panic(err)
	}
	req.Header.Set("authorization", *ServerKey)
	req.Header.Set("Content-Type", "application/json")
	Wg.Done()

}
