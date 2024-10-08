// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/forgot-password": {
            "post": {
                "description": "it send code to your email address",
                "tags": [
                    "auth"
                ],
                "summary": "Forgot Password",
                "parameters": [
                    {
                        "description": "enough",
                        "name": "token",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.GetUSerByEmailReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid date",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "error while reading from server",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "it generates new access and refresh tokens",
                "tags": [
                    "auth"
                ],
                "summary": "login user",
                "parameters": [
                    {
                        "description": "username and password",
                        "name": "userinfo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.LoginReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "tokens",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid date",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "error while reading from server",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "create new users",
                "tags": [
                    "auth"
                ],
                "summary": "Register user",
                "parameters": [
                    {
                        "description": "User info",
                        "name": "info",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.RegisterReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.RegisterRes"
                        }
                    },
                    "400": {
                        "description": "Invalid data",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/reset-password": {
            "post": {
                "description": "it Reset your Password",
                "tags": [
                    "auth"
                ],
                "summary": "Reset Password",
                "parameters": [
                    {
                        "description": "enough",
                        "name": "token",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.ResetPassReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid date",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "error while reading from server",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/change-password": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update User Profile by token",
                "tags": [
                    "user"
                ],
                "summary": "Update User Profile",
                "parameters": [
                    {
                        "description": "all",
                        "name": "userinfo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.ResetPasswordReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid date",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "error while reading from server",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/logout": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "logout",
                "tags": [
                    "user"
                ],
                "summary": "logout user",
                "responses": {
                    "200": {
                        "description": "message",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/photo": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Api for upload a new photo",
                "consumes": [
                    "multipart/form-data"
                ],
                "tags": [
                    "user"
                ],
                "summary": "UploadMediaUser",
                "parameters": [
                    {
                        "type": "file",
                        "description": "UploadMediaForm",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/profile": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get User Profile by token",
                "tags": [
                    "user"
                ],
                "summary": "Get User Profile",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.GetUserResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid date",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "error while reading from server",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update User Profile by token",
                "tags": [
                    "user"
                ],
                "summary": "Update User Profile",
                "parameters": [
                    {
                        "description": "all",
                        "name": "userinfo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.UpdateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid date",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "error while reading from server",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "user.GetUSerByEmailReq": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "user.GetUserResponse": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fullname": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "photo": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "user.LoginReq": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "user.RegisterReq": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fullname": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "user.RegisterRes": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "user.ResetPassReq": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "user.ResetPasswordReq": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "newpassword": {
                    "type": "string"
                },
                "oldpassword": {
                    "type": "string"
                }
            }
        },
        "user.UpdateUserRequest": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "fullname": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "photo": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "API Gateway",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
