// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://127.0.0.1/docs/index.html",
        "contact": {
            "name": "追梦小窝",
            "url": "http://github.com/iszmxw",
            "email": "mail@54zm.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/api/user/login": {
            "post": {
                "description": "提交注册的邮箱和密码即可登录",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "登录接口"
                ],
                "summary": "登录接口",
                "parameters": [
                    {
                        "type": "string",
                        "description": "邮箱",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "登录密码",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response._LoginHandler"
                        }
                    }
                }
            }
        },
        "/v1/api/user/send_email_register": {
            "post": {
                "description": "发送注册邮件",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "发送注册邮件"
                ],
                "summary": "发送注册邮件",
                "parameters": [
                    {
                        "type": "string",
                        "description": "邮箱",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response._OK"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "response._LoginHandler": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "msg": {
                    "type": "string"
                },
                "reqId": {
                    "type": "string"
                },
                "result": {
                    "type": "object",
                    "properties": {
                        "token": {
                            "description": "登录获取的token",
                            "type": "string"
                        },
                        "uid": {
                            "description": "用户id",
                            "type": "integer"
                        }
                    }
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "response._OK": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "msg": {
                    "type": "string"
                },
                "reqId": {
                    "type": "string"
                },
                "result": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
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
	Version:     "3.0",
	Host:        "127.0.0.1",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "用户端接口服务",
	Description: "3.0版本，基于之前的2.0改造的",
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
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
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
	swag.Register("swagger", &s{})
}
