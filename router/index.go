package router

//func IndexRouterHandler(r *gin.Engine) {
//	group := r.Group("/")
//	group.Use(middleware.Autenticate())
//	group.GET("/index", func(context *gin.Context) {
//		context.HTML(http.StatusOK, "index.tmpl", gin.H{})
//	})
//	group.GET("/", func(context *gin.Context) {
//		context.Redirect(http.StatusFound, "/index")
//	})
//}
