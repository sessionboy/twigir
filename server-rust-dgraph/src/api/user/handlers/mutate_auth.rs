use actix_web::{web, HttpRequest };
use crate::api::user::dgraph::auth as dgraph;
use crate::models::users as user_model;
use crate::models::users::input::{ RegisterInput };
use crate::models::users::_enum::{ Role };
use crate::models::users::auth::{ 
    AuthPayload,SlimUser,LoggedUser
};
use crate::lib::res::{ Res, ServiceResult };
use json_gettext::JSONGetText;
use crate::lib::dgraph::{ DgraphClient };
use crate::api::user::services::auth as service;
use crate::lib::validate::validate;
use std::collections::HashMap;
use actix_identity::Identity;
use crate::utils::password_hash::{ get_password_hash, verify_password };
use crate::utils::jwt::{ generate_token };
use crate::utils::parse::{ get_accept_language };
use crate::utils::{ date };
use serde_json::json;

/*
  1，功能：用户登录
  2，path: /auth/login
  3，body：账号/密码等用户参数
*/ 
pub async fn login(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    id: Identity,
    body: web::Json<user_model::args::LoginArg>,
) -> ServiceResult {
    let lang = get_accept_language(&req);
    let client = &data.into_inner();

    // 获取 ip 地址，
    // let addr = match req.peer_addr() {
    //     Some(addr) => addr,
    //     None => panic!(),
    // };
    // let ip = match addr {
    //     SocketAddr::V4(ipv4) => format!("{:?}", ipv4.ip()),
    //     SocketAddr::V6(ipv6) => format!("{:?}", ipv6.ip()),
    // };

    // 验证参数
    match validate(&body) {
        Err(errors) => {
            let msg = get_text!(lang_ctx, lang, "user.err.login.input_verify").unwrap();
            return Res::err(json!({ "message": msg }));
        },
        Ok(user) => user,
    };

    let (q, vars) = service::login_query(&body);
    let user: Option<user_model::auth::LoggedUserWithPassword> = client
        .find_one(q, vars)
        .await?;

    if user.is_none() {
        let msg = get_text!(lang_ctx, lang, "user.err.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    let login_user = user.unwrap();

    // 密码是否正确
    let is_password_matches = verify_password(&login_user.password, &body.password);
    if !is_password_matches {
        let msg = get_text!(lang_ctx, lang, "user.err.password_mismatch").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    // 组装登录数据
    let token = generate_token(SlimUser{
        uid: login_user.clone().uid,
        username: login_user.clone().username,
        role: Role::User
    });
    let logged_user = service::login_user(login_user.clone());
    let auth_payload = AuthPayload {
        token: token.clone(),
        user: logged_user
    };

    // 遗留问题 - 用户记录 (ip、注册平台等信息)

    id.remember(token);
    Res::ok(json!({ "data": auth_payload }))
}

/*
  1，功能：退出登录
  2，path: /auth/logout
*/ 
pub async fn logout(
    id: Identity,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest
) -> ServiceResult {
    let lang = get_accept_language(&req);
    // 清除cookie退出登录
    id.forget();
    let msg = get_text!(lang_ctx, lang, "user.err.logout").unwrap();
    Res::ok(json!({ "message": msg }))
}

/*
  1，功能： 用户注册
  2，path: /auth/register
  3，body：注册用户信息
  4，遗留问题： 用户记录 (ip、注册平台等信息)
*/ 
pub async fn register(
    data: web::Data<DgraphClient>,
    body: web::Json<user_model::input::RegisterInput>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    id: Identity,
    // generator: web::Data<CsrfTokenGenerator>
) -> ServiceResult {

    let lang = get_accept_language(&req);
    let client = &data.into_inner();

    // 验证参数
    match validate(&body) {
        Err(errors) => {
            let msg = get_text!(lang_ctx, lang, "user.err.input_verify").unwrap();
            return Res::err(json!({ "message": msg }));
        },
        Ok(user) => user,
    };
   
    let created_at = date::get_utc_now();
    let mut user = RegisterInput {
        lang: Some(lang.clone()),
        profile_default_cover: Some(false),
        followers_count: Some(0),
        followings_count: Some(0),
        friends_count: Some(0),
        statuses_count: Some(0),
        replies_count: Some(0),
        is_verified: Some(false),
        receive_chat: Some(true),
        receive_like_notification: Some(true),
        receive_reply_notification: Some(true),
        created_at: Some(created_at.clone()),
        ..body.into_inner()
    };

    // 检查name/username/phone_number是否已注册
    let mut vars = HashMap::new();
    vars.insert("$name", &user.name);
    vars.insert("$username", &user.username);
    vars.insert("$phone", &user.phone_number);
    let user_exist: user_model::exist::UserExist = client
        .query(dgraph::USER_EXIST, vars)
        .await?
        .try_into()?;

    let exist_errors = service::register_exist_errors(&user_exist, &user);
    if !exist_errors.is_empty() {
        let error = exist_errors[0].clone();
        let mut msg_txt: &str = "common.unknown";
        // 名称已存在
        if error["path"] == String::from("name") {
            msg_txt = "user.err.exist_name";
        } 
        // 主名已存在
        if error["path"] == String::from("username") {
            msg_txt = "user.err.exist_username";
        } 
        // 手机号已存在
        if error["path"] == String::from("phone_number") {
            msg_txt = "user.err.exist_phone";
        } 
        let msg = get_text!(lang_ctx, lang, msg_txt).unwrap();
        return Res::err(json!({ "message": msg }));
    }
    
    // 生成密码
    let password = get_password_hash(user.clone().password).hash;
    user.password = password;

    // 创建用户 
    let mut new_user = serde_json::to_value(&user).unwrap();
    new_user["dgraph.type"] = json!("User");
    new_user["uid"] = json!("_:user_id");

    let response = client.mutate(&new_user).await?;
    let user_id = response.uids["user_id"].clone();

    // 组装返回的数据
    let token = generate_token(SlimUser{
        uid: user_id.clone(),
        username: user.clone().username,
        role: user.clone().role.expect("User")
    });
    let logged_user = LoggedUser {
        uid: user_id,
        name: user.clone().name,
        username: user.clone().username,
        lang: user.clone().lang.expect("zh_CN"),
        role: user.clone().role.expect("User"),
        avatar_url: user.clone().avatar_url,
        is_verified: false,
        ..Default::default()
    };
    let auth_payload = AuthPayload {
        token: token.clone(),
        user: logged_user
    };

    // 存储cookie
    id.remember(token);

    // 遗留问题 - 用户记录 (ip、注册平台等信息)

    Res::ok(json!({ "data": auth_payload }))
}


/*
  1，功能： 发送手机验证码
  2，path: /auth/send_phonecode
*/ 
pub async fn send_phonecode(
    data: web::Data<DgraphClient>,
    body: web::Json<user_model::args::VerifyPhoneCode>,
) -> ServiceResult {
    // 验证参数
    match validate(&body) {
        Err(errors) => {
            return Res::err(json!({ "message": "参数错误" }));
        },
        Ok(r) => r
    };
    // 清除cookie退出登录
    println!("{:?}",body);
    Res::ok(json!({ "data": body.into_inner() }))
}


/*
  1，功能： 验证手机验证码
  2，path: /auth/verify_phonecode
*/ 
pub async fn verify_phonecode(
    data: web::Data<DgraphClient>,
    body: web::Json<user_model::args::VerifyPhoneCode>,
) -> ServiceResult {
    // 验证参数
    match validate(&body) {
        Err(errors) => {
            return Res::err(json!({ "message": "参数错误" }));  
        },
        Ok(r) => r,
    };
    // 清除cookie退出登录
    println!("{:?}",body);
    Res::ok(json!({ "data": body.into_inner() }))
}


/*
  1，功能： 发送邮箱验证码
  2，path: /auth/send_emailcode
*/ 
pub async fn send_emailcode(
    data: web::Data<DgraphClient>,
    body: web::Json<user_model::args::VerifyEmailCode>
) -> ServiceResult {
    // 验证参数
    match validate(&body) {
        Err(errors) => {
            return Res::err(json!({ "message": "参数错误" }));
        },
        Ok(r) => r,
    };
    // 清除cookie退出登录
    println!("{:?}",body);
    Res::ok(json!({ "data": body.into_inner() }))
}


/*
  1，功能： 验证邮箱地址
  2，path: /auth/verify_emailcode
*/ 
pub async fn verify_emailcode(
    data: web::Data<DgraphClient>,
    body: web::Json<user_model::args::VerifyEmailCode>
) -> ServiceResult {
    // 验证参数
    match validate(&body) {
        Err(errors) => {
            return Res::err(json!({ "message": "参数错误" }));
        },
        Ok(r) => r,
    };
    // 清除cookie退出登录
    println!("{:?}",body);
    Res::ok(json!({ "data": body.into_inner() }))
}