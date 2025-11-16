package ideas

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ChangeTypeCreated = "created"
	ChangeTypeEdited  = "edited"
	ChangeTypeMoved   = "moved"
	ChangeTypeBulk    = "bulk"
)

// Idea represents a graph node in the ideas canvas.
type Idea struct {
	ID          uuid.UUID      `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"column:user_id;type:uuid;index;not null;uniqueIndex:idx_ideas_user_slug" json:"userId"`
	Title       string         `gorm:"column:title;size:160;not null" json:"title"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	Slug        string         `gorm:"column:slug;size:200;uniqueIndex:idx_ideas_user_slug" json:"slug"`
	PosX        float64        `gorm:"column:pos_x;not null" json:"posX"`
	PosY        float64        `gorm:"column:pos_y;not null" json:"posY"`
	Color       string         `gorm:"column:color;size:16" json:"color"`
	ProjectID   *uuid.UUID     `gorm:"column:project_id;type:uuid;index" json:"projectId,omitempty"`
	Version     int            `gorm:"column:version;not null;default:1" json:"version"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
}

// IdeaLink represents a relationship between two ideas/projects.
type IdeaLink struct {
	ID            uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	SourceIdeaID  uuid.UUID `gorm:"column:source_idea_id;type:uuid;index;not null" json:"sourceIdeaId"`
	TargetIdeaID  uuid.UUID `gorm:"column:target_idea_id;type:uuid;index;not null" json:"targetIdeaId"`
	Relation      string    `gorm:"column:relation;size:64;not null" json:"relation"`
	Weight        float64   `gorm:"column:weight;default:1" json:"weight"`
	Bidirectional bool      `gorm:"column:bidirectional;default:false" json:"bidirectional"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"createdAt"`
}

// NewIdea constructs a new idea node.
func NewIdea(userID uuid.UUID, title, description string, posX, posY float64, color string, projectID *uuid.UUID) (*Idea, error) {
	idea := &Idea{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Slug:        "",
		PosX:        posX,
		PosY:        posY,
		Color:       strings.TrimSpace(color),
		ProjectID:   projectID,
		Version:     1,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	idea.Slug = generateIdeaSlug(idea.Title)

	return idea, idea.Validate()
}

// Validate ensures idea invariants.
func (i *Idea) Validate() error {
	if i == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilIdea)
	}

	if i.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyIdeaID)
	}

	if i.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if strings.TrimSpace(i.Title) == "" {
		return NewDomainError(ErrCodeInvalidTitle, ErrEmptyTitle)
	}

	if strings.TrimSpace(i.Slug) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyIdeaSlug)
	}

	return nil
}

// Touch updates the updated at timestamp.
func (i *Idea) Touch() {
	i.UpdatedAt = time.Now().UTC()
}

// Move updates canvas coordinates.
func (i *Idea) Move(posX, posY float64) {
	i.PosX = posX
	i.PosY = posY
	i.Version++
	i.Touch()
}

// UpdateDetails updates textual metadata.
func (i *Idea) UpdateDetails(title, description, color string) error {
	if title != "" {
		i.Title = strings.TrimSpace(title)
	}
	if description != "" {
		i.Description = strings.TrimSpace(description)
	}
	if color != "" {
		i.Color = strings.TrimSpace(color)
	}
	i.Version++
	return i.Validate()
}

// NewIdeaLink creates a connection between ideas.
func NewIdeaLink(userID, sourceID, targetID uuid.UUID, relation string, weight float64, bidirectional bool) (*IdeaLink, error) {
	link := &IdeaLink{
		ID:            uuid.New(),
		UserID:        userID,
		SourceIdeaID:  sourceID,
		TargetIdeaID:  targetID,
		Relation:      strings.TrimSpace(relation),
		Weight:        weight,
		Bidirectional: bidirectional,
		CreatedAt:     time.Now().UTC(),
	}

	return link, link.Validate()
}

