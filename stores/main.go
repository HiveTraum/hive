package stores

import (
	"auth/config"
	"auth/models"
	"auth/repositories"
	"auth/repositories/inMemoryRepository"
	"auth/repositories/postgresRepository"
	"auth/repositories/redisRepository"
	"context"
	"github.com/go-redis/redis/v7"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
	"time"
)

type IStore interface {

	// store can be used to combine multiple physical storage elements, like postgres, redis, elasticSearch and etc...

	// All store methods

	// Users

	CreateUser(ctx context.Context, password, email, phone string) (int, *models.User)
	GetUser(context context.Context, id uuid.UUID) *models.User
	GetUsers(context context.Context, query repositories.GetUsersQuery) []*models.User
	DeleteUser(ctx context.Context, id uuid.UUID) (int, *models.User)

	// User Views

	GetUsersView(context context.Context, query repositories.GetUsersViewStoreQuery) ([]*models.UserView, *models.PaginationResponse)
	GetUserView(context context.Context, id uuid.UUID) *models.UserView
	CreateOrUpdateUsersView(context context.Context, query repositories.CreateOrUpdateUsersViewStoreQuery) []*models.UserView
	CreateOrUpdateUsersViewByUsersID(context context.Context, id []uuid.UUID) []*models.UserView
	CreateOrUpdateUsersViewByRolesID(context context.Context, id []uuid.UUID) []*models.UserView
	CreateOrUpdateUsersViewByUserID(context context.Context, id uuid.UUID) []*models.UserView
	CreateOrUpdateUsersViewByRoleID(context context.Context, id uuid.UUID) []*models.UserView
	CacheUserView(ctx context.Context, userViews []*models.UserView)

	// Emails

	CreateEmail(ctx context.Context, userId uuid.UUID, value string) (int, *models.Email)
	GetEmail(ctx context.Context, email string) (int, *models.Email)
	CreateEmailConfirmationCode(ctx context.Context, email string, code string, duration time.Duration) *models.EmailConfirmation
	GetEmailConfirmationCode(ctx context.Context, email string) string
	GetRandomCodeForEmailConfirmation() string

	// Passwords

	CreatePassword(ctx context.Context, userId uuid.UUID, value string) (int, *models.Password)
	GetPasswords(ctx context.Context, userId uuid.UUID) []*models.Password
	GetLatestPassword(ctx context.Context, userId uuid.UUID) (int, *models.Password)

	// Phones

	CreatePhone(ctx context.Context, userId uuid.UUID, value string) (int, *models.Phone)
	GetPhone(ctx context.Context, phone string) (int, *models.Phone)
	CreatePhoneConfirmationCode(ctx context.Context, phone string, code string, duration time.Duration) *models.PhoneConfirmation
	GetPhoneConfirmationCode(ctx context.Context, phone string) string
	GetRandomCodeForPhoneConfirmation() string

	// Roles

	CreateRole(context context.Context, title string) (int, *models.Role)
	GetRole(context context.Context, id uuid.UUID) (int, *models.Role)
	GetRoles(context context.Context, query repositories.GetRolesQuery) ([]*models.Role, *models.PaginationResponse)
	GetRoleByTitle(ctx context.Context, title string) (int, *models.Role)
	GetAdminRole(ctx context.Context) (int, *models.Role)

	// User Roles

	CreateUserRole(ctx context.Context, userId uuid.UUID, roleId uuid.UUID) (int, *models.UserRole)
	GetUserRoles(ctx context.Context, query repositories.GetUserRoleQuery) ([]*models.UserRole, *models.PaginationResponse)
	DeleteUserRole(ctx context.Context, id uuid.UUID) (int, *models.UserRole)

	// Secrets

	GetSecret(ctx context.Context, id uuid.UUID) *models.Secret
	GetActualSecret(ctx context.Context) *models.Secret
	CreateSecret(ctx context.Context) *models.Secret

	// Sessions

	CreateSession(ctx context.Context, fingerprint string, userID uuid.UUID, secretID uuid.UUID, userAgent string) (int, *models.Session)
	GetSession(ctx context.Context, fingerprint string, refreshToken string, userID uuid.UUID) *models.Session
}

type DatabaseStore struct {
	db                 *pgxpool.Pool
	cache              *redis.Client
	inMemoryCache      *cache.Cache
	environment        *config.Environment
	postgresRepository postgresRepository.IPostgresRepository
	redisRepository    redisRepository.IRedisRepository
	inMemoryRepository inMemoryRepository.IInMemoryRepository
}

func InitStore(
	db *pgxpool.Pool,
	cache *redis.Client,
	inMemoryCache *cache.Cache,
	environment *config.Environment,
	postgresRepository postgresRepository.IPostgresRepository,
	redisRepository redisRepository.IRedisRepository,
	inMemoryRepository inMemoryRepository.IInMemoryRepository) *DatabaseStore {
	return &DatabaseStore{
		db:                 db,
		cache:              cache,
		inMemoryCache:      inMemoryCache,
		environment:        environment,
		postgresRepository: postgresRepository,
		redisRepository:    redisRepository,
		inMemoryRepository: inMemoryRepository,
	}
}
