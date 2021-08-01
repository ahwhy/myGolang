package hash

import (
	"fmt"

	"github.com/infraboard/mcube/types/ftime"
	"golang.org/x/crypto/bcrypt"
)

// User info
type User struct {
	Account  string
	Password *Password
}

// NewHashedPassword 生产hash后的密码对象
func NewHashedPassword(password string) (*Password, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	return &Password{
		Password: string(bytes),
		CreateAt: ftime.Now().Timestamp(),
		UpdateAt: ftime.Now().Timestamp(),
	}, nil
}

type Password struct {
	// hash过后的密码
	Password string
	// 密码创建时间
	CreateAt int64
	// 密码更新时间
	UpdateAt int64
	// 密码需要被重置
	NeedReset bool
	// 需要重置的原因
	ResetReason string
	// 历史密码
	History []string
	// 是否过期
	IsExpired bool
}

// Update 更新密码
func (p *Password) Update(new *Password, maxHistory uint, needReset bool) {
	p.rotaryHistory(maxHistory)
	p.Password = new.Password
	p.NeedReset = needReset
	p.UpdateAt = ftime.Now().Timestamp()
	if !needReset {
		p.ResetReason = ""
	}
}

// IsHistory 检测是否是历史密码
func (p *Password) IsHistory(password string) bool {
	for _, pass := range p.History {
		err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
		if err == nil {
			return true
		}
	}

	return false
}

// HistoryCount 保存了几个历史密码
func (p *Password) HistoryCount() int {
	return len(p.History)
}

func (p *Password) rotaryHistory(maxHistory uint) {
	if uint(p.HistoryCount()) < maxHistory {
		p.History = append(p.History, p.Password)
	} else {
		remainHistry := p.History[:maxHistory]
		p.History = []string{p.Password}
		p.History = append(p.History, remainHistry...)
	}
}

// CheckPassword 判断password 是否正确
// 输入的是明文: 123456
// 需要对比的Hash: $2a$10$ofPPqZ3m37Kp9ROK4ForAOXc5w6SsMKoJ9puCOgIO9yEFFknpYcsO
func (p *Password) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("user or password not connrect")
	}
	return nil
}
