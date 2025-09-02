package model

import (
	"github.com/pachirode/pkg/authn"
	"gorm.io/gorm"

	"github.com/pachirode/go_blog_study/internal/pkg/rid"
)

// AfterCreate 在创建数据库记录之后生成 postID.
func (m *PostM) AfterCreate(tx *gorm.DB) error {
	m.PostID = rid.PostID.NewResourceID(uint64(m.ID))

	return tx.Save(m).Error
}

// BeforeCreate 在创建数据库记录之前加密明文密码.
func (m *UserM) BeforeCreate(tx *gorm.DB) error {
	// Encrypt the user password.
	var err error
	m.Password, err = authn.Encrypt(m.Password)
	if err != nil {
		return err
	}

	return nil
}

// AfterCreate 在创建数据库记录之后生成 userID.
func (m *UserM) AfterCreate(tx *gorm.DB) error {
	m.UserID = rid.UserID.NewResourceID(uint64(m.ID))

	return tx.Save(m).Error
}
