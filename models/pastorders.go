package models

type PastOrder struct {
	Id 				int         `json:"id"`
	Mobile			int64 		`json:"mobile"`
	Time        	int64		`json:"time"`
	NameOne        	string  	`json:"nameone"`
	ImageOne       	string 		`json:"imageone"`
	NameTwo       	string  	`json:"nametwo"`
	ImageTwo       	string 		`json:"imagetwo"`
	NameThree       string  	`json:"nameThree"`
	ImageThree      string 		`json:"imagethree"`
	NameFour        string  	`json:"namefour"`
	ImageFour       string 		`json:"imagefour"`
	Price      		int    		`json:"price"`
	Description 	string 		`json:"description"`
	FoodOneId       int 		`json:"foodoneid"`
	FoodTwoId    	int 		`json:"foodtwoid"`
	FoodThreeId		int 		`json:"foodthreeid"`
	FoodFourId    	int			`json:"foodfourid"`

}