// Validate ensures IdeaLink invariants.
func (l *IdeaLink) Validate() error {
	if l == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilLink)
	}

	if l.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyLinkID)
	}

	if l.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if l.SourceIdeaID == uuid.Nil || l.TargetIdeaID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidRelation, ErrEmptyRelationNodes)
	}

	if l.SourceIdeaID == l.TargetIdeaID {
		return NewDomainError(ErrCodeInvalidRelation, ErrSelfRelation)
	}

	if strings.TrimSpace(l.Relation) == "" {
		return NewDomainError(ErrCodeInvalidRelation, ErrEmptyRelationLabel)
	}

	return nil
}

// IdeaVersion captures a snapshot of idea metadata.
type IdeaVersion struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	IdeaID      uuid.UUID `gorm:"column:idea_id;type:uuid;index;not null" json:"ideaId"`
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	EditorID    uuid.UUID `gorm:"column:editor_id;type:uuid;index;not null" json:"editorId"`
	Version     int       `gorm:"column:version;index;not null" json:"version"`
	Title       string    `gorm:"column:title;size:160;not null" json:"title"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	PosX        float64   `gorm:"column:pos_x;not null" json:"posX"`
	PosY        float64   `gorm:"column:pos_y;not null" json:"posY"`
	Color       string    `gorm:"column:color;size:16" json:"color"`
	ChangeType  string    `gorm:"column:change_type;size:32;not null" json:"changeType"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
}

// NewIdeaVersion builds a snapshot representation.
func NewIdeaVersion(idea *Idea, editorID uuid.UUID, changeType string) *IdeaVersion {
	if changeType == "" {
		changeType = ChangeTypeEdited
	}
	return &IdeaVersion{
		ID:          uuid.New(),
		IdeaID:      idea.ID,
		UserID:      idea.UserID,
		EditorID:    editorID,
		Version:     idea.Version,
		Title:       idea.Title,
		Description: idea.Description,
		PosX:        idea.PosX,
		PosY:        idea.PosY,
		Color:       idea.Color,
		ChangeType:  strings.ToLower(strings.TrimSpace(changeType)),
		CreatedAt:   time.Now().UTC(),
	}
}

var ideaSlugSanitizer = regexp.MustCompile(`[^a-z0-9]+`)

func generateIdeaSlug(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = ideaSlugSanitizer.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "idea"
	}
	return slug
}

// IdeaCollaborator tracks shared access to an idea canvas.
type IdeaCollaborator struct {
	ID             uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	OwnerID        uuid.UUID `gorm:"column:owner_id;type:uuid;not null;uniqueIndex:idx_owner_collaborator" json:"ownerId"`
	CollaboratorID uuid.UUID `gorm:"column:collaborator_id;type:uuid;not null;uniqueIndex:idx_owner_collaborator" json:"collaboratorId"`
	Role           string    `gorm:"column:role;size:32;not null" json:"role"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// NewIdeaCollaborator constructs a collaborator entry.
func NewIdeaCollaborator(ownerID, collaboratorID uuid.UUID, role string) (*IdeaCollaborator, error) {
	entry := &IdeaCollaborator{
		ID:             uuid.New(),
		OwnerID:        ownerID,
		CollaboratorID: collaboratorID,
		Role:           strings.ToLower(strings.TrimSpace(role)),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
	return entry, entry.Validate()
}

// Validate ensures collaborator invariants.
func (c *IdeaCollaborator) Validate() error {
	if c == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilCollaborator)
	}
	if c.OwnerID == uuid.Nil || c.CollaboratorID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if c.OwnerID == c.CollaboratorID {
		return NewDomainError(ErrCodeInvalidCollaborator, ErrSelfCollaborator)
	}
	if c.Role == "" {
		c.Role = "editor"
	}
	return nil
}

// ConnectionDirection represents the direction of a connection between nodes.
type ConnectionDirection string

const (
	DirectionNorth ConnectionDirection = "north"
	DirectionSouth ConnectionDirection = "south"
	DirectionEast  ConnectionDirection = "east"
	DirectionWest  ConnectionDirection = "west"
)

// IdeaNode represents a node within an idea's canvas.
type IdeaNode struct {
	ID          uuid.UUID      `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	IdeaID      uuid.UUID      `gorm:"column:idea_id;type:uuid;index;not null" json:"ideaId"`
	Title       string         `gorm:"column:title;size:160;not null" json:"title"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	PosX        float64        `gorm:"column:pos_x;not null" json:"posX"`
	PosY        float64        `gorm:"column:pos_y;not null" json:"posY"`
	Width       float64        `gorm:"column:width;default:200" json:"width"`
	Height      float64        `gorm:"column:height;default:100" json:"height"`
	Color       string         `gorm:"column:color;size:16" json:"color"`
	Type        string         `gorm:"column:type;size:64;default:default" json:"type"`
	Version     int            `gorm:"column:version;not null;default:1" json:"version"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
}

