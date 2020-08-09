package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"final-restaurant-backend1/goroutine"
	"final-restaurant-backend1/models"
	"strconv"
)

type UserController struct {
session *sql.DB
}
//var wg sync.WaitGroup {will need to call free(wg) for destroying this global attribute}

func NewUserController(s *sql.DB) *UserController {
return &UserController{s}
}

func (uc UserController) Starter(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf( w,`Your at homepage and Connection to db was success and check the terminal for dump diagnosis!`)
	fmt.Println("Request Dump:\n", string(dump))
}

func (uc UserController) MobileCreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.ParseMultipartForm(10 << 20)
	UserName := r.FormValue("name")
	EmailId := r.FormValue("email")
	PassWord := r.FormValue("password")
	MobileNumber := r.FormValue("mobile")
	i, _ := strconv.ParseInt(MobileNumber, 0, 64)
	AddRess := r.FormValue("address")
	ProfilePic, _ , err := r.FormFile("profilePic")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer ProfilePic.Close()
	_, err = os.Stat("test")

	if os.IsNotExist(err) {
		errDir := os.MkdirAll("profile-pics"+"/"+MobileNumber, 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}

	file, err := os.Create("profile-pics"+"/"+MobileNumber+"/"+"dp.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file,ProfilePic)
	if err != nil {
		log.Fatal(err)
	}


	query, err := uc.session.Prepare("Insert loggedin SET name=?, email=?,pwd=?,number=?,address=?,profilepic=?")
	if err != nil {
		panic(err)
	}
	_, err = query.Exec(UserName, EmailId,PassWord,i,AddRess,"http://3.19.133.66:8080/static/"+file.Name())
	if err != nil {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden) // 403

	}
	defer query.Close()
	var u []*models.User// declare a slice of courses that will hold all of the Course instances scanned from the rows object

	tent,err := uc.session.Query("SELECT name, email, pwd, number, address, profilepic from loggedin where email=?",EmailId)
	if err != nil{
		fmt.Println(err)
		return
	}
	defer tent.Close()
	for tent.Next(){
		c:=new(models.User)// initialize a new instance
		err = tent.Scan(&c.Name,&c.Email,&c.Password,&c.Mobile,&c.Address,&c.ProfilePic)
		if err != nil{
			fmt.Println(err)
			return
		}
		u = append(u,c) // add each instance to the slice
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		fmt.Println(err)
	}
	/*data,err := ioutil.ReadFile("/dev/deployment/go/src/restaurant_backend/profile-pics/dp.png")
	if err != nil { fmt.Fprint(w, err) }
	http.ServeContent(w,r,"dp.png",time.Now(), bytes.NewReader(data))*/
}

