package model

// Organization role constants for OrganizationMember.
const (
	OrgRoleOwner  = "owner"
	OrgRoleAdmin  = "admin"
	OrgRoleMember = "member"
)

// Organization represents a B2B tenant that can own projects and have members.
type Organization struct {
	Id          int    `json:"id"`
	Name        string `json:"name" gorm:"type:varchar(128);not null;uniqueIndex"`
	Description string `json:"description" gorm:"type:text"`
	OwnerId     int    `json:"owner_id" gorm:"not null;index"`
	Status      int    `json:"status" gorm:"default:1;not null"` // 1=active, 0=disabled
	CreatedAt   int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

// OrganizationMember records a user's membership and role within an organization.
type OrganizationMember struct {
	Id        int    `json:"id"`
	OrgId     int    `json:"org_id" gorm:"not null;uniqueIndex:idx_org_user"`
	UserId    int    `json:"user_id" gorm:"not null;uniqueIndex:idx_org_user"`
	Role      string `json:"role" gorm:"type:varchar(32);not null;default:'member'"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
}

// Project belongs to an Organization and groups API tokens and usage.
type Project struct {
	Id          int    `json:"id"`
	OrgId       int    `json:"org_id" gorm:"not null;index"`
	Name        string `json:"name" gorm:"type:varchar(128);not null"`
	Description string `json:"description" gorm:"type:text"`
	Status      int    `json:"status" gorm:"default:1;not null"` // 1=active, 0=disabled
	CreatedAt   int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

// ── Organization CRUD ─────────────────────────────────────────────────────────

func CreateOrganization(org *Organization) error {
	return DB.Create(org).Error
}

func GetOrganizationById(id int) (*Organization, error) {
	var org Organization
	err := DB.First(&org, "id = ?", id).Error
	return &org, err
}

func GetOrganizationsByOwner(ownerId int) ([]*Organization, error) {
	var orgs []*Organization
	err := DB.Where("owner_id = ?", ownerId).Find(&orgs).Error
	return orgs, err
}

func UpdateOrganization(org *Organization) error {
	return DB.Save(org).Error
}

func DeleteOrganization(id int) error {
	return DB.Delete(&Organization{}, "id = ?", id).Error
}

// ── OrganizationMember CRUD ───────────────────────────────────────────────────

func AddOrganizationMember(m *OrganizationMember) error {
	return DB.Create(m).Error
}

func GetOrganizationMembers(orgId int) ([]*OrganizationMember, error) {
	var members []*OrganizationMember
	err := DB.Where("org_id = ?", orgId).Find(&members).Error
	return members, err
}

func GetUserOrganizations(userId int) ([]*OrganizationMember, error) {
	var memberships []*OrganizationMember
	err := DB.Where("user_id = ?", userId).Find(&memberships).Error
	return memberships, err
}

func RemoveOrganizationMember(orgId, userId int) error {
	return DB.Where("org_id = ? AND user_id = ?", orgId, userId).
		Delete(&OrganizationMember{}).Error
}

// ── Project CRUD ──────────────────────────────────────────────────────────────

func CreateProject(p *Project) error {
	return DB.Create(p).Error
}

func GetProjectById(id int) (*Project, error) {
	var p Project
	err := DB.First(&p, "id = ?", id).Error
	return &p, err
}

func GetProjectsByOrg(orgId int) ([]*Project, error) {
	var projects []*Project
	err := DB.Where("org_id = ?", orgId).Find(&projects).Error
	return projects, err
}

func UpdateProject(p *Project) error {
	return DB.Save(p).Error
}

func DeleteProject(id int) error {
	return DB.Delete(&Project{}, "id = ?", id).Error
}
