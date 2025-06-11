package v1

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type itemGroup struct {
	OnlyName    string `json:"onlyName"`
	ItemGroupID uint   `json:"itemGroupID"`
	Title       string `json:"title"`
}

func (a *itemGroup) GetList(c *gin.Context) {

	userInfo, _ := base.GetCurrentUserInfo(c)
	itemIconGroups := []models.ItemIconGroup{}

	err := global.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Order("sort ,created_at").Where("user_id=?", userInfo.ID).Find(&itemIconGroups).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return err
		}

		// 判断分组是否为空，为空将自动创建默认分组
		if len(itemIconGroups) == 0 {
			defaultGroup := models.ItemIconGroup{
				Title:  "APP",
				UserId: userInfo.ID,
				Icon:   "material-symbols:ad-group-outline",
			}
			if err := tx.Create(&defaultGroup).Error; err != nil {
				apiReturn.ErrorDatabase(c, err.Error())
				return err
			}

			// 并将当前账号下所有无分组的图标更新到当前组
			if err := tx.Model(&models.ItemIcon{}).Where("user_id=?", userInfo.ID).Update("item_icon_group_id", defaultGroup.ID).Error; err != nil {
				apiReturn.ErrorDatabase(c, err.Error())
				return err
			}

			itemIconGroups = append(itemIconGroups, defaultGroup)
		}

		// 返回 nil 提交事务
		return nil
	})

	itemGroups := []itemGroup{}
	for _, group := range itemIconGroups {
		itemGroups = append(itemGroups, itemGroup{
			OnlyName:    group.Title,
			ItemGroupID: group.ID,
			Title:       group.Title,
		})
	}

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	} else {
		apiReturn.SuccessListData(c, itemGroups, int64(len(itemGroups)))
	}
}
