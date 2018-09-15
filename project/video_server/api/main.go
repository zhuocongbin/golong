package main 

//一个性能良好的HTTP请求路由器 julienschmidt/httprouter
import(
     
     "net/http" 
     "github.com/julienschmidt/httprouter"

)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	//创建用户
	router.POST("/user",CreateUser)
	//用户登录
	router.POST("/user/:user_name", Login)
	
	return router
	
}
func func main() {
	r := RegisterHandlers()
	http.ListenAndServe(":8000",r)
}