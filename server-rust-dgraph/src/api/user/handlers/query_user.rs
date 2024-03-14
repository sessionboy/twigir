use actix_web::{ web, HttpRequest };
use std::collections::HashMap;
use crate::api::user::dgraph::user as dgraph;
use crate::models::users as user_model;
use crate::lib::dgraph::{ DgraphClient };
use crate::lib::res::{ Res, ServiceResult };
use json_gettext::JSONGetText;
use crate::lib::auth::{ AuthUser, GuardAuthUser };
use crate::utils::parse::{ get_accept_language };
use serde_json::json;

/*
  1，功能：我的个人主页信息
  2，path: /users/me
  3，logged_user：当前登录用户信息
*/ 
pub async fn me(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser
) -> ServiceResult {
    let client = &data.into_inner();
    let user_id = logged_user.user_id();
    let mut vars = HashMap::new(); 
    vars.insert("$uid", user_id.as_str());

    let _user: Option<user_model::item::User> = client
        .find_one(dgraph::USER_ME, vars)
        .await?;

    if _user.is_none() {
        let lang = get_accept_language(&req);
        let msg = get_text!(lang_ctx, lang, "user.err.unauthorized").unwrap();
        return Res::err(json!({ "message": msg })); 
    }
    
    let user = _user.unwrap();
    Res::ok(json!({ "data": user }))
}


/*
  1，功能：{id}的个人主页信息
  2，path: /users/{id}
  3，user_id：目标用户{id}
  3，logged_user：当前登录用户信息，判断是否已关注(假如已登录的话)
*/ 
pub async fn user_detail(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    user_id: web::Path<String>,
    logged_user: AuthUser
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id_some();
    let mut vars = HashMap::new();
    vars.insert("$uid", user_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());

    let _user: Option<user_model::item::User> = client
        .find_one(dgraph::USER, vars)
        .await?;

    if _user.is_none() {
        let lang = get_accept_language(&req);
        let msg = get_text!(lang_ctx, lang, "user.err.not_found").unwrap();
        return Res::err(json!({ "message": msg })); 
    }

    let user = _user.unwrap();
    Res::ok(json!({ "data": user }))
}

