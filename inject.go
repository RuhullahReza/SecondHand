package main

import (
	"github.com/RuhullahReza/SecondHand/controller"
	"github.com/RuhullahReza/SecondHand/helper"
	"github.com/RuhullahReza/SecondHand/middleware"
	"github.com/RuhullahReza/SecondHand/repository"
	"github.com/RuhullahReza/SecondHand/service"
	"github.com/cloudinary/cloudinary-go"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Inject(db *sqlx.DB, cld *cloudinary.Cloudinary) *gin.Engine {

	accountRepository := repository.NewAccountRepository(db)
	profileRepository := repository.NewProfileRepository(db)
	productRepository := repository.NewProductRepository(db)
	dataRepository := repository.NewDataRepository(db)
	imageRepository := repository.NewImageRepository(cld,db)
	transactionRepostory := repository.NewTransactionRepository(db)

	userService := service.NewUserService(accountRepository, profileRepository, imageRepository, dataRepository)
	dataService := service.NewDataService(dataRepository)
	productService := service.NewProductService(productRepository, profileRepository, dataRepository, imageRepository)
	transactionService := service.NewTransactionSerive(productRepository, profileRepository, transactionRepostory)

	translator := helper.InitTranslator()

	userController := controller.NewUserController(userService, translator)
	dataController := controller.NewDataController(dataService, translator)
	productController := controller.NewProductController(productService, translator)
	transactionController := controller.NewTransactionController(transactionService, translator)

	router := gin.Default()

	router.POST("/register", userController.Register)
	router.POST("admin/register", userController.RegisterAdmin)
	router.POST("/login", userController.Login)
	router.GET("/profile", middleware.Auth(), userController.MyProfile)
	router.GET("/profile/:id",userController.Profile)
	router.PUT("/profile", middleware.Auth(), userController.Update)
	router.PUT("/profile_image", middleware.Auth(), userController.UpdateImage)

	router.GET("/data/city", dataController.GetAllCity)
	router.POST("/data/city", middleware.Auth(), middleware.IsAdmin(), dataController.AddCity)
	router.DELETE("/data/city/:id", middleware.Auth(), middleware.IsAdmin(), dataController.DeleteCity)
	router.GET("/data/category", dataController.GetAllCategory)
	router.POST("/data/category", middleware.Auth(), middleware.IsAdmin(), dataController.AddCategory)
	router.DELETE("/data/category/:id", middleware.Auth(), middleware.IsAdmin(), dataController.DeleteCategory)

	router.POST("/product", middleware.Auth(), productController.AddProduct)
	router.GET("/product", productController.GetAllProduct)
	router.GET("/product/my-product", middleware.Auth(), productController.GetMyProduct)
	router.GET("/product/id/:id", middleware.Auth(), productController.GetProductById)
	router.GET("/product/account/:id", middleware.Auth(), productController.GetProductByAccount)
	router.GET("/product/category/:path", middleware.Auth(), productController.GetByCategory)
	router.PUT("/product/:id", middleware.Auth(), productController.UpdateProduct)
	router.POST("/product/image/:id", middleware.Auth(), productController.AddProductImage)
	router.PUT("/product/thumbnail", middleware.Auth(), productController.UpdateProductThumbnail)
	router.PUT("/product/publish/:id", middleware.Auth(), productController.UpdatePublishStatus)
	router.PUT("/product/status/:id", middleware.Auth(), productController.UpdateSoldStatus)
	router.DELETE("/product/:id", middleware.Auth(), productController.DeleteProduct)
	router.DELETE("/product/image", middleware.Auth(), productController.DeleteProductImage)

	
	router.POST("/transaction", middleware.Auth(), transactionController.AddTransaction)
	router.GET("/transaction/offer", middleware.Auth(), transactionController.GetTransactionByAccount)
	router.GET("/transaction/my-transaction", middleware.Auth(), transactionController.GetMyTransaction)
	router.GET("/transaction/id/:id", middleware.Auth(), transactionController.GetTransactionDetail)
	router.GET("/transaction/product/:id", middleware.Auth(), transactionController.GetTransactionByProduct)
	router.GET("/transaction/buyer/:id", middleware.Auth(), transactionController.GetTransactionByBuyer)
	router.PUT("/transaction", middleware.Auth(), transactionController.UpdatePriceOffer)
	router.PUT("/transaction/id/:id", middleware.Auth(), transactionController.UpdateTransactionStatus)


	return router
}