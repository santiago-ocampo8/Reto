package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"google.golang.org/grpc"
)

//Structs

type buyer struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type product struct {
	Idproduct   string `json:"idproduct"`
	Nameproduct string `json:"nameproduct"`
	Price       int    `json:"price"`
}
type listProducts struct {
	Idproduct string `json:"idproduct"`
}
type idbuyer struct {
	Id string `json:"id"`
}

type transaction struct {
	Idtransaction string    `json:"idtransaction"`
	Buyerc        buyer     `json:"buyerc"`
	Ip            string    `json:"ip"`
	Device        string    `json:"device"`
	Product       []product `json:"product"`
}

var DB *dgo.Dgraph

//Main
func main() {
	DB = newClient()
	port := ":3000"
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},

		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Get("/", indexRoute)
	r.Post("/buyers/{date}", AddBuyer)
	r.Get("/buyers", getBuyers)
	r.Get("/buyerproducts/{id}", getProductsBuyer)
	r.Get("/buyersip/{ip}", getBuyersIp)
	r.Get("/product/{idproduct}", getProducts)
	http.ListenAndServe(port, r)

}

// Create Database

func newClient() *dgo.Dgraph {
	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		fmt.Println("Err")
	}
	fmt.Println("Db connect")
	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}

//Upload information to the buyers, products and transactions in database
func AddBuyer(w http.ResponseWriter, r *http.Request) {

	date := chi.URLParam(r, "id")
	clienteHttp := &http.Client{}
	request, err := http.NewRequest("GET", "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/buyers?date="+date, nil)
	if err != nil {
		fmt.Println("Error")
	}

	request.Header.Add("Content-Type", "application/json")
	response, err := clienteHttp.Do(request)
	if err != nil {
		fmt.Println("Error")
	}

	defer response.Body.Close()
	bodyResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error")
	}

	content1 := []buyer{}
	json.Unmarshal([]byte(bodyResponse), &content1)

	clienteHttp = &http.Client{}
	request, err = http.NewRequest("GET", "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/products?date="+date, nil)

	if err != nil {
		fmt.Println("Error")
	}

	request.Header.Add("Content-Type", "application/json")

	response, err = clienteHttp.Do(request)

	if err != nil {
		fmt.Println("Error")
	}

	defer response.Body.Close()
	bodyResponse, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error")
	}

	responseString := string(bodyResponse)
	list := strings.Split(responseString, "\n")
	var content2 []product

	for i := 0; i < len(list)-1; i++ {
		temp := strings.Split(list[i], "'")

		var newProduct product
		newProduct.Idproduct = temp[0]
		newProduct.Nameproduct = temp[1]
		newProduct.Price, err = strconv.Atoi(temp[2])
		content2 = append(content2, newProduct)
	}

	clienteHttp = &http.Client{}
	request, err = http.NewRequest("GET", "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/transactions?date="+date, nil)

	if err != nil {
		fmt.Println("Error")
	}

	request.Header.Add("Content-Type", "application/json")

	response, err = clienteHttp.Do(request)

	if err != nil {
		fmt.Println("Error")
	}

	defer response.Body.Close()
	bodyResponse, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error")
	}

	//responseString := string(bodyResponse)
	for i := 0; i < len(bodyResponse); i++ {
		if bodyResponse[i] == 0 {
			bodyResponse[i] = 32
		}

	}
	responseString = string(bodyResponse)
	responseString = strings.TrimRight(responseString, "  ")

	list = strings.Split(responseString, "  #")
	var content3 []transaction
	for i := 0; i < len(list); i++ {
		temp := strings.Split(list[i], " ")

		var newTransaction transaction
		newTransaction.Idtransaction = temp[0]

		for h := 0; h < len(content1); h++ {
			if content1[h].Id == temp[1] {
				newTransaction.Buyerc = content1[h]
				break
			}
		}
		newTransaction.Ip = temp[2]
		newTransaction.Device = temp[3]

		temp[4] = strings.TrimLeft(temp[4], "(")
		temp[4] = strings.TrimRight(temp[4], ")")
		list := strings.Split(temp[4], ",")
		var list2 []product
		for j := 0; j < len(list); j++ {
			var idTemporal = list[j]
			for n := 0; n < len(content2); n++ {
				if content2[n].Idproduct == idTemporal {
					list2 = append(list2, content2[n])
				}
			}

		}

		newTransaction.Product = list2
		content3 = append(content3, newTransaction)
	}

	op := &api.Operation{}
	op.Schema = `
	idtransaction: string @index(term) .
	buyerc: uid @reverse .
	ip: string @index(term) .
	device: string .
	product: [uid] @reverse .
	id: string @index(term) .
	idproduct: string @index(term) .
	
`
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	if err := DB.Alter(ctx, op); err != nil {
		fmt.Println("eBackground")
		log.Fatal(err)

	}

	mu := &api.Mutation{
		CommitNow: true,
	}

	pb, err := json.Marshal(content3)

	if err != nil {
		log.Fatal(err)
		fmt.Println("eMarshal")
	}

	mu.SetJson = pb

	DB.NewTxn().Mutate(ctx, mu)

}

