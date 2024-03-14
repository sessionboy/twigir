use actix_web::{ web,HttpRequest };
use std::collections::HashMap;
use crate::lib::dgraph::{ DgraphClient };
use crate::api::user::dgraph::account as dgraph;
use crate::models::users as user_model;
use crate::lib::res::{ Res, ServiceResult };
use json_gettext::JSONGetText;
use crate::lib::auth::{ GuardAuthUser };
use crate::utils::password_hash::{ get_password_hash, verify_password };
use crate::lib::validate::validate;
use crate::utils::parse::{ get_accept_language };
use crate::utils::{ date };
use serde_json::{ json };

/*
  1，功能： 更新用户资料
  2，path: account/profile
  3，body：要更新的用户信息
*/ 
pub async fn profile(
    data: web::Data<DgraphClient>,
    body: web::Json<user_model::input::ProfileInput>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser
) -> ServiceResult {
    let lang = get_accept_language(&req);
    let client = &data.into_inner();

    // 验证参数
    match validate(&body) {
        Err(errors) => {
            let msg = get_text!(lang_ctx, lang, "user.err.input_verify").unwrap();
            return Res::err(json!({ "message": msg })); 
        },
        Ok(r) => r
    };

    let login_user = logged_user.to_slim();
    let mut profile = body.into_inner();
    profile.uid = Some(login_user.uid);
    profile.updated_at = Some(date::get_utc_now());

    client.mutate(&profile).await?;
    let msg = get_text!(lang_ctx, lang, "common.success").unwrap();
    Res::ok(json!({ 
        "data": {
            "uid": profile.uid,
        },
        "message": msg  
    }))
}


/*
  1，功能： 更新用户名字
  2，path: account/name
  3，body：要更新的用户信息
*/ 
pub async fn account_name(
    data: web::Data<DgraphClient>,
    body: web::Json<user_model::args::AccountName>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser
) -> ServiceResult {
    let lang = get_accept_language(&req);
    // 验证参数
    match validate(&body) {
        Err(errors) => {
            let msg = get_text!(lang_ctx, lang, "user.err.input_verify").unwrap();
            return Res::err(json!({ "message": msg })); 
        },
        Ok(r) => r,
    };
  
    let client = &data.into_inner();
    let login_user = logged_user.to_slim();
    let mut user = body.into_inner();
    user.uid = Some(login_user.uid);
    user.updated_at = Some(date::get_utc_now());

    // 查询参数
    let mut vars = HashMap::new();
    vars.insert("$name", &user.name);
    let exists: user_model::exist::ExistName = client
        .query(dgraph::EXIST_NAME, vars)
        .await?
        .try_into()?;

    // 名称已存在
    if !exists.name.is_empty() {
        let msg = get_text!(lang_ctx, lang, "user.err.exist_name").unwrap();
        return Res::err(json!({ "message": msg })); 
    }

    // 更新名字
    client.mutate(&user).await?;
    Res::ok(json!({ 
        "data": {
            "uid": user.uid,
            "name": user.name
        }
    }))
}


/*
  1，功能： 更新用户主名
  2，path: account/username
  3，body：要更新的用户信息
*/ 
pub async fn account_username(
    data: web::Data<DgraphClient>,
    body: web::Json<user_model::args::AccountUsername>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser
) -> ServiceResult {
    let lang = get_accept_language(&req);
    // 验证参数
    match validate(&body) {
        Err(errors) => {
            let msg = get_text!(lang_ctx, lang, "user.err.input_verify").unwrap();
            return Res::err(json!({ "message": msg })); 
        },
        Ok(r) => r,
    };

    let client = &data.into_inner();
    let login_user = logged_user.to_slim();
    let mut user = body.into_inner();
    user.uid = Some(login_user.uid);
    user.updated_at = Some(date::get_utc_now());

    // 查询参数
    let mut vars = HashMap::new();
    vars.insert("$username", &user.username);
    // 是否已存在
    let exists: user_model::exist::ExistUsername = client
        .query(dgraph::EXIST_USERNAME, vars)
        .await?
        .try_into()?;
    if !exists.username.is_empty() {
        let msg = get_text!(lang_ctx, lang, "user.err.exist_username").unwrap();
        return Res::err(json!({ "message": msg })); 
    }
    // 更新主名
    client.mutate(&user).await?;
    Res::ok(json!({ 
        "data": {
            "uid": user.uid,
            "username": user.username
        }
    }))
}


