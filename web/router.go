package web

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"wechat-mall-backend/handler/cms"
	"wechat-mall-backend/handler/portal"
	"wechat-mall-backend/service"
)

func NewRouter(app *App) *mux.Router {
	router := mux.NewRouter()
	conf := app.Conf
	s := service.NewService(app.Conf)
	cmsHandler := cms.NewHandler(conf, s)
	portalHandler := portal.NewHandler(conf, s)

	registerHandler(router, cmsHandler, portalHandler)
	return router
}

func registerHandler(router *mux.Router, cmsHandler *cms.Handler, portalHandler *portal.Handler) {
	mw := &Middleware{}
	chain := alice.New(mw.LoggingHandler, mw.RecoverPanic, mw.CORSHandler, mw.ValidateAuthToken)
	router.Handle("/api/wxapp/login", chain.ThenFunc(portalHandler.Login)).Methods("GET").Queries("code", "{code}")
	router.Handle("/api/wxapp/user-info", chain.ThenFunc(portalHandler.UserInfo)).Methods("GET")
	router.Handle("/api/wxapp/auth-phone", chain.ThenFunc(portalHandler.AuthPhone)).Methods("POST")
	router.Handle("/api/wxapp/auth-info", chain.ThenFunc(portalHandler.AuthUserInfo)).Methods("POST")
	router.Handle("/api/home/banner", chain.ThenFunc(portalHandler.HomeBanner)).Methods("GET").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/api/home/grid", chain.ThenFunc(portalHandler.GetGridCategoryList)).Methods("GET").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/api/category/list", chain.ThenFunc(portalHandler.GetSubCategoryList)).Methods("GET")
	router.Handle("/api/goods/list", chain.ThenFunc(portalHandler.GetGoodsList)).Methods("GET").Queries("k", "{k}").Queries("s", "{s}").Queries("c", "{c}").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/api/goods/detail", chain.ThenFunc(portalHandler.GetGoodsDetail)).Methods("GET").Queries("id", "{id}")
	router.Handle("/api/cart/list", chain.ThenFunc(portalHandler.GetCartGoodsList)).Methods("GET").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/api/cart/add", chain.ThenFunc(portalHandler.AddCartGoods)).Methods("POST")
	router.Handle("/api/cart/edit", chain.ThenFunc(portalHandler.EditCartGoods)).Methods("POST")
	router.Handle("/api/coupon/list", chain.ThenFunc(portalHandler.GetCouponList)).Methods("GET").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/api/coupon/take", chain.ThenFunc(portalHandler.TakeCoupon)).Methods("POST")
	router.Handle("/api/user/coupon/list", chain.ThenFunc(portalHandler.GetUserCouponList)).Methods("GET").Queries("status", "{status}").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/api/user/coupon", chain.ThenFunc(portalHandler.DoDeleteCouponLog)).Methods("DELETE").Queries("id", "{id}")
	router.Handle("/api/user/address/list", chain.ThenFunc(portalHandler.GetAddressList)).Methods("GET").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/api/user/address/edit", chain.ThenFunc(portalHandler.EditAddress)).Methods("POST")
	router.Handle("/api/user/address", chain.ThenFunc(portalHandler.GetAddress)).Methods("GET").Queries("id", "{id}")
	router.Handle("/api/user/address", chain.ThenFunc(portalHandler.DoDeleteAddress)).Methods("DELETE").Queries("id", "{id}")
	router.Handle("/api/user/default_address", chain.ThenFunc(portalHandler.GetDefaultAddress)).Methods("GET")
	router.Handle("/api/placeorder", chain.ThenFunc(portalHandler.PlaceOrder)).Methods("POST")
	router.Handle("/api/order/list", chain.ThenFunc(portalHandler.GetOrderList)).Methods("GET").Queries("status", "{status}").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/api/order/detail", chain.ThenFunc(portalHandler.GetOrderDetail)).Methods("GET").Queries("orderNo", "{orderNo}")
	router.Handle("/api/order/cancel", chain.ThenFunc(portalHandler.CancelOrder)).Methods("PUT").Queries("id", "{id}")
	router.Handle("/api/order", chain.ThenFunc(portalHandler.DeleteOrder)).Methods("DELETE").Queries("id", "{id}")
	router.Handle("/api/order/confirm_goods", chain.ThenFunc(portalHandler.ConfirmTakeGoods)).Methods("PUT").Queries("id", "{id}")
	router.Handle("/api/order/refund", chain.ThenFunc(portalHandler.RefundApply)).Methods("PUT")
	router.Handle("/api/browse/list", chain.ThenFunc(portalHandler.UserBrowseHistory)).Methods("GET").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/api/browse/clear", chain.ThenFunc(portalHandler.ClearBrowseHistory)).Methods("POST")
	router.Handle("/wxpay/notify", chain.ThenFunc(portalHandler.WxPayNotify)).Methods("POST")
	router.Handle("/cms/user/login", chain.ThenFunc(cmsHandler.Login)).Methods("POST", "OPTIONS")
	router.Handle("/cms/user/refresh", chain.ThenFunc(cmsHandler.Refresh)).Methods("GET", "OPTIONS")
	router.Handle("/cms/user/info", chain.ThenFunc(cmsHandler.GetUserInfo)).Methods("GET", "OPTIONS")
	router.Handle("/cms/user/change_password", chain.ThenFunc(cmsHandler.DoChangePassword)).Methods("PUT", "OPTIONS")
	router.Handle("/cms/admin/users", chain.ThenFunc(cmsHandler.GetUserList)).Methods("GET", "OPTIONS").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/cms/admin/user", chain.ThenFunc(cmsHandler.DoEditUser)).Methods("POST", "OPTIONS")
	router.Handle("/cms/admin/user", chain.ThenFunc(cmsHandler.GetUser)).Methods("GET", "OPTIONS").Queries("id", "{id}")
	router.Handle("/cms/admin/user", chain.ThenFunc(cmsHandler.DoDeleteCMSUser)).Methods("DELETE", "OPTIONS").Queries("id", "{id}")
	router.Handle("/cms/admin/reset_password", chain.ThenFunc(cmsHandler.DoResetCMSUserPassword)).Methods("POST", "OPTIONS")
	router.Handle("/cms/admin/groups", chain.ThenFunc(cmsHandler.GetUserGroupList)).Methods("GET", "OPTIONS").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/cms/admin/group", chain.ThenFunc(cmsHandler.DoEditUserGroup)).Methods("POST", "OPTIONS")
	router.Handle("/cms/admin/group", chain.ThenFunc(cmsHandler.GetUserGroup)).Methods("GET", "OPTIONS").Queries("id", "{id}")
	router.Handle("/cms/admin/group", chain.ThenFunc(cmsHandler.DoDeleteUserGroup)).Methods("DELETE", "OPTIONS").Queries("id", "{id}")
	router.Handle("/cms/admin/authority", chain.ThenFunc(cmsHandler.GetModuleList)).Methods("GET", "OPTIONS")
	router.Handle("/cms/banner/list", chain.ThenFunc(cmsHandler.GetBannerList)).Methods("GET", "OPTIONS").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/cms/banner/{id:[0-9]+}", chain.ThenFunc(cmsHandler.GetBanner)).Methods("GET", "OPTIONS")
	router.Handle("/cms/banner/edit", chain.ThenFunc(cmsHandler.DoEditBanner)).Methods("POST", "OPTIONS")
	router.Handle("/cms/banner/{id:[0-9]+}", chain.ThenFunc(cmsHandler.DoDeleteBanner)).Methods("DELETE", "OPTIONS")
	router.Handle("/cms/category/list", chain.ThenFunc(cmsHandler.GetCategoryList)).Methods("GET", "OPTIONS").Queries("pid", "{pid}").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/cms/category/{id:[0-9]+}", chain.ThenFunc(cmsHandler.GetCategoryById)).Methods("GET", "OPTIONS")
	router.Handle("/cms/category/edit", chain.ThenFunc(cmsHandler.DoEditCategory)).Methods("POST", "OPTIONS")
	router.Handle("/cms/category/{id:[0-9]+}", chain.ThenFunc(cmsHandler.DoDeleteCategory)).Methods("DELETE", "OPTIONS")
	router.Handle("/cms/category/all", chain.ThenFunc(cmsHandler.GetChooseCategory)).Methods("GET", "OPTIONS")
	router.Handle("/cms/grid_category/list", chain.ThenFunc(cmsHandler.GetGridCategoryList)).Methods("GET", "OPTIONS").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/cms/grid_category/{id:[0-9]+}", chain.ThenFunc(cmsHandler.GetGridCategory)).Methods("GET", "OPTIONS")
	router.Handle("/cms/grid_category/edit", chain.ThenFunc(cmsHandler.DoEditGridCategory)).Methods("POST", "OPTIONS")
	router.Handle("/cms/grid_category/{id:[0-9]+}", chain.ThenFunc(cmsHandler.DoDeleteGridCategory)).Methods("DELETE", "OPTIONS")
	router.Handle("/cms/spec/list", chain.ThenFunc(cmsHandler.GetSpecificationList)).Methods("GET", "OPTIONS").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/cms/spec/{id:[0-9]+}", chain.ThenFunc(cmsHandler.GetSpecification)).Methods("GET", "OPTIONS")
	router.Handle("/cms/spec/edit", chain.ThenFunc(cmsHandler.DoEditSpecification)).Methods("POST", "OPTIONS")
	router.Handle("/cms/spec/{id:[0-9]+}", chain.ThenFunc(cmsHandler.DoDeleteSpecification)).Methods("DELETE", "OPTIONS")
	router.Handle("/cms/spec/attr/list", chain.ThenFunc(cmsHandler.GetSpecificationAttrList)).Methods("GET", "OPTIONS").Queries("specId", "{specId}")
	router.Handle("/cms/spec/attr/{id:[0-9]+}", chain.ThenFunc(cmsHandler.GetSpecificationAttr)).Methods("GET", "OPTIONS")
	router.Handle("/cms/spec/attr/edit", chain.ThenFunc(cmsHandler.DoEditSpecificationAttr)).Methods("POST", "OPTIONS")
	router.Handle("/cms/spec/attr/{id:[0-9]+}", chain.ThenFunc(cmsHandler.DoDeleteSpecificationAttr)).Methods("DELETE", "OPTIONS")
	router.Handle("/cms/goods/list", chain.ThenFunc(cmsHandler.GetGoodsList)).Methods("GET", "OPTIONS").Queries("k", "{k}").Queries("c", "{c}").Queries("o", "{o}").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/cms/goods/{id:[0-9]+}", chain.ThenFunc(cmsHandler.GetGoods)).Methods("GET", "OPTIONS")
	router.Handle("/cms/goods/edit", chain.ThenFunc(cmsHandler.DoEditGoods)).Methods("POST", "OPTIONS")
	router.Handle("/cms/goods/{id:[0-9]+}", chain.ThenFunc(cmsHandler.DoDeleteGoods)).Methods("DELETE", "OPTIONS")
	router.Handle("/cms/goods/spec", chain.ThenFunc(cmsHandler.GetGoodsSpecList)).Methods("GET", "OPTIONS").Queries("id", "{id}")
	router.Handle("/cms/goods/all", chain.ThenFunc(cmsHandler.GetChooseCategoryGoods)).Methods("GET", "OPTIONS")
	router.Handle("/cms/sku/list", chain.ThenFunc(cmsHandler.GetSKUList)).Methods("GET", "OPTIONS").Queries("page", "{page}").Queries("size", "{size}").Queries("goodsId", "{goodsId}").Queries("k", "{k}").Queries("o", "{o}")
	router.Handle("/cms/sku/{id:[0-9]+}", chain.ThenFunc(cmsHandler.GetSKU)).Methods("GET", "OPTIONS")
	router.Handle("/cms/sku/edit", chain.ThenFunc(cmsHandler.DoEditSKU)).Methods("POST", "OPTIONS")
	router.Handle("/cms/sku/{id:[0-9]+}", chain.ThenFunc(cmsHandler.DoDeleteSKU)).Methods("DELETE", "OPTIONS")
	router.Handle("/cms/coupon/list", chain.ThenFunc(cmsHandler.GetCouponList)).Methods("GET", "OPTIONS").Queries("page", "{page}").Queries("size", "{size}")
	router.Handle("/cms/coupon/{id:[0-9]+}", chain.ThenFunc(cmsHandler.GetCoupon)).Methods("GET", "OPTIONS")
	router.Handle("/cms/coupon/edit", chain.ThenFunc(cmsHandler.DoEditCoupon)).Methods("POST", "OPTIONS")
	router.Handle("/cms/coupon/{id:[0-9]+}", chain.ThenFunc(cmsHandler.DoDeleteCoupon)).Methods("DELETE", "OPTIONS")
	router.Handle("/cms/oss/policy-token", chain.ThenFunc(cmsHandler.GetOSSPolicyToken)).Methods("GET", "OPTIONS").Queries("dir", "{dir}")
	router.Handle("/cms/market_metrics", chain.ThenFunc(cmsHandler.GetMarketMetrics)).Methods("GET", "OPTIONS")
	router.Handle("/cms/order/order_statement", chain.ThenFunc(cmsHandler.GetSaleTableData)).Methods("GET", "OPTIONS").Queries("page", "{page}").Queries("size", "{size}")
}
