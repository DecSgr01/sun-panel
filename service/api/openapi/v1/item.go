package v1

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/siteFavicon"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type item struct {
	OnlyName          string `json:"onlyName"`
	Title             string `json:"title"`
	Url               string `json:"url"`
	LanUrl            string `json:"lanUrl"`
	IconUrl           string `json:"iconUrl"`
	Desc              string `json:"description"`
	ItemGroupID       int    `json:"itemGroupID"`
	ItemGroupOnlyName string `json:"itemGroupOnlyName"`
	IsSaveIcon        bool   `json:"isSaveIcon"`
}

func (a *item) Create(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	// type Request
	item := item{}
	itemIcon := models.ItemIcon{}
	if err := c.ShouldBindBodyWith(&item, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if item.ItemGroupID == 0 && item.ItemGroupOnlyName == "" {
		// apiReturn.Error(c, "Group is mandatory")
		apiReturn.ErrorParamFomat(c, "Group is mandatory")
		return
	}

	if item.IsSaveIcon {
		// 下载网站图标
		var err error
		parsedURL, err := url.Parse(item.IconUrl)
		if err != nil {
			apiReturn.Error(c, err.Error())
			return
		}
		item.IconUrl, err = GetSiteFavicon(item.IconUrl)
		if err != nil {
			apiReturn.Error(c, err.Error())
			return
		}

		// 保存到数据库
		ext := path.Ext(item.IconUrl)
		mFile := models.File{}
		if _, err := mFile.AddFile(userInfo.ID, parsedURL.Host, ext, item.IconUrl); err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
		itemIcon.IconJson = `{"itemType":2,"src":"` + item.IconUrl[1:] + `","backgroundColor":"#2a2a2a6b"}`
	} else {
		itemIcon.IconJson = `{"itemType":2,"src":"` + item.IconUrl + `","backgroundColor":"#2a2a2a6b"}`
	}

	itemIcon.Title = item.Title
	itemIcon.Url = item.Url
	itemIcon.LanUrl = item.LanUrl
	itemIcon.Description = item.Desc
	itemIcon.Private = 0
	itemIcon.OpenMethod = 2
	itemIcon.ItemIconGroupId = item.ItemGroupID
	itemIcon.UserId = userInfo.ID
	itemIcon.Sort = 9999
	global.Db.Create(&itemIcon)

	apiReturn.SuccessData(c, itemIcon)
}

// 支持获取并直接下载对方网站图标到服务器
func GetSiteFavicon(IconUrl string) (string, error) {
	parsedURL, err := url.Parse(IconUrl)
	if err != nil {
		return "", err
	}

	protocol := parsedURL.Scheme
	global.Logger.Debug("protocol:", protocol)
	global.Logger.Debug("IconUrl:", IconUrl)

	// 如果URL以双斜杠（//）开头，则使用当前页面协议
	if strings.HasPrefix(IconUrl, "//") {
		IconUrl = protocol + "://" + IconUrl[2:]
	} else if !strings.HasPrefix(IconUrl, "http://") && !strings.HasPrefix(IconUrl, "https://") {
		// 如果URL既不以http://开头也不以https://开头，则默认为http协议
		IconUrl = "http://" + IconUrl
	}
	global.Logger.Debug("IconUrl:", IconUrl)
	// 去除图标的get参数
	{
		parsedIcoURL, err := url.Parse(IconUrl)
		if err != nil {
			return "", err
		}
		IconUrl = parsedIcoURL.Scheme + "://" + parsedIcoURL.Host + parsedIcoURL.Path
	}
	global.Logger.Debug("IconUrl:", IconUrl)

	// 生成保存目录
	configUpload := global.Config.GetValueString("base", "source_path")
	savePath := fmt.Sprintf("%s/%d/%d/%d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())
	isExist, _ := cmn.PathExists(savePath)
	if !isExist {
		os.MkdirAll(savePath, os.ModePerm)
	}

	// 下载
	var imgInfo *os.File
	{
		var err error
		if imgInfo, err = siteFavicon.DownloadImage(IconUrl, savePath, 1024*1024); err != nil {
			return "", err
		}
	}
	return imgInfo.Name(), nil

}
