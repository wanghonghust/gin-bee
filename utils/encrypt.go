package utils

import "golang.org/x/crypto/bcrypt"

func Password(pwd string) (pwdHash string, err error) {
	var byteHashPwd []byte
	byteHashPwd, err = bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	pwdHash = string(byteHashPwd)
	return
}

func PasswordVerify(pwd string, encryptPwd string) (status bool) {
	err := bcrypt.CompareHashAndPassword([]byte(encryptPwd), []byte(pwd))
	return err == nil
}
