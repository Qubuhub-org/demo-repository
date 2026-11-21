package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/argon2"
)

const (
	sessionCookie = "qubuhub_session"
	sessionTTL = 12 * time.Hour
	maxBodyBytes = 8 << 10
)

type Service struct {
	db *sql.DB
	secureCookie bool
	mu sync.Mutex
	attempts map[string][]time.Time
}

type credentials struct { Email string `json:"email"`; Password string `json:"password"` }
type user struct { ID string `json:"id"`; Email string `json:"email"`; CreatedAt time.Time `json:"createdAt"` }

func New(ctx context.Context, dsn string, secure bool) (*Service, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil { return nil, err }
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30*time.Minute)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil { db.Close(); return nil, err }
	return &Service{db: db, secureCookie: secure, attempts: make(map[string][]time.Time)}, nil
}
func (s *Service) Close() { _ = s.db.Close() }

func (s *Service) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, http.StatusOK, map[string]string{"status":"ok"}) })
	mux.HandleFunc("POST /api/auth/register", s.register)
	mux.HandleFunc("POST /api/auth/login", s.login)
	mux.HandleFunc("POST /api/auth/logout", s.logout)
	mux.HandleFunc("GET /api/auth/me", s.me)
}

func (s *Service) register(w http.ResponseWriter, r *http.Request) {
	if !s.validMutation(r) { http.Error(w, "invalid request origin", http.StatusForbidden); return }
	if !s.allow(r) { http.Error(w, "too many attempts", http.StatusTooManyRequests); return }
	var input credentials
	if !decode(w,r,&input) { return }
	email := normalizeEmail(input.Email)
	if _, err := mail.ParseAddress(email); err != nil || len(email)>254 { http.Error(w,"valid email required",http.StatusBadRequest); return }
	if len(input.Password)<12 || len(input.Password)>128 { http.Error(w,"password must be 12–128 characters",http.StatusBadRequest); return }
	salt:=make([]byte,16); if _,err:=rand.Read(salt); err!=nil { http.Error(w,"internal error",500); return }
	hash:=argon2.IDKey([]byte(input.Password),salt,3,64*1024,2,32)
	encoded:=base64.RawStdEncoding.EncodeToString(append(salt,hash...))
	var created user
	err:=s.db.QueryRowContext(r.Context(),`INSERT INTO users(email,password_hash) VALUES($1,$2) RETURNING id::text,email,created_at`,email,encoded).Scan(&created.ID,&created.Email,&created.CreatedAt)
	if err!=nil { http.Error(w,"account could not be created",http.StatusConflict); return }
	if err:=s.createSession(w,r,created.ID); err!=nil { http.Error(w,"could not create session",500); return }
	writeJSON(w,http.StatusCreated,map[string]any{"user":created})
}

func (s *Service) login(w http.ResponseWriter, r *http.Request) {
	if !s.validMutation(r) { http.Error(w,"invalid request origin",403); return }
	if !s.allow(r) { http.Error(w,"too many attempts",429); return }
	var input credentials
	if !decode(w,r,&input) { return }
	var u user
	var encoded string
	err:=s.db.QueryRowContext(r.Context(),`SELECT id::text,email,created_at,password_hash FROM users WHERE email=$1`,normalizeEmail(input.Email)).Scan(&u.ID,&u.Email,&u.CreatedAt,&encoded)
	if err!=nil || !verifyPassword(encoded,input.Password) { http.Error(w,"invalid email or password",http.StatusUnauthorized); return }
	if err:=s.createSession(w,r,u.ID); err!=nil { http.Error(w,"could not create session",500); return }
	writeJSON(w,http.StatusOK,map[string]any{"user":u})
}

