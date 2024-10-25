// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
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
        "/billowners": {
            "get": {
                "description": "Get all Bill Owner",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "billowners"
                ],
                "summary": "Get all Bill Owner",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.BillOwner"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create new bill owner",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "billowners"
                ],
                "summary": "Create new bill owner",
                "parameters": [
                    {
                        "description": "Bill Owner data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.BillOwner"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.BillOwner"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    }
                }
            }
        },
        "/billowners/{id}": {
            "get": {
                "description": "Get Bill Owner by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "billowners"
                ],
                "summary": "Get a Bill Owner by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "data",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.BillOwner"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Edit Bill Owner",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "billowners"
                ],
                "summary": "Edit Bill Owner",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Bill Owner ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Edit Bill Owner Data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.BillOwner"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.BillOwner"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete bill owner by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "billowners"
                ],
                "summary": "Delete bill owner by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Bill Owner ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.BillOwner"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    }
                }
            }
        },
        "/bills": {
            "get": {
                "description": "Get all Bills from table",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bills"
                ],
                "summary": "Get all bills",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Bill"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    }
                }
            }
        },
        "/bills/{id}": {
            "get": {
                "description": "Get bill by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bills"
                ],
                "summary": "Get a bill by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bill ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Bill"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update bill by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bills"
                ],
                "summary": "Update a bill by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bill ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Bill Data",
                        "name": "bill",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Bill"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Bill"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete bill by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bills"
                ],
                "summary": "Delete a bill by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bill ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/helpers.ErrResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "helpers.ErrResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "models.Bill": {
            "type": "object",
            "properties": {
                "billData": {
                    "description": "has many BillData",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.BillData"
                    }
                },
                "billOwner": {
                    "$ref": "#/definitions/models.BillOwner"
                },
                "billOwnerId": {
                    "description": "belongs to a BillOwner",
                    "type": "integer"
                },
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "rawImage": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "models.BillData": {
            "type": "object",
            "properties": {
                "bill": {
                    "$ref": "#/definitions/models.Bill"
                },
                "billId": {
                    "description": "belongs to a Bill",
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "discount": {
                    "type": "number"
                },
                "id": {
                    "type": "integer"
                },
                "misc": {
                    "type": "string"
                },
                "service": {
                    "type": "number"
                },
                "storeName": {
                    "type": "string"
                },
                "subTotal": {
                    "type": "number"
                },
                "tax": {
                    "type": "number"
                },
                "total": {
                    "type": "number"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "models.BillOwner": {
            "type": "object",
            "properties": {
                "bankAccount": {
                    "type": "string"
                },
                "bill": {
                    "description": "has many Bill",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Bill"
                    }
                },
                "contact": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Split Bill API",
	Description:      "Split Bill Swagger Documentation",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
