use actix_web::{ web, HttpRequest };
use crate::lib::dgraph::{ DgraphClient };
use crate::lib::res::{ Res, ServiceResult };
use json_gettext::JSONGetText;
use crate::lib::auth::{ GuardAuthUser };
use crate::utils::parse::{ get_accept_language };
use serde_json::json;
use crate::utils::{ date };

/*
  1，功能：关注某人
  2，path: /users/{id}/follow
  3，user_id：被关注的用户的id
  4，logged_user：当前登录用户信息
*/ 
pub async fn follow(
    data: web::Data<DgraphClient>,
    logged_user: GuardAuthUser,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    user_id: web::Path<String>
) -> ServiceResult {
    let lang = get_accept_language(&req);
    let login_user = logged_user.to_slim();
    if login_user.uid == user_id.as_str() {
        let msg = get_text!(lang_ctx, lang, "user.err.cannot_follow_self").unwrap();
        return Res::err(json!({ "message": msg })); 
    }
    
    let client = &data.into_inner();
    let p = json!({
        "uid": login_user.uid,
        "follows": {
            "uid": user_id.as_str()
        },
        "updated_at": date::get_utc_now()
    });
    client.mutate(&p).await?;
    Res::ok(json!({ 
        "data": { "uid": login_user.uid }
    }))
}


/*
  1，功能：取关某人
  2，path: /users/{id}/unfollow
  3，user_id：被取关的用户的id
  4，logged_user：当前登录用户信息
*/ 
pub async fn unfollow(
    data: web::Data<DgraphClient>,
    logged_user: GuardAuthUser,
    user_id: web::Path<String>
) -> ServiceResult {
    let login_user = logged_user.to_slim();
    let client = &data.into_inner();
    let p = json!({
        "uid": login_user.uid,
        "follows": {
            "uid": user_id.as_str()
        },
        "updated_at": date::get_utc_now()
    });
    client.delete(&p).await?;
    Res::ok(json!({ 
        "data": { "uid": login_user.uid }
    }))
}


/*
  1，功能：屏蔽某人
  2，path: /users/{id}/shield 
  3，user_id：被屏蔽的用户的id
  4，logged_user：当前登录用户信息
*/ 
pub async fn shield(
    data: web::Data<DgraphClient>,
    logged_user: GuardAuthUser,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    user_id: web::Path<String>
) -> ServiceResult {
    let lang = get_accept_language(&req);
    let login_user = logged_user.to_slim();
    if login_user.uid == user_id.as_str() {
        let msg = get_text!(lang_ctx, lang, "user.err.cannot_shield_self").unwrap();
        return Res::err(json!({ "message": msg })); 
    }
    let client = &data.into_inner();
    let p = json!({
        "uid": login_user.uid,
        "shields": {
            "uid": user_id.as_str(),
        },
        "updated_at": date::get_utc_now()
    });
    client.mutate(&p).await?;
    Res::ok(json!({ 
        "data": { "uid": login_user.uid }
    }))
}


/*
  1，功能：解除屏蔽某人
  2，path: /users/{id}/unshield
  3，user_id：被解除屏蔽的用户的id
  4，logged_user：当前登录用户信息
*/ 
pub async fn unshield(
    data: web::Data<DgraphClient>,
    logged_user: GuardAuthUser,
    user_id: web::Path<String>
) -> ServiceResult {
    let login_user = logged_user.to_slim();
    let client = &data.into_inner();
    let p = json!({
        "uid": login_user.uid,
        "shields": {
            "uid": user_id.as_str(),
        },
        "updated_at": date::get_utc_now()
    });
    client.delete(&p).await?;
    Res::ok(json!({ 
        "data": { "uid": login_user.uid }
    }))
}

