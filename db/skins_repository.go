package db

import (
	"context"
	"database/sql"
	"fmt"
)

type SkinsRepository interface {
	GetAllColors(ctx context.Context) ([]Color, error)
	GetAllPatterns(ctx context.Context) ([]Pattern, error)
	CheckColorExists(ctx context.Context, colorID string) (bool, error)
	CheckPatternsExists(ctx context.Context, patternID string) (bool, error)
	CheckAccountHasSkin(ctx context.Context, accountID string) (bool, error)
	GetAccountSkin(ctx context.Context, accountID string) (*Skin, error)
	SetAccountSkin(ctx context.Context, accountID string, colorID string, patternID string) error
}

type skinsRepository struct {
	dbConn *sql.DB
}

type Color struct {
	ID    string
	Color string
}

type Pattern struct {
	ID     string
	Source string
	Type   string
}

type Skin struct {
	ColorID   string
	PatternID string
}

func NewSkinsRepository(dbConn *sql.DB) SkinsRepository {
	return &skinsRepository{dbConn}
}

func (sr skinsRepository) GetAllColors(ctx context.Context) ([]Color, error) {
	rows, err := sr.dbConn.QueryContext(ctx, "SELECT id, color FROM skin_colors")
	if err != nil {
		return nil, err
	}

	colors := make([]Color, 0)

	for rows.Next() {
		color := Color{}

		err := rows.Scan(&color.ID, &color.Color)
		if err != nil {
			break
		}

		colors = append(colors, color)
	}

	if closeErr := rows.Close(); closeErr != nil {
		return nil, closeErr
	}

	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return colors, nil
}

func (sr skinsRepository) GetAllPatterns(ctx context.Context) ([]Pattern, error) {
	rows, err := sr.dbConn.QueryContext(ctx, "SELECT id, source, type FROM skin_patterns")
	if err != nil {
		return nil, err
	}

	patterns := make([]Pattern, 0)

	for rows.Next() {
		pattern := Pattern{}

		err := rows.Scan(&pattern.ID, &pattern.Source, &pattern.Type)
		if err != nil {
			break
		}

		patterns = append(patterns, pattern)
	}

	if closeErr := rows.Close(); closeErr != nil {
		return nil, closeErr
	}

	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return patterns, nil
}

func (sr skinsRepository) CheckColorExists(ctx context.Context, colorID string) (bool, error) {
	row := sr.dbConn.QueryRowContext(
		ctx,
		"SELECT count(id) FROM skin_colors WHERE id=$1",
		colorID,
	)

	var count uint8

	err := row.Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	if count >= 1 {
		return true, nil
	}

	return false, nil
}

func (sr skinsRepository) CheckPatternsExists(ctx context.Context, patternID string) (bool, error) {
	row := sr.dbConn.QueryRowContext(
		ctx,
		"SELECT count(id) FROM skin_patterns WHERE id=$1",
		patternID,
	)

	var count uint8

	err := row.Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	if count >= 1 {
		return true, nil
	}

	return false, nil
}

func (sr skinsRepository) CheckAccountHasSkin(ctx context.Context, accountID string) (bool, error) {
	row := sr.dbConn.QueryRowContext(
		ctx,
		"SELECT count(account) FROM account_skin WHERE account=$1",
		accountID,
	)

	var count uint8

	err := row.Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	if count >= 1 {
		return true, nil
	}

	return false, nil
}

func (sr skinsRepository) GetAccountSkin(ctx context.Context, accountID string) (*Skin, error) {
	row := sr.dbConn.QueryRowContext(ctx, "SELECT color, pattern FROM account_skin WHERE account=$1", accountID)

	var skin Skin

	err := row.Scan(&skin.ColorID, &skin.PatternID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &skin, nil
}

func (sr skinsRepository) SetAccountSkin(ctx context.Context, accountID string, colorID string, patternID string) error {
	accountHasSkin, err := sr.CheckAccountHasSkin(ctx, accountID)
	if err != nil {
		return err
	}

	if accountHasSkin {
		result, err := sr.dbConn.ExecContext(ctx, "UPDATE account_skin SET color=$1, pattern=$2 WHERE account=$3", colorID, patternID, accountID)
		if err != nil {
			return err
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rows != 1 {
			return fmt.Errorf("expected to affect 1 row, affected %d", rows)
		}
	} else {
		stmt, err := sr.dbConn.PrepareContext(ctx, "INSERT INTO account_skin (account, color, pattern) VALUES ($1, $2, $3)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(accountID, colorID, patternID)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