func (uc UserController) DishCreate(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u := models.Dishes{}

	json.NewDecoder(r.Body).Decode(&u)

	query, err := uc.session.Prepare("Insert dishes SET name=?, image=?,category=?,price=?,description=?,label=?,featured=?")
	if err != nil {
		panic(err)
	}
	_, err = query.Exec(u.Name,u.Image,u.Category,u.Price,u.Description,u.Label,u.Rating,u.Time,u.Featured)
	if err != nil {
		panic(err)
	}
	defer query.Close()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		fmt.Println(err)
	}
}
//reactjs get dishes
func (uc UserController) GetDishes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	mainId := r.URL.Query().Get("id")
	sideId := r.URL.Query().Get("idd")
	var u []*models.Dishes// declare a slice of courses that will hold all of the Course instances scanned from the rows object
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	switch mainId == "" {
	case true:
		tent,err := uc.session.Query("SELECT * from dishes")
		if err != nil{
			fmt.Println(err)
			return
		}
		defer tent.Close()
		for tent.Next(){
			c:=new(models.Dishes)// initialize a new instance
			err = tent.Scan(&c.Id,&c.Name,&c.Image,&c.Category,&c.Price,&c.Description,&c.Label,&c.Rating,&c.Time,&c.Featured)
			if err != nil{
				fmt.Println(err)
				return
			}
			u = append(u,c) // add each instance to the slice
		}
		getjson(w,u)
	case false:
		switch sideId == "" {
		case true:
			tent,err := uc.session.Prepare("SELECT * from dishes where category=?")
			if err != nil{
				fmt.Println(err)
				return
			}
			defer tent.Close()
			row,err := tent.Query(mainId)
			for row.Next(){
				c:=new(models.Dishes)// initialize a new instance
				err = row.Scan(&c.Id,&c.Name,&c.Image,&c.Category,&c.Price,&c.Description,&c.Label,&c.Rating,&c.Time,&c.Featured)
				if err != nil{
					fmt.Println(err)
					return
				}
				u = append(u,c) // add each instance to the slice
			}
			getjson(w,u)
		case false:
			tent,err := uc.session.Prepare("SELECT * from dishes where category=?")
			if err != nil{
				fmt.Println(err)
				return
			}
			defer tent.Close()
			row,err := tent.Query(mainId)
			for row.Next(){
				c:=new(models.Dishes)// initialize a new instance
				err = row.Scan(&c.Id,&c.Name,&c.Image,&c.Category,&c.Price,&c.Description,&c.Label,&c.Rating,&c.Time,&c.Featured)
				if err != nil{
					fmt.Println(err)
					return
				}
				u = append(u,c) // add each instance to the slice
			}
			rows,err := tent.Query(sideId)
			for rows.Next(){
				c:=new(models.Dishes)// initialize a new instance
				err = rows.Scan(&c.Id,&c.Name,&c.Image,&c.Category,&c.Price,&c.Description,&c.Label,&c.Rating,&c.Time,&c.Featured)
				if err != nil{
					fmt.Println(err)
					return
				}
				u = append(u,c) // add each instance to the slice
			}
			getjson(w,u)
		}
	}

}

func (uc UserController) PlaceOrder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idOne := r.URL.Query().Get("starter")
	idTwo := r.URL.Query().Get("main")
	idThree := r.URL.Query().Get("bread")
	idFour := r.URL.Query().Get("sweet")
	idFive := r.URL.Query().Get("mob")
	QuantOne := r.URL.Query().Get("qone")
	QuantTwo := r.URL.Query().Get("qtwo")
	QuantThree := r.URL.Query().Get("qthree")
	QuantFour := r.URL.Query().Get("qfour")
	var u []*models.Dishes
	var x models.PastOrder
	x.Mobile, _ = strconv.ParseInt(idFive, 0, 64)
	var price int
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	smnt,err := uc.session.Prepare("SELECT * from dishes where id = ?")
	if err != nil{
		fmt.Println(err)
		return
	}
	defer smnt.Close()
	rows,err := smnt.Query(idOne)
	for rows.Next(){
		c:=new(models.Dishes)// initialize a new instance
		err = rows.Scan(&c.Id,&c.Name,&c.Image,&c.Category,&c.Price,&c.Description,&c.Label,&c.Rating,&c.Time,&c.Featured)
		if err != nil{
			fmt.Println(err)
			return
		}
		price = c.Price
		x.FoodOneId = c.Id
		x.ImageOne = c.Image
		x.NameOne = c.Name
		x.Time = c.Time
		x.Description = c.Name + " x " + QuantOne
		u = append(u,c) // add each instance to the slice
	}
	row,err := smnt.Query(idTwo)
	for row.Next(){
		c:=new(models.Dishes)// initialize a new instance
		err = row.Scan(&c.Id,&c.Name,&c.Image,&c.Category,&c.Price,&c.Description,&c.Label,&c.Rating,&c.Time,&c.Featured)
		if err != nil{
			fmt.Println(err)
			return
		}
		price = price + c.Price
		x.FoodTwoId = c.Id
		x.ImageTwo = c.Image
		x.NameTwo = c.Name
		x.Time = x.Time + c.Time
		x.Description = x.Description + ", " + c.Name + " x " + QuantTwo
		u = append(u,c) // add each instance to the slice
	}
	roww,err := smnt.Query(idThree)
	for roww.Next(){
		c:=new(models.Dishes)// initialize a new instance
		err = roww.Scan(&c.Id,&c.Name,&c.Image,&c.Category,&c.Price,&c.Description,&c.Label,&c.Rating,&c.Time,&c.Featured)
		if err != nil{
			fmt.Println(err)
			return
		}
		price = price + c.Price
		x.FoodThreeId = c.Id
		x.ImageThree = c.Image
		x.NameThree = c.Name
		x.Time = x.Time + c.Time
		x.Description = x.Description + ", " + c.Name + " x " + QuantThree
		u = append(u,c) // add each instance to the slice
	}
	roow,err := smnt.Query(idFour)
	for roow.Next(){
		c:=new(models.Dishes)// initialize a new instance
		err = roow.Scan(&c.Id,&c.Name,&c.Image,&c.Category,&c.Price,&c.Description,&c.Label,&c.Rating,&c.Time,&c.Featured)
		if err != nil{
			fmt.Println(err)
			return
		}
		price = price + c.Price
		x.FoodFourId = c.Id
		x.ImageFour = c.Image
		x.NameFour = c.Name
		x.Time = x.Time + c.Time
		x.Description = x.Description + ", " + c.Name + " x " + QuantFour
		u = append(u,c) // add each instance to the slice
	}
	x.Price = price
	query, err := uc.session.Prepare("Insert lastorder SET number=?, time=?,foodnameone=?,foodnametwo=?,foodnamethree=?,foodnamefour=?,price=?,description=?,foodoneid=?,foodtwoid=?,foodthreeid=?,foodfourid=?,foodoneimage=?,foodtwoimage=?,foodthreeimage=?,foodfourimage=?")
	if err != nil {
		panic(err)
	}
	_, err = query.Exec(x.Mobile,x.Time,x.NameOne,x.NameTwo,x.NameThree,x.NameFour,x.Price,x.Description,x.FoodOneId,x.FoodTwoId,x.FoodThreeId,x.FoodFourId,x.ImageOne,x.ImageTwo,x.ImageThree,x.ImageFour)
	if err != nil {
		panic(err)
	}
	defer query.Close()
	getjson(w,u)
}

