package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123"

	// 生成哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("❌ 生成哈希失败:", err)
		return
	}

	fmt.Println("✅ 密码:", password)
	fmt.Println("✅ 哈希:", string(hash))

	// 验证
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		fmt.Println("❌ 验证失败:", err)
	} else {
		fmt.Println("✅ 验证成功")
	}
}
