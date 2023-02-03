package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Req struct {
	Fruit    string `json:"fruit"`
	Quantity int    `json:"quantity"`
}

type Res struct {
	Fruits map[string]int `json:"fruits"`
}

var rdb *redis.Client

func get(fruitName string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	total, err := rdb.IncrBy(ctx, fruitName, 0).Result()
	return int(total), err
}

func incrBy(fruitName string, quantity int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	total, err := rdb.IncrBy(ctx, fruitName, int64(quantity)).Result()
	return int(total), err
}

func decrBy(fruitname string, quantity int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	total, err := rdb.DecrBy(ctx, fruitname, int64(quantity)).Result()
	return int(total), err
}

func list() (map[string]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	fruits, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}
	m := make(map[string]int)
	for _, fruitName := range fruits {
		total, err := get(fruitName)
		if err != nil {
			return nil, err
		}
		m[fruitName] = total
	}
	return m, nil
}

func main() {

	redis_url := os.Getenv("REDIS_URL")
	if len(redis_url) == 0 {
		redis_url = "localhost:6379"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     redis_url,
		Password: "",
		DB:       0,
	})

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				respondWithError(w, "Invalid API request at: "+r.URL.Path+"\n")
				return
			}
			respond(w)
		})
	http.HandleFunc("/buy", buy)
	http.HandleFunc("/sell", sell)
	port := os.Getenv("API_PORT")
	if len(port) == 0 {
		port = "8081"
	}
	fmt.Println("listening on port: ", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	(*w).Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
}

func getReq(w http.ResponseWriter, r *http.Request) (*Req, error) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		return nil, nil
	}
	fmt.Printf("request: %+v\n", r.URL)
	var req Req
	fruit := r.URL.Query().Get("fruit")
	quantityStr := r.URL.Query().Get("quantity")
	fromParams := true
	if len(fruit) == 0 || len(quantityStr) == 0 {
		fromParams = false
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return nil, fmt.Errorf("either fruit or quantity is not provided but required\n")
		}
	}
	if fromParams {
		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			return nil, fmt.Errorf("quantity must be a number\n")
		}
		req.Fruit = fruit
		req.Quantity = quantity
	}
	if len(req.Fruit) == 0 || req.Quantity == 0 {
		return nil, fmt.Errorf("either fruit or quantity is not provided but required\n")
	}
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be a positive number\n")
	}
	return &req, nil
}

func buy(w http.ResponseWriter, r *http.Request) {
	req, err := getReq(w, r)
	if err != nil {
		respondWithError(w, err.Error())
		return
	}
	if req == nil {
		respond(w)
		return
	}
	fruit := req.Fruit
	quantity := req.Quantity

	_, err = incrBy(fruit, quantity)
	if err != nil {
		respondWithError(w, err.Error())
		return
	}

	respond(w)
}

func sell(w http.ResponseWriter, r *http.Request) {
	req, err := getReq(w, r)
	if err != nil {
		respondWithError(w, err.Error())
		return
	}
	if req == nil {
		respond(w)
		return
	}
	fruit := req.Fruit
	quantity := req.Quantity

	current, err := get(fruit)
	if err != nil {
		respondWithError(w, err.Error())
		return
	}

	if current < quantity {
		respondWithError(w, "not enough fruits\n")
		return
	}

	_, err = decrBy(fruit, quantity)
	if err != nil {
		respondWithError(w, err.Error())
		return
	}
	respond(w)
}

func respond(w http.ResponseWriter) {
	enableCors(&w)
	fruits, err := list()
	if err != nil {
		respondWithError(w, err.Error())
		return
	}
	b, _ := json.Marshal(Res{fruits})
	fmt.Fprint(w, string(b))
}

func respondWithError(w http.ResponseWriter, message string) {
	fmt.Print("error: ", message)
	enableCors(&w)
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(message))
}
