package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"spider-movie/hleper"
	"strings"
)

func main(){
	request := hleper.NewRequest()
	request.Header.Add("Authorization", "bearereyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOlwvXC9kZXZ0My5vZmZpY2VtYXRlLmNuXC9iYWNrZW5kXC9zaW5nbGVBdXRoIiwiaWF0IjoxNjUwMzQ5Mzg3LCJleHAiOjE3MjIzNDkzODcsIm5iZiI6MTY1MDM0OTM4NywianRpIjoiVGtOZk1rVjlDdDlvTkNGdCIsInN1YiI6MSwicHJ2IjoiODdlMGFmMWVmOWZkMTU4MTJmZGVjOTcxNTNhMTRlMGIwNDc1NDZhYSJ9.pZxxge_mhBRHSBeTsmZ-hrBoxEq0GKuwLnWrvjdsI9M")
	request.Header.Set("Content-Type", hleper.APPLICATION_JOSN)
	request.Url = hleper.Url{
		Host: "http://t3.com:8080",
		Path: "backend/admin/getS4CustomerList?company_code=1000&name=东北大学",
		Query: url.Values{
			"company_code":{"1000"},
			"name":{"东北大学"},
		},
	}
	params := make(map[string]interface{})
	params["company_code"] = "1020"
	params["create_date"] = []string{"2022-04-22 00:00:00", "2022-04-22 23:59:59"}

	result := request.Get()

	fmt.Println(result)
	os.Exit(1)
}

