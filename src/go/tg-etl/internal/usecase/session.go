package usecase

import (
	"context"
	"tg-etl/internal/models"
	"tg-etl/pkg/slogger/wsl"
)

func (uc *UseCase) CreateSession(ctx context.Context, session models.Session) error {
	err := uc.etlDB.CreateSession(ctx, session)
	if err != nil {
		uc.l.ErrorContext(ctx, "create session",
			wsl.Err(err))
		return err
	}

	//uc.l.DebugContext(ctx, "success create session")
	return nil
}

func (uc *UseCase) GetSessionByID(ctx context.Context, id uint) (*models.Session, error) {
	session, err := uc.etlDB.GetSessionByID(ctx, id)
	if err != nil {
		uc.l.ErrorContext(ctx, "get session by id", wsl.Int("id", int(id)),
			wsl.Err(err))
		return nil, err
	}
	//uc.l.DebugContext(ctx, "success got session", wsl.Int("id", int(id)))
	return session, nil
}

func (uc *UseCase) GetSessions(ctx context.Context, limit, offset int) ([]models.Session, error) {
	sessions, err := uc.etlDB.GetSessions(ctx)
	if err != nil {
		uc.l.ErrorContext(ctx, "get sessions",
			wsl.Err(err))
		return nil, err
	}
	//uc.l.DebugContext(ctx, "success got sessions", wsl.Int("count", len(sessions)))
	return sessions, nil
}
