package model

import "testing"

// validates M9-F01 — Organization struct fields
func TestOrganizationFields(t *testing.T) {
	org := Organization{
		Name:        "Acme Corp",
		Description: "Test org",
		OwnerId:     1,
		Status:      1,
	}
	if org.Name != "Acme Corp" {
		t.Errorf("Name = %q, want 'Acme Corp'", org.Name)
	}
	if org.OwnerId != 1 {
		t.Errorf("OwnerId = %d, want 1", org.OwnerId)
	}
	if org.Status != 1 {
		t.Errorf("Status = %d, want 1", org.Status)
	}
}

// validates M9-F02 — OrganizationMember role constants and struct
func TestOrganizationMemberRoles(t *testing.T) {
	if OrgRoleOwner != "owner" {
		t.Errorf("OrgRoleOwner = %q, want 'owner'", OrgRoleOwner)
	}
	if OrgRoleAdmin != "admin" {
		t.Errorf("OrgRoleAdmin = %q, want 'admin'", OrgRoleAdmin)
	}
	if OrgRoleMember != "member" {
		t.Errorf("OrgRoleMember = %q, want 'member'", OrgRoleMember)
	}

	m := OrganizationMember{OrgId: 1, UserId: 42, Role: OrgRoleOwner}
	if m.OrgId != 1 || m.UserId != 42 || m.Role != OrgRoleOwner {
		t.Errorf("OrganizationMember fields incorrect: %+v", m)
	}
}

// validates M9-F03 — Project belongs to Organization via OrgId
func TestProjectBelongsToOrg(t *testing.T) {
	p := Project{
		OrgId:       7,
		Name:        "Alpha Project",
		Description: "First project",
		Status:      1,
	}
	if p.OrgId != 7 {
		t.Errorf("OrgId = %d, want 7", p.OrgId)
	}
	if p.Name != "Alpha Project" {
		t.Errorf("Name = %q, want 'Alpha Project'", p.Name)
	}
}

// validates proof: org→project and user→org relationships are representable
func TestOrgProjectUserRelationships(t *testing.T) {
	org := Organization{Id: 1, Name: "TestOrg", OwnerId: 10}
	member := OrganizationMember{OrgId: org.Id, UserId: 10, Role: OrgRoleOwner}
	project := Project{OrgId: org.Id, Name: "TestProject"}

	if member.OrgId != org.Id {
		t.Error("member.OrgId should match org.Id")
	}
	if project.OrgId != org.Id {
		t.Error("project.OrgId should match org.Id")
	}
}
