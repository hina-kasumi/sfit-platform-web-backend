package entities

type RoleEnum string

const (
	RoleEnumAdmin  RoleEnum = "ADMIN"
	RoleEnumHead   RoleEnum = "HEAD"
	RoleEnumVice   RoleEnum = "VICE"
	RoleEnumMember RoleEnum = "MEMBER"
	RoleEnumUser   RoleEnum = "USER"
)

type Role struct {
	ID RoleEnum `gorm:"type:varchar;primaryKey"`
}
