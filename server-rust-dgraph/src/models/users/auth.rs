use serde::{Deserialize, Serialize};
use crate::models::users::_enum::{ Role };
use crate::models::commons::verify::{ SlimVerified };

// 登录接口返回给客户端的信息
#[derive(Debug, Default, Serialize, Deserialize, Clone)]
pub struct AuthPayload {
    pub token: String,
    pub user: LoggedUser
}

// 存储于jwt/cookie的登录用户信息
#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SlimUser {
    pub uid: String,
    pub username: String,
    pub role: Role
}

// 返回给客户端的登录用户信息 (密码用于登录时验证密码是否正确，然后会过滤掉)
#[derive(Debug, Default, Serialize, Deserialize, Clone)]
pub struct LoggedUserWithPassword {
    pub uid: String,
    pub name: String,
    pub username: String,
    pub password: String,
    pub role: Role,
    pub avatar_url: String,
    pub is_verified: bool,
    pub lang: String,
    pub description: Option<String>,
    pub verified: Option<SlimVerified>
}

// 返回给客户端的登录用户信息
#[derive(Debug, Default, Serialize, Deserialize, Clone)]
pub struct LoggedUser {
    pub uid: String,
    pub name: String,
    pub username: String,
    pub role: Role,
    pub avatar_url: String,
    pub is_verified: bool,
    pub lang: String,
    pub description: Option<String>,
    pub verified: Option<SlimVerified>
}
