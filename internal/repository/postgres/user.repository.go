package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/fares7elsadek/Social-Golang/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
    db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
    return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	userQuery := `
		INSERT INTO users(username, email, password,is_active)
		VALUES($1, $2, $3, $4)
		RETURNING id
	`

	err = tx.QueryRow(ctx, userQuery,
		user.Username,
		user.Email,
		user.Password,
		user.IsActive,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	roleQuery := `
		INSERT INTO user_roles(user_id, role_id)
		VALUES($1, $2)
	`

	for _, role := range user.Roles {
		roleID, ok := domain.RoleMapping[string(role)]
		if !ok {
			return fmt.Errorf("invalid role: %s", role)
		}

		_, err := tx.Exec(ctx, roleQuery,
			user.ID,
			roleID,
		)

		if err != nil {
			return fmt.Errorf("insert role: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	query := `
        SELECT id, username, email, created_at, updated_at 
        FROM users 
        WHERE id = $1
    `

	user := &domain.User{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err,pgx.ErrNoRows){
            return nil, domain.ErrNotFound
        }
		return nil, fmt.Errorf("userRepository.GetUserByID: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
    query := `
        SELECT id, username, email, created_at, updated_at 
        FROM users 
        WHERE email = $1
    `

	user := &domain.User{}

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if(errors.Is(err,pgx.ErrNoRows)){
			return nil,domain.ErrNotFound
		}
		return nil, fmt.Errorf("userRepository.GetUserByID: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
    query := `
        SELECT id, username, email, created_at, updated_at 
        FROM users 
        WHERE username = $1
    `

	user := &domain.User{}

	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if(errors.Is(err,pgx.ErrNoRows)){
			return nil,domain.ErrNotFound
		}
		return nil, fmt.Errorf("userRepository.GetUserByUsername: %w", err)
	}

	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
        UPDATE users
        SET username = $1,
            email = $2,
			password = $3
        WHERE id = $4
    `

	result, err := r.db.Exec(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("userRepository.UpdateUser: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
    query := `
        DELETE FROM users
        WHERE id=$1
    `

	result, err := r.db.Exec(ctx, query,
		id,
	)

	if err != nil {
		return fmt.Errorf("userRepository.DeleteUser: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}