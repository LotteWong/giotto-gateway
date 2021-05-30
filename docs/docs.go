// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/tokens/jwt": {
            "post": {
                "description": "生成认证令牌",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "认证接口"
                ],
                "summary": "生成认证令牌接口",
                "operationId": "/tokens/jwt",
                "parameters": [
                    {
                        "description": "generate jwt request body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.JwtReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/common_middleware.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.JwtRes"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "common_middleware.Response": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "errmsg": {
                    "type": "string"
                },
                "errno": {
                    "type": "integer"
                },
                "stack": {
                    "type": "object"
                },
                "trace_id": {
                    "type": "object"
                }
            }
        },
        "dto.JwtReq": {
            "type": "object",
            "properties": {
                "permission": {
                    "description": "读写权限",
                    "type": "string"
                },
                "type": {
                    "description": "令牌类型",
                    "type": "string"
                }
            }
        },
        "dto.JwtRes": {
            "type": "object",
            "properties": {
                "expire_at": {
                    "description": "失效时间",
                    "type": "integer"
                },
                "permission": {
                    "description": "读写权限",
                    "type": "string"
                },
                "token": {
                    "description": "令牌内容",
                    "type": "string"
                },
                "type": {
                    "description": "令牌类型",
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}