package main


import "testing"

//test if encoding is being done properly from base 10 to base 62, the test checks if 62 is encoded rightly as 10 and 
//61 is encoded rightly as Z

func TestEncode(t *testing.T) {  
    Efirst := Encode(62)
    Esecond := Encode(61)
    if Efirst != "10" {
       t.Errorf("incorrect encode %s", Efirst)
    }
    if Esecond != "Z" {
       t.Errorf("incorrect encode %s", Esecond)
    }
}

//test if decoding from base62 to base10 is being done proerly
func TestDecode(t *testing.T) {  
    Dfirst := Decode("Z")
    Dsecond := Decode("10")
    if Dfirst != 61 {
       t.Errorf("incorrect encode %s", Dfirst)
    }
    if Dsecond != 62 {
       t.Errorf("incorrect encode %s", Dsecond)
    }
}

//test if we are getting the right id of a particular url from the database
func TestDb_process(t *testing.T) {  
    Db_process_first := db_process("https:/google.com")
    
    if Db_process_first != 1 {
       t.Errorf("incorrect encode %s", Db_process_first)
    }
    
}

//test if we are getting the right url of a particular id from the database

func TestGet_url_long(t *testing.T) {  
    Db_Get_url_long := get_url_long(1)
    
    if Db_Get_url_long != "https:/google.com" {
       t.Errorf("incorrect encode %s", Db_Get_url_long)
    }
    
}