/*
  1，功能： 更改手机号
  2，path: account/phone
  3，body：要更新的用户信息
*/ 
pub async fn account_phone(
    data: web::Data<DgraphClient>,
    body: web::Json<user_model::args::AccountPhone>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser
) -> ServiceResult {
    let lang = get_accept_language(&req);
    // 验证参数
    match validate(&body) {
        Err(errors) => {
            let msg = get_text!(lang_ctx, lang, "user.err.input_verify").unwrap();
            return Res::err(json!({ "message": msg })); 
        },
        Ok(r) => r,
    };
    let client = &data.into_inner();
    let login_user = logged_user.to_slim();
    let mut user = body.into_inner();
    user.uid = Some(login_user.uid);
    user.updated_at = Some(date::get_utc_now());

    // 遗留项： 验证手机号验证码

    // 查询参数
    let mut vars = HashMap::new();
    vars.insert("$phone_number", &user.phone_number);
    // 是否已存在
    let exists: user_model::exist::ExistPhone = client
        .query(dgraph::EXIST_PHONE, vars)
        .await?
        .try_into()?;
    if !exists.phone.is_empty() {
        let msg = get_text!(lang_ctx, lang, "user.err.exist_phone").unwrap();
        return Res::err(json!({ "message": msg })); 
    }
    // 更新号码
    client.mutate(&user).await?;
    Res::ok(json!({ 
        "data": {
            "uid": user.uid,
            "phone_number": user.phone_number
        }
    }))
}


/*
  1，功能： 更改邮箱地址
  2，path: account/email
  3，body：要更新的用户信息
*/ 
pub async fn account_email(
    data: web::Data<DgraphClient>,
    body: web::Json<user_model::args::AccountEmail>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser
) -> ServiceResult {
    let lang = get_accept_language(&req);
    // 验证参数
    match validate(&body) {
        Err(errors) => {
            let msg = get_text!(lang_ctx, lang, "user.err.input_verify").unwrap();
            return Res::err(json!({ "message": msg }));  
        },
        Ok(r) => r
    };

    let client = &data.into_inner();
    let login_user = logged_user.to_slim();
    let mut user = body.into_inner();
    user.uid = Some(login_user.uid);
    user.updated_at = Some(date::get_utc_now());

    // 遗留项： 验证邮箱验证码
    
    // 查询参数
    let mut vars = HashMap::new();
    vars.insert("$email", &user.email);
    // 是否已存在
    let exists: user_model::exist::ExistEmail = client
        .query(dgraph::EXIST_EMAIL, vars)
        .await?
        .try_into()?;
    if !exists.email.is_empty() {
        let msg = get_text!(lang_ctx, lang, "user.err.exist_email").unwrap();
        return Res::err(json!({ "message": msg })); 
    }
    // 更新email
    client.mutate(&user).await?;
    Res::ok(json!({ 
        "data": {
            "uid": user.uid,
            "email": user.email
        }
    }))
}


/*
  1，功能： 更改密码
  2，path: account/password
  3，body：要更新的用户信息
*/ 
pub async fn account_password(
    data: web::Data<DgraphClient>,
    body: web::Json<user_model::args::AccountPassword>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser
) -> ServiceResult {
    let lang = get_accept_language(&req);
    // 验证参数
    match validate(&body) {
        Err(errors) => {
            let msg = get_text!(lang_ctx, lang, "user.err.input_verify").unwrap();
            return Res::err(json!({ "message": msg })); 
        },
        Ok(r) => r,
    };

    let client = &data.into_inner();
    let login_user = logged_user.to_slim();
    let args = body.into_inner();

     // 遗留项： 验证邮箱/手机验证码

    // 查询参数
    let mut vars = HashMap::new();
    vars.insert("$uid", login_user.clone().uid);
    let user: Option<user_model::item::Password> = client
        .find_one(dgraph::QUERY_PASSWORD, vars)
        .await?;
    if user.is_none() {
        let msg = get_text!(lang_ctx, lang, "user.err.not_found").unwrap();
        return Res::err(json!({ "message": msg })); 
    }

    let me = user.unwrap();

    // 验证旧密码是否正确
    let is_password_matches = verify_password(&me.password, &args.old_password);
    if !is_password_matches {
        let msg = get_text!(lang_ctx, lang, "user.err.old_password_mismatch").unwrap();
        return Res::err(json!({ "message": msg })); 
    }

    // 生成新密码
    let password = get_password_hash(args.clone().password).hash;
    // 更新密码
    let p = json!({
        "uid": login_user.clone().uid,
        "password": password.as_str(),
        "updated_at": date::get_utc_now()
    });   
    client.mutate(&p).await?;
    Res::ok(json!({ 
        "data": {
            "uid": login_user.uid
        }
    }))
}