// NewIdeaNode constructs a new node within an idea canvas.
func NewIdeaNode(ideaID uuid.UUID, title, description string, posX, posY, width, height float64, color, nodeType string) (*IdeaNode, error) {
	node := &IdeaNode{
		ID:          uuid.New(),
		IdeaID:      ideaID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		PosX:        posX,
		PosY:        posY,
		Width:       width,
		Height:      height,
		Color:       strings.TrimSpace(color),
		Type:        strings.TrimSpace(nodeType),
		Version:     1,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	if node.Type == "" {
		node.Type = "default"
	}
	if node.Width <= 0 {
		node.Width = 200
	}
	if node.Height <= 0 {
		node.Height = 100
	}

	return node, node.Validate()
}

// Validate ensures IdeaNode invariants.
func (n *IdeaNode) Validate() error {
	if n == nil {
		return NewDomainError(ErrCodeInvalidPayload, "ideas: node entity is nil")
	}

	if n.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, "ideas: node id cannot be empty")
	}

	if n.IdeaID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, "ideas: idea id cannot be empty")
	}

	if strings.TrimSpace(n.Title) == "" {
		return NewDomainError(ErrCodeInvalidTitle, ErrEmptyTitle)
	}

	return nil
}

// Touch updates the updated at timestamp.
func (n *IdeaNode) Touch() {
	n.UpdatedAt = time.Now().UTC()
}

// Move updates canvas coordinates.
func (n *IdeaNode) Move(posX, posY float64) {
	n.PosX = posX
	n.PosY = posY
	n.Version++
	n.Touch()
}

// UpdateDetails updates textual metadata.
func (n *IdeaNode) UpdateDetails(title, description, color, nodeType string) error {
	if title != "" {
		n.Title = strings.TrimSpace(title)
	}
	if description != "" {
		n.Description = strings.TrimSpace(description)
	}
	if color != "" {
		n.Color = strings.TrimSpace(color)
	}
	if nodeType != "" {
		n.Type = strings.TrimSpace(nodeType)
	}
	n.Version++
	return n.Validate()
}

// Resize updates node dimensions.
func (n *IdeaNode) Resize(width, height float64) {
	if width > 0 {
		n.Width = width
	}
	if height > 0 {
		n.Height = height
	}
	n.Version++
	n.Touch()
}

