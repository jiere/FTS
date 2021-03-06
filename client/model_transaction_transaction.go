/*
 * FTS API
 *
 * This is FTS(Funds Transfer Service) server API document.
 *
 * API version: 1.0
 * Contact: dickrj@163.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import "time"

type TransactionTransaction struct {
	ID        uint       `json:"ID,omitempty"`
	CreatedAt time.Time  `json:"CreatedAt,omitempty"`
	UpdatedAt time.Time  `json:"UpdatedAt,omitempty"`
	DeletedAt *time.Time `json:"DeletedAt,omitempty"`
	DstName   string     `json:"dst_name,omitempty"`
	Money     int32      `json:"money,omitempty"`
	SrcName   string     `json:"src_name,omitempty"`
	Status    string     `json:"status,omitempty"`
}
