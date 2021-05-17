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
type allBuyer []buyer
type allProducts []product
type allTransactions []transaction
type alllistProducts []listProducts

var buyers = allBuyer{}
var transsactions = allTransactions{}
var listsProducts = alllistProducts{}
var products = allProducts{}
var DB *dgo.Dgraph

func main() {
	DB = newClient()
	port := ":3000"
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Get("/", indexRoute)
	r.Post("/buyers/{date}", AddBuyer)
	r.Get("/buyers", getBuyers)
	r.Get("/buyerproducts/{id}", getProductsBuyer)
	r.Get("/buyersip/{ip}", getBuyersIp)
	r.Get("/product/{idproduct}", getProducts)
	r.Post("/borrar", borrar)
	http.ListenAndServe(port, r)

}

func AddBuyer(w http.ResponseWriter, r *http.Request) {

	j := chi.URLParam(r, "id")

	clienteHttp := &http.Client{}

	peticion, err := http.NewRequest("GET", "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/buyers?date="+j, nil)

	if err != nil {
		fmt.Println("Error")
	}

	peticion.Header.Add("Content-Type", "application/json")

	respuesta, err := clienteHttp.Do(peticion)

	if err != nil {
		fmt.Println("Error")
	}

	defer respuesta.Body.Close()
	cuerpoRespuesta, err := ioutil.ReadAll(respuesta.Body)
	if err != nil {
		fmt.Println("Error")
	}

	content1 := []buyer{}
	json.Unmarshal([]byte(cuerpoRespuesta), &content1)

	clienteHttp = &http.Client{}
	peticion, err = http.NewRequest("GET", "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/products?date="+j, nil)

	if err != nil {
		fmt.Println("Error")
	}

	peticion.Header.Add("Content-Type", "application/json")

	respuesta, err = clienteHttp.Do(peticion)

	if err != nil {
		fmt.Println("Error")
	}

	defer respuesta.Body.Close()
	cuerpoRespuesta, err = ioutil.ReadAll(respuesta.Body)
	if err != nil {
		fmt.Println("Error")
	}

	respuestaString := string(cuerpoRespuesta)
	arreglo := strings.Split(respuestaString, "\n")
	var content2 []product
	for i := 0; i < len(arreglo)-1; i++ {
		temp := strings.Split(arreglo[i], "'")

		var newProduct product
		newProduct.Idproduct = temp[0]
		newProduct.Nameproduct = temp[1]
		newProduct.Price, err = strconv.Atoi(temp[2])
		content2 = append(content2, newProduct)
	}

	clienteHttp = &http.Client{}
	peticion, err = http.NewRequest("GET", "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/transactions?date="+j, nil)

	if err != nil {
		fmt.Println("Error")
	}

	peticion.Header.Add("Content-Type", "application/json")

	respuesta, err = clienteHttp.Do(peticion)

	if err != nil {
		fmt.Println("Error")
	}

	defer respuesta.Body.Close()
	cuerpoRespuesta, err = ioutil.ReadAll(respuesta.Body)
	if err != nil {
		fmt.Println("Error")
	}

	//respuestaString := string(cuerpoRespuesta)
	for i := 0; i < len(cuerpoRespuesta); i++ {
		if cuerpoRespuesta[i] == 0 {
			cuerpoRespuesta[i] = 32
		}

	}
	respuestaString = string(cuerpoRespuesta)
	respuestaString = strings.TrimRight(respuestaString, "  ")

	arreglo = strings.Split(respuestaString, "  #")
	var content3 []transaction
	for i := 0; i < len(arreglo); i++ {
		temp := strings.Split(arreglo[i], " ")

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

			//list2 = append(list2, list[j])
			var idTemporal = list[j]
			for n := 0; n < len(content2); n++ {
				//fmt.Println(n, "<- Aca va n, aca va el tamañp->", len(content2), " Y aca va j->", j)
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

/*func getBuyers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buyers)
}*/

//Products

//Database

func newClient() *dgo.Dgraph {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		fmt.Println("Esto es un error")
	}
	fmt.Println("Esto esta bien")
	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}

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
	//fmt.Println(" no se que putas es esto ", resp.Json)
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
func getProducts(w http.ResponseWriter, r *http.Request) {
	k := chi.URLParam(r, "idproduct")
	fmt.Println(k)
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
	fmt.Println(string(resp.Json))
	fmt.Println(decode)
	json.Unmarshal(resp.Json, &decode)
	json.NewEncoder(w).Encode(decode)
}
func setup() {
	// Install a schema into dgraph. Accounts have a `name` and a `balance`.
	err := DB.Alter(context.Background(), &api.Operation{
		Schema: `
			id: string @index(term) .
			name: string .
			age: int .
		`,
	})
	fmt.Println(" no se que putas es esto ", err)
	const q = `
		{
			all(func: has(name)) {
				uid
				name
			}
		}
	`
	resp, err := DB.NewTxn().Query(context.Background(), q)
	if err != nil {
		log.Fatal(err)
		fmt.Println("esta consulta no se acaba de realizar", resp)
	}
	fmt.Println("esta consulta se acaba de realizar", string(resp.Json))

}

func indexRoute(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Welcome to my API")

}
func borrar(w http.ResponseWriter, r *http.Request) {

	cont := context.Background()
	op := &api.Operation{}
	op.DropAll = true
	err := DB.Alter(cont, op)
	if err != nil {
		log.Fatal(err)
	}
}