func (s *Service) logout(w http.ResponseWriter,r *http.Request) {
	if !s.validMutation(r) { http.Error(w,"invalid request origin",403); return }
	if c,err:=r.Cookie(sessionCookie); err==nil { digest:=sha256.Sum256([]byte(c.Value)); _,_=s.db.ExecContext(r.Context(),`DELETE FROM sessions WHERE token_hash=$1`,digest[:]) }
	http.SetCookie(w,&http.Cookie{Name:sessionCookie,Value:"",Path:"/",HttpOnly:true,Secure:s.secureCookie,SameSite:http.SameSiteStrictMode,MaxAge:-1})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) me(w http.ResponseWriter,r *http.Request) {
	c,err:=r.Cookie(sessionCookie); if err!=nil { http.Error(w,"unauthorized",401); return }
	digest:=sha256.Sum256([]byte(c.Value)); var u user
	err=s.db.QueryRowContext(r.Context(),`SELECT u.id::text,u.email,u.created_at FROM sessions s JOIN users u ON u.id=s.user_id WHERE s.token_hash=$1 AND s.expires_at>now()`,digest[:]).Scan(&u.ID,&u.Email,&u.CreatedAt)
	if err!=nil { http.Error(w,"unauthorized",401); return }
	writeJSON(w,200,map[string]any{"user":u})
}

func (s *Service) createSession(w http.ResponseWriter,r *http.Request,userID string) error {
	token:=make([]byte,32); if _,err:=rand.Read(token); err!=nil { return err }
	raw:=base64.RawURLEncoding.EncodeToString(token); digest:=sha256.Sum256([]byte(raw)); expires:=time.Now().Add(sessionTTL)
	if _,err:=s.db.ExecContext(r.Context(),`INSERT INTO sessions(user_id,token_hash,expires_at) VALUES($1,$2,$3)`,userID,digest[:],expires); err!=nil { return err }
	http.SetCookie(w,&http.Cookie{Name:sessionCookie,Value:raw,Path:"/",HttpOnly:true,Secure:s.secureCookie,SameSite:http.SameSiteStrictMode,Expires:expires,MaxAge:int(sessionTTL.Seconds())})
	return nil
}
func (s *Service) validMutation(r *http.Request) bool {
	origin:=r.Header.Get("Origin"); if origin=="" { return false }
	host:=r.Host
	return strings.EqualFold(strings.TrimPrefix(strings.TrimPrefix(origin,"https://"),"http://"),host)
}
func (s *Service) allow(r *http.Request) bool {
	key:=r.RemoteAddr; now:=time.Now(); s.mu.Lock(); defer s.mu.Unlock()
	recent:=s.attempts[key][:0]; for _,t:=range s.attempts[key] { if now.Sub(t)<time.Minute { recent=append(recent,t) } }
	if len(recent)>=10 { s.attempts[key]=recent; return false }
	s.attempts[key]=append(recent,now); return true
}
func decode(w http.ResponseWriter,r *http.Request,v any) bool {
	r.Body=http.MaxBytesReader(w,r.Body,maxBodyBytes); defer r.Body.Close()
	d:=json.NewDecoder(r.Body); d.DisallowUnknownFields()
	if err:=d.Decode(v); err!=nil { http.Error(w,"invalid JSON request",400); return false }; return true
}
func verifyPassword(encoded,password string) bool {
	raw,err:=base64.RawStdEncoding.DecodeString(encoded); if err!=nil || len(raw)!=48 { return false }
	salt:=raw[:16]; expected:=raw[16:]; actual:=argon2.IDKey([]byte(password),salt,3,64*1024,2,32)
	return subtleCompare(actual,expected)
}
func subtleCompare(a,b []byte) bool { if len(a)!=len(b) { return false }; var diff byte; for i:=range a { diff |= a[i]^b[i] }; return diff==0 }
func normalizeEmail(v string) string { return strings.ToLower(strings.TrimSpace(v)) }
func writeJSON(w http.ResponseWriter,status int,v any) { w.Header().Set("Content-Type","application/json"); w.WriteHeader(status); _=json.NewEncoder(w).Encode(v) }

