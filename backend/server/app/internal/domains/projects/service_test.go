package projects

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func setupTestService(t *testing.T) Service {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&Project{}, &Milestone{}, &KanbanColumn{}, &KanbanCard{}, &ProjectDependency{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}

	repo := NewGormRepository(db)
	return NewService(repo, nil)
}

func TestDuplicateProjectCopiesBoardAndMilestones(t *testing.T) {
	svc := setupTestService(t)
	userID := uuid.New()

	project, err := svc.CreateProject(contextWithTODO(), CreateProjectRequest{
		UserID:      userID,
		Name:        "Template",
		Description: "Source project",
		Status:      ProjectStatusPlanning,
		HealthScore: 80,
		MRR:         1200,
		CAC:         200,
		LTV:         4000,
		ChurnRate:   5,
	})
	if err != nil {
		t.Fatalf("create project: %v", err)
	}

	milestone, err := svc.AddMilestone(contextWithTODO(), AddMilestoneRequest{
		ProjectID:   project.ID,
		UserID:      userID,
		Title:       "Kickoff",
		Description: "Initial planning",
		DueDate:     time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("add milestone: %v", err)
	}

	board, err := svc.CreateKanbanColumn(contextWithTODO(), CreateKanbanColumnRequest{
		ProjectID: project.ID,
		UserID:    userID,
		Name:      "Backlog",
	})
	if err != nil {
		t.Fatalf("create column: %v", err)
	}
	if len(board.Columns) != 1 {
		t.Fatalf("expected 1 column, got %d", len(board.Columns))
	}

	columnID := board.Columns[0].Column.ID
	if columnID == uuid.Nil {
		t.Fatalf("expected column id to be set")
	}

	due := time.Now().Add(48 * time.Hour)
	milestoneID := milestone.ID
	if _, err := svc.CreateKanbanCard(contextWithTODO(), CreateKanbanCardRequest{
		ProjectID:   project.ID,
		UserID:      userID,
		ColumnID:    columnID,
		Title:       "Define roadmap",
		Description: "Outline major milestones",
		DueDate:     &due,
		MilestoneID: &milestoneID,
	}); err != nil {
		t.Fatalf("create card: %v", err)
	}

	duplicate, err := svc.DuplicateProject(contextWithTODO(), DuplicateProjectRequest{
		TemplateProjectID: project.ID,
		UserID:            userID,
		Name:              "Template Copy",
		CopyBoard:         true,
		CopyMilestones:    true,
		CopyDependencies:  false,
	})
	if err != nil {
		t.Fatalf("duplicate project: %v", err)
	}

	if duplicate.ID == project.ID {
		t.Fatalf("duplicate should have new id")
	}

	dupBoard, err := svc.GetKanbanBoard(contextWithTODO(), duplicate.ID, userID)
	if err != nil {
		t.Fatalf("get duplicate board: %v", err)
	}

	if len(dupBoard.Columns) != 1 {
		t.Fatalf("expected 1 column in duplicate, got %d", len(dupBoard.Columns))
	}
	if len(dupBoard.Columns[0].Cards) != 1 {
		t.Fatalf("expected 1 card in duplicate column, got %d", len(dupBoard.Columns[0].Cards))
	}

	dupMilestones, err := svc.ListMilestones(contextWithTODO(), duplicate.ID, userID)
	if err != nil {
		t.Fatalf("list duplicate milestones: %v", err)
	}
	if len(dupMilestones) != 1 {
		t.Fatalf("expected 1 milestone in duplicate, got %d", len(dupMilestones))
	}
	if dupMilestones[0].Title != milestone.Title {
		t.Fatalf("expected milestone title %q, got %q", milestone.Title, dupMilestones[0].Title)
	}
}

func TestMoveKanbanCardUpdatesBoard(t *testing.T) {
	svc := setupTestService(t)
	userID := uuid.New()

	project, err := svc.CreateProject(contextWithTODO(), CreateProjectRequest{
		UserID:      userID,
		Name:        "Kanban",
		Description: "Board",
		Status:      ProjectStatusExecuting,
		HealthScore: 75,
	})
	if err != nil {
		t.Fatalf("create project: %v", err)
	}

	columnAResp, err := svc.CreateKanbanColumn(contextWithTODO(), CreateKanbanColumnRequest{ProjectID: project.ID, UserID: userID, Name: "Todo"})
	if err != nil {
		t.Fatalf("create column a: %v", err)
	}
	columnAID := columnAResp.Columns[0].Column.ID
	if columnAID == uuid.Nil {
		t.Fatalf("expected column a id to be set")
	}

	columnBResp, err := svc.CreateKanbanColumn(contextWithTODO(), CreateKanbanColumnRequest{ProjectID: project.ID, UserID: userID, Name: "Doing"})
	if err != nil {
		t.Fatalf("create column b: %v", err)
	}
	columnBID := columnBResp.Columns[1].Column.ID
	if columnBID == uuid.Nil {
		t.Fatalf("expected column b id to be set")
	}

	cardResp, err := svc.CreateKanbanCard(contextWithTODO(), CreateKanbanCardRequest{
		ProjectID:   project.ID,
		UserID:      userID,
		ColumnID:    columnAID,
		Title:       "Task",
		Description: "Do this",
	})
	if err != nil {
		t.Fatalf("create card: %v", err)
	}

	cardID := cardResp.Columns[0].Cards[0].ID
	if cardID == uuid.Nil {
		t.Fatalf("expected card id to be set")
	}

	board, err := svc.MoveKanbanCard(contextWithTODO(), MoveKanbanCardRequest{
		ProjectID:      project.ID,
		UserID:         userID,
		CardID:         cardID,
		TargetColumnID: columnBID,
		TargetPosition: 0,
	})
	if err != nil {
		t.Fatalf("move card: %v", err)
	}

	if len(board.Columns[0].Cards) != 0 {
		t.Fatalf("expected column A empty, got %d cards", len(board.Columns[0].Cards))
	}
	if len(board.Columns[1].Cards) != 1 {
		t.Fatalf("expected column B to have 1 card, got %d", len(board.Columns[1].Cards))
	}
	if board.Columns[1].Cards[0].ColumnID != columnBID {
		t.Fatalf("card column mismatch: expected %s got %s", columnBID, board.Columns[1].Cards[0].ColumnID)
	}
}

func TestBulkUpdateMilestones(t *testing.T) {
	svc := setupTestService(t)
	userID := uuid.New()

	project, err := svc.CreateProject(contextWithTODO(), CreateProjectRequest{
		UserID:      userID,
		Name:        "Milestones",
		Description: "Bulk",
		Status:      ProjectStatusPlanning,
		HealthScore: 60,
	})
	if err != nil {
		t.Fatalf("create project: %v", err)
	}

	milestone, err := svc.AddMilestone(contextWithTODO(), AddMilestoneRequest{
		ProjectID:   project.ID,
		UserID:      userID,
		Title:       "Phase 1",
		Description: "Initial",
		DueDate:     time.Now(),
	})
	if err != nil {
		t.Fatalf("add milestone: %v", err)
	}

	newTitle := "Phase 1 Updated"
	completed := true
	newDue := time.Now().Add(72 * time.Hour)

	updates, err := svc.BulkUpdateMilestones(contextWithTODO(), BulkUpdateMilestonesRequest{
		ProjectID: project.ID,
		UserID:    userID,
		Updates: []MilestoneUpdate{
			{
				MilestoneID: milestone.ID,
				Title:       &newTitle,
				DueDate:     &newDue,
				Completed:   &completed,
			},
		},
	})
	if err != nil {
		t.Fatalf("bulk update milestones: %v", err)
	}

	if len(updates) != 1 {
		t.Fatalf("expected 1 milestone updated, got %d", len(updates))
	}

	updated := updates[0]
	if updated.Title != newTitle {
		t.Fatalf("expected title %q, got %q", newTitle, updated.Title)
	}
	if !updated.Completed {
		t.Fatalf("expected milestone completed")
	}
	if !updated.DueDate.Equal(newDue.UTC()) {
		t.Fatalf("expected due date %v, got %v", newDue.UTC(), updated.DueDate)
	}
}

// contextWithTODO returns a cancellable background context for tests.
func contextWithTODO() context.Context {
	return context.Background()
}