// IdeaNodeConnection represents a directional connection between two nodes within an idea's canvas.
type IdeaNodeConnection struct {
	ID           uuid.UUID          `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	IdeaID       uuid.UUID          `gorm:"column:idea_id;type:uuid;index;not null" json:"ideaId"`
	SourceNodeID uuid.UUID          `gorm:"column:source_node_id;type:uuid;index;not null" json:"sourceNodeId"`
	TargetNodeID uuid.UUID          `gorm:"column:target_node_id;type:uuid;index;not null" json:"targetNodeId"`
	Direction    ConnectionDirection `gorm:"column:direction;size:16;not null" json:"direction"`
	Label        string             `gorm:"column:label;size:128" json:"label"`
	CreatedAt    time.Time          `gorm:"column:created_at" json:"createdAt"`
}

// NewIdeaNodeConnection creates a connection between two nodes.
func NewIdeaNodeConnection(ideaID, sourceNodeID, targetNodeID uuid.UUID, direction ConnectionDirection, label string) (*IdeaNodeConnection, error) {
	conn := &IdeaNodeConnection{
		ID:           uuid.New(),
		IdeaID:       ideaID,
		SourceNodeID: sourceNodeID,
		TargetNodeID: targetNodeID,
		Direction:    direction,
		Label:        strings.TrimSpace(label),
		CreatedAt:    time.Now().UTC(),
	}

	return conn, conn.Validate()
}

// Validate ensures IdeaNodeConnection invariants.
func (c *IdeaNodeConnection) Validate() error {
	if c == nil {
		return NewDomainError(ErrCodeInvalidPayload, "ideas: connection entity is nil")
	}

	if c.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, "ideas: connection id cannot be empty")
	}

	if c.IdeaID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, "ideas: idea id cannot be empty")
	}

	if c.SourceNodeID == uuid.Nil || c.TargetNodeID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidRelation, "ideas: connection requires source and target nodes")
	}

	if c.SourceNodeID == c.TargetNodeID {
		return NewDomainError(ErrCodeInvalidRelation, ErrSelfRelation)
	}

	switch c.Direction {
	case DirectionNorth, DirectionSouth, DirectionEast, DirectionWest:
		// Valid direction
	default:
		return NewDomainError(ErrCodeInvalidRelation, "ideas: invalid connection direction, must be north, south, east, or west")
	}

	return nil
}

// Document represents a text document related to an idea, optionally linked to a node.
type Document struct {
	ID          uuid.UUID      `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	IdeaID      uuid.UUID      `gorm:"column:idea_id;type:uuid;index;not null" json:"ideaId"`
	NodeID      *uuid.UUID     `gorm:"column:node_id;type:uuid;index" json:"nodeId,omitempty"`
	Title       string         `gorm:"column:title;size:200;not null" json:"title"`
	Content     string         `gorm:"column:content;type:text" json:"content"`
	Version     int            `gorm:"column:version;not null;default:1" json:"version"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
}

// NewDocument constructs a new document for an idea.
func NewDocument(ideaID uuid.UUID, nodeID *uuid.UUID, title, content string) (*Document, error) {
	doc := &Document{
		ID:        uuid.New(),
		IdeaID:    ideaID,
		NodeID:    nodeID,
		Title:     strings.TrimSpace(title),
		Content:   strings.TrimSpace(content),
		Version:   1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	return doc, doc.Validate()
}

// Validate ensures Document invariants.
func (d *Document) Validate() error {
	if d == nil {
		return NewDomainError(ErrCodeInvalidPayload, "ideas: document entity is nil")
	}

	if d.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, "ideas: document id cannot be empty")
	}

	if d.IdeaID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, "ideas: idea id cannot be empty")
	}

	if strings.TrimSpace(d.Title) == "" {
		return NewDomainError(ErrCodeInvalidTitle, ErrEmptyTitle)
	}

	return nil
}

// Touch updates the updated at timestamp.
func (d *Document) Touch() {
	d.UpdatedAt = time.Now().UTC()
}

// UpdateContent updates the document content.
func (d *Document) UpdateContent(content string) {
	d.Content = strings.TrimSpace(content)
	d.Version++
	d.Touch()
}

// UpdateDetails updates document title and content.
func (d *Document) UpdateDetails(title, content string) error {
	if title != "" {
		d.Title = strings.TrimSpace(title)
	}
	if content != "" {
		d.Content = strings.TrimSpace(content)
	}
	d.Version++
	return d.Validate()
}
