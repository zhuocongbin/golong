package main 

//一个性能良好的HTTP请求路由器 julienschmidt/httprouter
import(
     
     "net/http" 
     "github.com/julienschmidt/httprouter"

)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.POST("/user",CreateUser)
	return router
	
}
func func main() {
	r := RegisterHandlers()
	http.ListenAndServe(":8000",r)
}