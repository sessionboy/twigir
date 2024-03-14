use actix_web::{ web,HttpRequest };
use std::collections::HashMap;
use dgraph_tonic::{ Mutation, Mutate };
use crate::models::statuses as status_model;
use crate::models::replies as reply_model;
use crate::api::status::dgraph::mutation as dgraph;
use crate::lib::dgraph::{ DgraphClient };
use crate::api::status::services::status as service;
use crate::lib::auth::{ GuardAuthUser };
use json_gettext::JSONGetText;
use crate::lib::res::{ Res, ServiceResult };
use crate::utils::{ date, parse };
use serde_json::json;

/*
  1，功能：发布回复
  2，path: /reply
  3，body：贴文信息
  4，logged_user：当前登录用户信息
*/ 
pub async fn reply_publish(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser,
    body: web::Json<reply_model::input::CreateReplyInput>,
) -> ServiceResult {
    let client = &data.into_inner();
    let inputs = &body.into_inner();
    let login_user = logged_user.to_slim();
    let login_user_id = login_user.uid.clone();
    let lang = parse::get_accept_language(&req);

    // 1，检查最近一次回复时间是否少于60秒，防止恶意刷回复
    let mut user_vars = HashMap::new(); 
    user_vars.insert("$uid", login_user_id.as_str());

    let _user: Option<status_model::item::StatusCreater> = client
        .find_one(dgraph::STATUS_CREATE, user_vars)
        .await?;
    let user = _user.unwrap();
    if user.last_reply_at.is_some() {
        let last_reply_at = user.last_reply_at.unwrap();
        if date::compare_from_now_secs(&last_reply_at) < 60 {
            let msg = get_text!(lang_ctx, lang, "common.err.too_frequently").unwrap();
            return Res::err(json!({ "message": msg }));
        }
    }

    // 2，查询帖子信息
    let mut vars = HashMap::new(); 
    vars.insert("$uid", inputs.status_id.as_str());
    let _status: Option<status_model::item::StatusWithCount> = client
        .find_one(dgraph::STATUS_WITH_COUNT, vars)
        .await?;    
    // 帖子不存在
    if _status.is_none() {
        let msg = get_text!(lang_ctx, lang, "status.err.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    let status = _status.unwrap();
    let mut txn = client.mutate_txn();

    // 3，创建回复
    let mut reply_json = service::parse_reply(inputs,login_user_id.clone());
    reply_json["status"] = json!({ "uid": inputs.status_id });
    let created_at = date::get_utc_now();
    let p = json!({
        "uid": inputs.status_id.as_str(),
        "replies": reply_json,
        "replies_count": status.replies_count + 1
    });
    let mut mu = Mutation::new();
    mu.set_set_json(&p)?;
    let resp = txn.mutate(mu).await?;
    let reply_id = &resp.uids["reply"];

    // 4，更新用户的回复时间
    let user_reply_date = json!({
        "uid": login_user_id.as_str(),
        "last_reply_at": created_at.clone()
    });
    let mut mu_at = Mutation::new();
    mu_at.set_set_json(&user_reply_date)?;
    txn.mutate(mu_at).await?;
    
    // 遗留：创建feed

    // 遗留：通知提及用户，以及粉丝
    // if inputs.entities.is_some() {
    //     let mentions = &inputs.clone().entities.unwrap().mentions;
    // }

    let res = txn.commit().await;
    if res.is_err() {
        let msg = get_text!(lang_ctx, lang, "common.action_failure").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    return Res::ok(json!({ "data": { "reply_id": reply_id } }));
}

/*
  1，功能：喜欢/点赞该回复
  2，path: /reply/{id}/favorite
  3，logged_user：当前登录用户信息
*/ 
pub async fn reply_favorite(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser,
    reply_id: web::Path<String>
) -> ServiceResult {
    let client = &data.into_inner();
    let login_user = logged_user.to_slim();
    let logged_user_id = login_user.uid.clone();
    let lang = parse::get_accept_language(&req);

    // 1，查询回复
    let mut vars = HashMap::new(); 
    vars.insert("$uid", reply_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());
    let _reply: Option<status_model::item::Favorite> = client
        .find_one(dgraph::STATUS_OR_REPLY_FAVORITE_INFO, vars)
        .await?;
    
    // 回复不存在
    if _reply.is_none() {
        let msg = get_text!(lang_ctx, lang, "reply.err.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    let reply = _reply.unwrap();
    let mut txn = client.mutate_txn();

    // 2-1，已点赞，则取消点赞
    if reply.is_favorite {
        // 1，删除边
        let p1 = json!({
            "uid": reply_id.as_str(),
            "favorites":{
                "uid": logged_user_id
            }
        });
        let mut mu_del = Mutation::new();
        mu_del.set_delete_json(&p1)?;
        txn.mutate(mu_del).await?;

        // 2，更新点赞数量
        let p2 = json!({
            "uid": reply_id.as_str(),
            "favorites_count": reply.favorites_count - 1
        });
        let mut mu_mut = Mutation::new();        
        mu_mut.set_set_json(&p2)?;        
        txn.mutate(mu_mut).await?;
    } 
    // 2-2，没点赞，则添加点赞
    else {        
        let p = json!({
            "uid": reply_id.as_str(),
            "favorites":{
                "uid": logged_user_id
            },
            "favorites_count": reply.favorites_count + 1
        });
        let mut mu = Mutation::new();
        mu.set_set_json(&p)?;
        txn.mutate(mu).await?;
    }

    let res = txn.commit().await;
    if res.is_err() {
        let msg = get_text!(lang_ctx, lang, "common.action_failure").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    
    let msg = get_text!(lang_ctx, lang, "common.success").unwrap();
    Res::ok(json!({ "message": msg }))
}

/*
  1，功能：删除回复
  2，path: /reply
  3，body：贴文信息
  4，logged_user：当前登录用户信息
*/ 
pub async fn reply_delete(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser,
    reply_id: web::Path<String>
) -> ServiceResult {
    let client = &data.into_inner();
    let login_user = logged_user.to_slim();
    let login_user_id = login_user.uid.clone();
    let lang = parse::get_accept_language(&req);

    // 查询回复信息
    let mut vars = HashMap::new(); 
    vars.insert("$uid", reply_id.as_str());
    let _reply: Option<reply_model::query::ReplyInfo> = client
        .find_one(dgraph::REPLY_INFO, vars)
        .await?;
    
    // 回复不存在
    if _reply.is_none() {
        let msg = get_text!(lang_ctx, lang, "reply.err.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    // 非回复的作者，不能删除该回复
    let reply = _reply.unwrap();
    if reply.user.uid != login_user_id {
        let msg = get_text!(lang_ctx, lang, "reply.err.delete.not_owner").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    let mut txn = client.mutate_txn();

    // 删除回复，以及实体信息
    let mut _vars = HashMap::new(); 
    _vars.insert("$uid", reply_id.as_str());
    let mut mu = Mutation::new();
    mu.set_delete_nquads(dgraph::STATUS_OR_REPLY_DELETE);
    txn.upsert_with_vars(dgraph::STATUS_OR_REPLY_DELETE_QUERY, _vars, mu).await?;

    if reply.status.is_some() {
        let status = reply.status.unwrap();
        // 删除帖子关联的回复的边
        let p1 = json!({
            "uid": status.uid.as_str(),
            "replies":{
                "uid": reply_id.as_str()
            }
        });
        let mut mu_del = Mutation::new();
        mu_del.set_delete_json(&p1)?;
        txn.mutate(mu_del).await?;

        // 更新帖子的回复总数
        let p2 = json!({
            "uid": status.uid.as_str(),
            "replies_count": status.replies_count - 1
        });
        let mut mu_mut = Mutation::new();        
        mu_mut.set_set_json(&p2)?;        
        txn.mutate(mu_mut).await?;
    }

    let res = txn.commit().await;
    if res.is_err() {
        let msg = get_text!(lang_ctx, lang, "common.action_failure").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    
    let msg = get_text!(lang_ctx, lang, "common.success").unwrap();
    return Res::ok(json!({ "message": msg }));
}