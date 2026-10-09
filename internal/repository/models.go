package repository

import "time"

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` //do not expose the hash
}

type Video struct {
	ID        string    `json:"id"`         // UUID o identificador único
	Title     string    `json:"title"`      // Título del video
	Status    string    `json:"status"`     // "processing", "ready", "failed"
	HLSPath   string    `json:"hls_path"`   // Ruta al archivo .m3u8 (ej: /uploads/id/index.m3u8)
	UserID    int64     `json:"user_id"`    // ID del usuario creador
	CreatedAt time.Time `json:"created_at"` // Fecha de subida
}
