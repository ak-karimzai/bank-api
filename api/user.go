package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/ak-karimzai/ak-karimzai/simpleb/internal/db"
	"github.com/ak-karimzai/ak-karimzai/simpleb/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username            string    `json:"username"`
	Fullname            string    `json:"fullname"`
	Email               string    `json:"email"`
	PasswordLastChanged time.Time `json:"password_last_changed"`
	CreatedAt           time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:            user.Username,
		Fullname:            user.Fullname,
		Email:               user.Email,
		PasswordLastChanged: user.PasswordLastChanged,
		CreatedAt:           user.CreatedAt,
	}
}

func (server *Server) createUser(
	ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,
			errResponse(err))
		return
	}

	hashedPwd, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError, errResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPwd,
		Fullname:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(
					http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(
			http.StatusInternalServerError, errResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
