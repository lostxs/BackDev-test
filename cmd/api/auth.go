package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lostxs/BackDev-test/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type CreateTokenResponse struct {
	*store.User
	AccessToken string `json:"access_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

func (a *app) createTokensHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		a.badRequestException(w, r, fmt.Errorf("user_id not provided"))
		return
	}

	ipAddress := r.RemoteAddr

	user, err := a.getUser(r.Context(), userID)
	if err != nil {
		switch err {
		case store.ErrInvalidUserID:
			a.badRequestException(w, r, err)
		case store.ErrUserNotFound:
			a.unauthorizedException(w, r, err)
		default:
			a.internalServerException(w, r, err)
		}
		return
	}

	accessToken, err := a.createAccessToken(user.ID, ipAddress, a.config.auth.accessToken.exp)
	if err != nil {
		a.internalServerException(w, r, err)
		return
	}

	refreshToken, err := a.authenticator.GenerateRefreshToken()
	if err != nil {
		a.internalServerException(w, r, err)
		return
	}

	hash, err := hashValue(refreshToken)
	if err != nil {
		a.internalServerException(w, r, err)
		return
	}

	session := &store.Session{
		UserID:           user.ID,
		RefreshTokenHash: string(hash),
	}
	if err := a.store.Sessions.Upsert(r.Context(), session); err != nil {
		a.internalServerException(w, r, err)
		return
	}

	setCookie(w, "refresh_token", refreshToken, "/", true)

	if err := a.jsonResponse(w, http.StatusOK, CreateTokenResponse{
		User:        user,
		AccessToken: accessToken,
	}); err != nil {
		a.internalServerException(w, r, err)
	}
}

func (a *app) refreshTokensHandler(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil || refreshCookie.Value == "" {
		a.unauthorizedException(w, r, fmt.Errorf("refresh token not provided"))
		return
	}

	refreshToken := refreshCookie.Value

	newIPAddress := r.RemoteAddr

	userID := r.Context().Value(userCtx).(*store.User).ID

	session, err := a.store.Sessions.GetByUserID(r.Context(), userID)
	if err != nil {
		switch err {
		case store.ErrSessionNotFound:
			a.unauthorizedException(w, r, err)
		case store.ErrInvalidUserID:
			a.badRequestException(w, r, err)
		default:
			a.internalServerException(w, r, err)
		}
		return
	}

	if !compareHashAndValue(session.RefreshTokenHash, refreshToken) {
		a.unauthorizedException(w, r, fmt.Errorf("refresh token mismatch"))
		return
	}

	userEmail := r.Context().Value(userCtx).(*store.User).Email
	tokenIPAddress := r.Context().Value(ipAddressCtx).(string)

	if tokenIPAddress != newIPAddress {
		mockSendEmail(userEmail, "IP address mismatch", "your IP address has changed")
	}

	newAccessToken, err := a.createAccessToken(session.UserID, newIPAddress, a.config.auth.accessToken.exp)
	if err != nil {
		a.internalServerException(w, r, err)
		return
	}

	newRefreshToken, err := a.authenticator.GenerateRefreshToken()
	if err != nil {
		a.internalServerException(w, r, err)
		return
	}

	hash, err := hashValue(newRefreshToken)
	if err != nil {
		a.internalServerException(w, r, err)
		return
	}

	session.RefreshTokenHash = string(hash)
	if err := a.store.Sessions.Upsert(r.Context(), session); err != nil {
		a.internalServerException(w, r, err)
		return
	}

	setCookie(w, "refresh_token", newRefreshToken, "/", true)

	if err := a.jsonResponse(w, http.StatusOK, RefreshResponse{
		AccessToken: newAccessToken,
	}); err != nil {
		a.internalServerException(w, r, err)
	}
}

func (a *app) createAccessToken(userID string, ipAddress string, exp time.Duration) (string, error) {
	accessClaims := jwt.MapClaims{
		"sub":        userID,
		"ip_address": ipAddress,
		"exp":        time.Now().Add(exp).Unix(),
	}

	return a.authenticator.GenerateAccessToken(accessClaims)
}

func mockSendEmail(to, subject, body string) {
	log.Printf("mock email sent to %s:\nSubject: %s\nBody: %s", to, subject, body)
}

func hashValue(value string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func compareHashAndValue(hash, value string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(value)) == nil
}

func setCookie(w http.ResponseWriter, name, value string, path string, httpOnly bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		HttpOnly: httpOnly,
	})
}