func (uc UserController) MobileLogin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	EmailId := r.FormValue("email")
	PassWord := r.FormValue("password")
	var u []*models.User
	tent,err := uc.session.Query("SELECT name, email, pwd, number, address, profilepic from loggedin where email=? AND pwd=?",EmailId,PassWord)
	if err != nil{
		fmt.Println(err)
		return
	}
	defer tent.Close()
	for tent.Next(){
		c:=new(models.User)// initialize a new instance
		err = tent.Scan(&c.Name,&c.Email,&c.Password,&c.Mobile,&c.Address,&c.ProfilePic)
		if err != nil{
			fmt.Println(err)
			return
		}
		if c == nil{
			w.WriteHeader(http.StatusNotFound) //404
		}
		u = append(u,c) // add each instance to the slice
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 201
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		fmt.Println(err)
	}
}

func (uc UserController) GetImage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Request Dump:\n", string(dump))
	w.WriteHeader(http.StatusOK) // 201
	_ = r.URL.Query().Get("imagi")
	http.ServeFile(w, r, "")
}

func (uc UserController) CreateComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u := models.Comments{}

	json.NewDecoder(r.Body).Decode(&u)
	query, err := uc.session.Prepare("Insert comments SET dishid=?, rating=?,authorid=?,comment=?,date=?,id=?")
	if err != nil {
		panic(err)
	}
	_, err = query.Exec(u.DishID,u.Rating,u.Author,u.Comment,u.Date,u.ID)
	if err != nil {
		panic(err)
	}
	defer query.Close()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		fmt.Println(err)
	}
}
//reactjs get comments
func (uc UserController) GetComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var u []*models.Comments
	tent,err := uc.session.Query("SELECT dishid, rating,authorid,comment,date,id from comments")
	if err != nil{
		fmt.Println(err)
		return
	}
	defer tent.Close()
	for tent.Next(){
		c:=new(models.Comments)// initialize a new instance
		err = tent.Scan(&c.DishID,&c.Rating,&c.Author,&c.Comment,&c.Date,&c.ID)
		if err != nil{
			fmt.Println(err)
			return
		}
		u = append(u,c) // add each instance to the slice
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		fmt.Println(err)
	}
}
//reactjs promotions
func (uc UserController) Promo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var u []*models.Promotions
	tent,err := uc.session.Query("SELECT * from promotions")
	if err != nil{
		fmt.Println(err)
		return
	}
	defer tent.Close()
	for tent.Next(){
		c:=new(models.Promotions)// initialize a new instance
		err = tent.Scan(&c.ID,&c.Name,&c.Image,&c.Label,&c.Price,&c.Featured,&c.Description)
		if err != nil{
			fmt.Println(err)
			return
		}
		u = append(u,c) // add each instance to the slice
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		fmt.Println(err)
	}
}
//reactjs leaders
func (uc UserController) Leader(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var u []*models.Leader

	tent,err := uc.session.Query("SELECT * from leaders")
	if err != nil{
		fmt.Println(err)
		return
	}
	defer tent.Close()
	for tent.Next(){
		c:=new(models.Leader)// initialize a new instance
		err = tent.Scan(&c.Id,&c.Name,&c.Image,&c.Designation,&c.Abbr,&c.Featured,&c.Description)
		if err != nil{
			fmt.Println(err)
			return
		}
		u = append(u,c) // add each instance to the slice
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w,"hi")
}

