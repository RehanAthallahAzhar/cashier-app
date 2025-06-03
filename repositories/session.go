package repositories

// import (
// 	"fmt"
// 	"time"

// 	"github.com/rehanazhar/cashier-app/models"
// 	"gorm.io/gorm"
// )

// type SessionsRepository struct {
// 	db *gorm.DB
// }

// func NewSessionsRepository(db *gorm.DB) SessionsRepository {
// 	return SessionsRepository{db}
// }

// func (u *SessionsRepository) AddSessions(session models.Session) error {
// 	if err := u.db.Table("sessions").Create(&session).Error; err != nil {
// 		// return any error will rollback
// 		return err
// 	}
// 	return nil
// }

// func (u *SessionsRepository) DeleteSessions(tokenTarget string) error {
// 	results := []models.Session{}
// 	rows, err := u.db.Table("sessions").Select("*").Where("deleted_at is NULL").Rows()
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()
// 	for rows.Next() { // Next akan menyiapkan hasil baris berikutnya untuk dibaca dengan metode Scan.
// 		u.db.ScanRows(rows, &results)
// 	}

// 	err = u.db.Table("sessions").Delete(&results).Where("token = ?", tokenTarget).Error
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (u *SessionsRepository) UpdateSessions(session models.Session) error {
// 	// UPDATE sessions SET (token = {token}, expiry = {expiry}) where username = {username}
// 	result := u.db.Table("sessions").Where("username = ?", session.Username).Update("token", session.Token).Update("expiry", session.Expiry)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }

// func (u *SessionsRepository) TokenValidity(token string) (models.Session, error) {
// 	session, err := u.SessionAvailToken(token)
// 	if err != nil {
// 		return models.Session{}, fmt.Errorf("asd")
// 	}

// 	if u.TokenExpired(session) {
// 		err := u.DeleteSessions(token)
// 		if err != nil {
// 			return models.Session{}, err
// 		}
// 		return models.Session{}, fmt.Errorf("Token is Expired!")
// 	}
// 	return session, nil
// }

// func (u *SessionsRepository) SessionAvailName(name string) (models.Session, error) {
// 	var data models.Session
// 	result := u.db.Table("sessions").Select("*").Where("username = ?", name).Scan(&data)
// 	if result.Error != nil {
// 		return models.Session{}, result.Error
// 	}
// 	if data == (models.Session{}) {
// 		return models.Session{}, fmt.Errorf("Session token tidak ditemukan")
// 	}
// 	return data, nil
// }

// func (u *SessionsRepository) SessionAvailToken(token string) (models.Session, error) {
// 	var data models.Session
// 	result := u.db.Table("sessions").Select("*").Where("token = ?", token).Where("deleted_at is null").Scan(&data)
// 	if result.Error != nil {
// 		return models.Session{}, result.Error
// 	}
// 	if data == (models.Session{}) {
// 		return models.Session{}, fmt.Errorf("Session token tidak ditemukan")
// 	}

// 	return data, nil
// }

// func (u *SessionsRepository) TokenExpired(s models.Session) bool {
// 	return s.Expiry.Before(time.Now())
// }
