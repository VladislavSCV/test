package seed

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"w/backend/internal/models"
)

func Run(db *gorm.DB) error {
	login := "Admin26"
	var n int64
	if err := db.Model(&models.User{}).Where("login = ?", login).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("Demo20"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return db.Create(&models.User{
		FullName: "Администратор", Phone: "+70000000000",
		Email: "admin@local", Login: login, PasswordHash: string(hash), IsAdmin: true,
	}).Error
}
