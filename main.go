package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jsvm/twitter-transaction-payload-gen/payload"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type input struct {
	Key             string   `json:"key"`
	Frames          []string `json:"frames"`
	FramesConverted [][][]int
}

var (
	regex = regexp.MustCompile(`[^\d]+`)
)

func newAPI(c *gin.Context) {
	//fmt.Println("HIII")
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	//fmt.Println(string(b))
	f := &input{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		c.AbortWithStatus(501)
		return
	}
	formatted := [][][]int{}
	for _, str := range f.Frames {
		str = str[9:]
		parts := strings.Split(str, "C")
		tism := [][]int{}
		for i, p := range parts {
			things := regex.FindAllString(p, -1)
			for _, thing := range things {
				parts[i] = strings.Replace(parts[i], thing, " ", 1)
			}
			p = parts[i]
			if string(p[len(p)-1]) == " " {
				parts[i] = parts[i][:len(p)-1]
			}
			if string(p[0]) == " " {
				parts[i] = parts[i][1:]
			}
			nums := []int{}
			for _, a := range strings.Split(parts[i], " ") {
				num, _ := strconv.Atoi(a)
				nums = append(nums, num)
			}
			tism = append(tism, nums)
		}
		formatted = append(formatted, tism)
	}
	locker.Lock()
	if _, ok := keys[f.Key]; ok {
		locker.Unlock()
		c.AbortWithStatus(502)
		return
	}
	f.FramesConverted = formatted
	keys[f.Key] = f
	locker.Unlock()
	c.Status(200)
}

var (
	keys   = make(map[string]*input)
	locker sync.RWMutex
)

type getInput struct {
	Key    string `json:"key"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

type output struct {
	Header string `json:"header"`
}

func getAPI(c *gin.Context) {
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	f := &getInput{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		c.AbortWithStatus(501)
		return
	}
	locker.Lock()
	if _, ok := keys[f.Key]; !ok {
		locker.Unlock()
		c.AbortWithStatus(502)
		return
	}
	obj := keys[f.Key]
	locker.Unlock()

	c.JSON(200, &output{Header: payload.GenerateHeader(f.Path, f.Method, f.Key, obj.FramesConverted)})
}

func main() {
	router := gin.New()
	router.Use(gin.Recovery())
	router.POST("new", newAPI)
	router.POST("get", getAPI)
	s := &http.Server{
		// ! change after ready to go out
		Addr:           ":6969",
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("httpServe: %s", err)
	}
}