func (uc UserController) CreateLeader(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u := models.Leader{}
	json.NewDecoder(r.Body).Decode(&u)
	tent,err := uc.session.Prepare("insert leaders SET id=?, name=?, image=?, designation=?, abbr=?, featured=?, description=?")
	if err != nil {
		fmt.Println(err)
	}
	defer tent.Close()
	_, err = tent.Exec(u.Id,u.Name,u.Image,u.Designation,u.Abbr,u.Featured,u.Description)
	if err != nil {
		panic(err)
	}
	defer tent.Close()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		fmt.Println(err)
	}

}

func (uc UserController) FeedBack(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u := models.Feedback{}
	json.NewDecoder(r.Body).Decode(&u)
	tent,err := uc.session.Prepare("insert feedback SET firstname=?, lastname=?, telnum=?, email=?, contacttype=?, message=?,  date=?,id=?,agree=?")
	if err != nil {
		fmt.Println(err)
	}
	defer tent.Close()
	_, err = tent.Exec(u.Firstname,u.Lastname,u.Telnum,u.Email,u.ContactType,u.Message,u.Date,u.ID,u.Agree)
	if err != nil {
		panic(err)
	}
	defer tent.Close()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		fmt.Println(err)
	}
}

func (uc UserController) NotifyUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	UserKey := r.FormValue("user")
	ServerKey := r.FormValue("server")
	Title:= r.FormValue("title")
	Content := r.FormValue("content")
	var x models.Concurrencyy
	x.Mutexx.Lock()// locks so that no one run this go routine
	x.Wg.Add(1)
	go goroutine.FcmRequest(&x.Wg,&UserKey,&ServerKey,&Title,&Content)
	x.Wg.Wait()
	x.Mutexx.Unlock()
}

func (uc UserController) PastOrders(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func (uc UserController) LastOrder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idOne := r.URL.Query().Get("mob")
	y , _:= strconv.ParseInt(idOne, 0, 64)
	var x []*models.PastOrder// initialize a new instance

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	smnt,err := uc.session.Query("SELECT * from lastorder where number=?",y)
	if err != nil{
		fmt.Println(err)
		return
	}
	defer smnt.Close()
	for smnt.Next(){
		c:=new(models.PastOrder)
		err = smnt.Scan(&c.Mobile,&c.Time,&c.NameOne,&c.NameTwo,&c.NameThree,&c.NameFour,&c.Description,&c.Price,&c.FoodOneId,&c.FoodTwoId,&c.FoodThreeId,&c.FoodFourId,&c.ImageOne,&c.ImageTwo,&c.ImageThree,&c.ImageFour,&c.Id)
		if err != nil{
			fmt.Println(err)
			return
		}
	x = append(x,c)
	}
	err = json.NewEncoder(w).Encode(x)
	if err != nil{
		fmt.Println(err)
		return
	}
}

func getjson(w http.ResponseWriter,u []*models.Dishes)  {
	err := json.NewEncoder(w).Encode(u)
	if err != nil {
		fmt.Println(err)
	}
}
/*curl -X POST -H "Content-Type: application/json" -d '{"name":"Grilled Chicken","image":"https://imgur.com/a/dbkjXoN","category":"starters","price":150,"description":"Grilled Chicken with Cherry salsa ","label":"Grilled","comments":"this is a test comment"}' http://40.114.117.222:8080/user*/
