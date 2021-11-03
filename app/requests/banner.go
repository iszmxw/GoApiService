package requests

type BannerType struct {
	Type string `json:"type" form:"type"` // 广告位置：1-pc 2-h5 3-app
}