//Get List of Buyers
func getBuyers(w http.ResponseWriter, r *http.Request) {

	const q = `
		{
			find(func:has(id)) @groupby(uid,id,name,age) {
			count(uid)
	  
			}
		}
	`
	resp, err := DB.NewTxn().Query(context.Background(), q)
	if err != nil {
		log.Fatal(err)

	}
	var decode struct {
		Find []struct {
			Groupby []struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Age   int    `json:"age"`
				Count int    `json:"count"`
			} `json:"@groupby"`
		} `json:"find"`
	}

	json.Unmarshal(resp.Json, &decode)
	json.NewEncoder(w).Encode(decode)
}

//Get products buy by a buyer
func getProductsBuyer(w http.ResponseWriter, r *http.Request) {
	k := chi.URLParam(r, "id")

	const q = `
	query find_buyer($a:string){
		find_buyer(func:eq(id,$a)){
		id
	  name
	  age
	  ~buyerc{
			idtransaction
		ip
		device
		product{
		idproduct
		  nameproduct
		  price
		}
	  }
	}
	}
	
	
	`
	resp, err := DB.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$a": k})
	if err != nil {
		log.Fatal(err)

	}
	var decode struct {
		FindBuyer []struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Age    int    `json:"age"`
			Buyerc []struct {
				Idtransaction string `json:"idtransaction"`
				IP            string `json:"ip"`
				Device        string `json:"device"`
				Product       []struct {
					Idproduct   string `json:"idproduct"`
					Nameproduct string `json:"nameproduct"`
					Price       int    `json:"price"`
				} `json:"product"`
			} `json:"~buyerc"`
		} `json:"find_buyer"`
	}

	json.Unmarshal(resp.Json, &decode)
	json.NewEncoder(w).Encode(decode)
}

// Get buyers with the same ip
func getBuyersIp(w http.ResponseWriter, r *http.Request) {
	k := chi.URLParam(r, "ip")

	const q = `
	query find_buyer($a:string){
		find_ip(func: eq(ip,$a)){
			ip
		  buyerc{
			id
			name
			age
		  }
		}
		}
		
	`
	resp, err := DB.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$a": k})
	if err != nil {
		log.Fatal(err)

	}
	var decode struct {
		FindIP []struct {
			IP     string `json:"ip"`
			Buyerc struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Age  int    `json:"age"`
			} `json:"buyerc"`
		} `json:"find_ip"`
	}

	json.Unmarshal(resp.Json, &decode)
	json.NewEncoder(w).Encode(decode)
}

//Product Recommendations
func getProducts(w http.ResponseWriter, r *http.Request) {
	k := chi.URLParam(r, "idproduct")
	const q = `
	query find_buyer($a:string){
		find_products(func:eq(idproduct,$a),first: 2){
		~product {
			idtransaction
		  product{
				idproduct
		  nameproduct
		  price
		}
	}
}
}	
	`
	resp, err := DB.NewTxn().QueryWithVars(context.Background(), q, map[string]string{"$a": k})
	if err != nil {
		log.Fatal(err)

	}
	var decode struct {
		FindProducts []struct {
			Products []struct {
				Idtransaction string `json:"idtransaction"`
				Product       []struct {
					Idproduct   string `json:"idproduct"`
					Nameproduct string `json:"nameproduct"`
					Price       int    `json:"price"`
				} `json:"product"`
			} `json:"~product"`
		} `json:"find_products"`
	}
	json.Unmarshal(resp.Json, &decode)
	json.NewEncoder(w).Encode(decode)
}

//
func indexRoute(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Welcome to my API")

}
