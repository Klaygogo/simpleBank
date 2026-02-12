package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/Klaygogo/simplebank/db/sqlc"
	util "github.com/Klaygogo/simplebank/util"
	pgtype "github.com/jackc/pgx/v5/pgtype"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"

	"errors"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Fullname string `json:"fullname" binding:"required"`
}

type userResponse struct {
	Username          string
	FullName          string
	Email             string
	CreatedAt         pgtype.Timestamptz
	PasswordChangedAt pgtype.Timestamptz
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}

func (server *Server) createUser(c *gin.Context) {
	var req createUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Email:          req.Email,
		FullName:       req.Fullname,
	}

	user, err := server.store.CreateUser(c, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
			c.JSON(400, errorResponse(err))
			return
		}
	}
	resp := newUserResponse(user)
	c.JSON(http.StatusOK, resp)
}

type getUserRequest struct {
	Username string `uri:"username" binding:"required,min=1"`
}

func (server *Server) getUser(c *gin.Context) {
	var req getUserRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(c, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, user)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken           string       `json:"access_token"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	sessionID             pgtype.UUID  `json:"session_id"`
	User                  userResponse `json:"user"`
}

func (server *Server) loginUser(c *gin.Context) {
	var req loginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := server.store.GetUser(c, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	session, err := server.store.CreateSession(c, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpireAt:     refreshPayload.ExpiredAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	resp := loginUserResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt.Time,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt.Time,
		sessionID:             session.ID,
		User:                  newUserResponse(user),
	}
	c.JSON(http.StatusOK, resp)
}
