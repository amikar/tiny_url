package main

//import packages

import (
	_ "github.com/go-sql-driver/mysql"
    "fmt"
    "database/sql"
    "net/http"
    "strings"
    "log"
    "strconv"

)


//declare constants 
const(
		alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		length = int64(len(alphabet))

		//mysql id and password required for using mysql and this application

		mysql_id = "root"
		mysql_pass = "12"
)

// Encode takes an int token and returns a string base62 value.

func Encode(n int64) string {
	if n == 0 {
		return string(alphabet[0])
	}

	s := ""
	for ; n > 0; n = n / length {
		s = string(alphabet[n%length]) + s
	}
	return s
}

// Decode converts a base62 token to int.

func Decode(key string) (int64) {
	var n int64
	for _, c := range []byte(key) {
		i := strings.IndexByte(alphabet, c)
		if i < 0 {
			return 0
		}
		n = length*n + int64(i)
	}
	return n
}

var text string


//when we enter a site link after http://localhost:8080/s/ the characters after /s/ are taken and checked with the database
//to give an id , this id is then put with http://localhost:8080/g/ to generate a new link

func get_url_and_convert(w http.ResponseWriter, r *http.Request) {

	en := string(r.URL.Path[3:])
	tf := db_process(en)
	gk := Encode(tf) 
    fmt.Fprintf(w, "Hi there, Your new url is : http://localhost:8080/g/%s", gk)
}

//when the link http://localhost:8080/g/"id value" is given the "id value" is given and checked with the database
//the url corresponding to this id value is taken and then redirected

func get_url_and_find(w http.ResponseWriter, r *http.Request) {

	old := string(r.URL.Path[3:])
	//tf := db_process(en)
	de := strconv.FormatInt(Decode(old),10)

	get_url,err :=strconv.Atoi(de)
	if err != nil {
            panic(err.Error())
        }
    get_f_url := get_url_long(get_url)
  	

  	http.Redirect(w, r, get_f_url, 307)

}

//takes an id value as int and checks with the database to find the corresponding actual_url 
//this url_l is returned 

func get_url_long(id int) (string){

	db2 ,err1:= sql.Open("mysql", ""+mysql_id+":"+mysql_pass+"@tcp(127.0.0.1:3306)/")
	if err1 != nil {
    	panic(err1.Error())
	}
	defer db2.Close()

	_, err1 = db2.Exec("use url") 
	if err1 != nil {
		panic(err1.Error()) 
	}

	var actual_url string

	err1 = db2.QueryRow("select url_long from url_s where id = ?",id).Scan(&actual_url)
	if err1 != nil {
		panic(err1.Error()) 
	}

	return actual_url

}

//opens on http://localhost:8080/ as index page


func indexpage(w http.ResponseWriter, r *http.Request){

	fmt.Fprintf(w, "<title>Url Based Url Shortner</title>"+ "<p> Welcome to a url based URL Shortner </p>"+ "<p> Please enter your Url after http://localhost:8080/s/ to shorten it</p>"+
		"<p> To get the actual URL type the shortened url in the url </p>"+
		"<h2> Make sure the url is in the format : https:/yoururl.com </h2>"+
        "</form>")


}


//takes the url_l value and checks if it exists in the database url, if it exists it returns the id of that url_l
//if value does not exist it adds the url_l to the table and gives the id of that url


func db_process(text string) (int64){


	db ,err:= sql.Open("mysql", ""+mysql_id+":"+mysql_pass+"@tcp(127.0.0.1:3306)/")
	if err != nil {
    	panic(err.Error()) 
	}
	defer db.Close()
	

	//database name = url
	_, err = db.Exec("create database if not exists url") 
	if err != nil {
		panic(err.Error()) 
	}
	

	_, err = db.Exec("use url") 
	if err != nil {
		panic(err.Error()) 
	}

	//table name is url_s
	
	_, err = db.Exec("create table if not exists url_s (id int(11) NOT NULL AUTO_INCREMENT, url_long varchar(225), PRIMARY KEY (id))") // ? = placeholder
	if err != nil {
		panic(err.Error()) 
	}

	var fid int64

	
	
	res, err := db.Exec("INSERT INTO url_s (url_long) select '"+text+"' where not exists (select url_long from url_s where url_long = '"+text+"')")
    if err != nil {
        log.Fatal("Cannot prepare DB statement ", err)
    }
    id, _ := res.LastInsertId()


    if id < 1{
    	
    	err := db.QueryRow("select id from url_s where url_long = ?",text).Scan(&fid)
  		if err != nil {
        log.Fatal("Cannot prepare DB statement ", err)}
    	
    	
    }else{
    	fid = id
    }

    return (fid)
}

func main() {

	

	//id, _ := stmtIns.LastInsertId()
	
	//takes the url value from http://localhost:8080/s/ "url"
	http.HandleFunc("/s/", get_url_and_convert)
	
	//redirects to the url from the id http://localhost:8080/g/"id"
	http.HandleFunc("/g/", get_url_and_find)

	//opens to http://localhost:8080/ as index page
	http.HandleFunc("/", indexpage)

	//port in use
    http.ListenAndServe(":8080", nil)



//	fmt.Println("https://linl.com/" + strconv.FormatInt(fid,10))

//	fmt.Println("https://linl.com/" + strconv.FormatInt(de,10))


}
