package repositories

// import (
// 	"fmt"

// 	"gorm.io/gorm"

// 	model "github.com/rehanazhar/cashier-app/models"
// )

// type UserRepository struct {
// 	db *gorm.DB
// }

// func NewUserRepository(db *gorm.DB) UserRepository {
// 	return UserRepository{db}
// }

// func (u *UserRepository) AddUser(user model.User) error {
// 	if err := u.db.Create(&user).Error; err != nil {
// 		// return any error will rollback
// 		return err
// 	}
// 	return nil
// }

// func (u *UserRepository) UserAvail(cred model.User) error {
// 	var data model.User
// 	result := u.db.Table("users").Select("*").Where("username = ?", cred.Username).Where("password = ?", cred.Password).Scan(&data)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	if data == (model.User{}) {
// 		return fmt.Errorf("user tidak available")
// 	}
// 	return nil
// }

// func (u *UserRepository) CheckPassLength(pass string) bool {

// 	return len(pass) <= 5
// }

// func (u *UserRepository) CheckPassAlphabet(pass string) bool {
// 	for _, charVariable := range pass {
// 		if (charVariable < 'a' || charVariable > 'z') && (charVariable < 'A' || charVariable > 'Z') {
// 			return false
// 		}
// 	}
// 	return true
// }
