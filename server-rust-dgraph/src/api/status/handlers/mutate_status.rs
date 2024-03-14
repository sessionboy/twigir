use actix_web::{ web,HttpRequest };
use std::collections::HashMap;
use dgraph_tonic::{ Mutation, Mutate };
use crate::models::groups::member::IsGroupMember;
use crate::models::statuses as status_model;
use crate::api::status::dgraph::mutation as dgraph;
use crate::lib::dgraph::{ DgraphClient };
use crate::api::status::services::status as service;
use crate::lib::auth::{ GuardAuthUser };
use json_gettext::JSONGetText;
use crate::lib::res::{ Res, ServiceResult };
use crate::utils::{ date, parse };
use serde_json::json;

/*
  1，功能：发布贴子
  2，path: /status
  3，body：贴文信息
  4，logged_user：当前登录用户信息
*/ 
pub async fn status_publish(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser,
    body: web::Json<status_model::input::CreateStatusInput>,
) -> ServiceResult {
    let client = &data.into_inner();
    let inputs = &body.into_inner();
    let login_user = logged_user.to_slim();
    let login_user_id = login_user.uid.clone();
    let lang = parse::get_accept_language(&req);

    let mut vars = HashMap::new(); 
    vars.insert("$uid", login_user_id.as_str());
    let _user: Option<status_model::item::StatusCreater> = client
        .find_one(dgraph::STATUS_CREATE, vars)
        .await?;
    let user = _user.unwrap();
    
    // 检查最近一次发帖时间是否少于60秒，防止恶意刷帖
    if user.last_publish_at.is_some() {
        let last_publish_at = user.last_publish_at.unwrap();
        if date::compare_from_now_secs(&last_publish_at) < 60 {
            let msg = get_text!(lang_ctx, lang, "common.err.too_frequently").unwrap();
            return Res::err(json!({ "message": msg }));
        }
    }

    let mut txn = client.mutate_txn();

    // 1，发布帖子
    let status_json = service::parse_status(inputs,login_user_id.clone());
    let created_at = date::get_utc_now();
    let p = json!({
        "uid": login_user_id.as_str(),
        "statuses": status_json,
        "last_publish_at": created_at.clone(),
        "statuses_count": user.statuses_count + 1
    });
    let mut mu = Mutation::new();
    mu.set_set_json(&p)?;
    let user_resp = txn.mutate(mu).await?;
    let status_id = &user_resp.uids["status"];

    // 2，如果是转发，更新原帖子的转发信息
    if inputs.forward_to_status.is_some() {
        let forward_id = inputs.clone().forward_to_status.unwrap();
        let mut vars = HashMap::new(); 
        vars.insert("$uid", forward_id.as_str());
        let _forward: Option<status_model::item::StatusInfo> = client
            .find_one(dgraph::STATUS_INFO, vars)
            .await?;
        // 如果原贴还在，则更新原贴转发数量
        if _forward.is_some() {
            let forward = _forward.unwrap();
            let p = json!({ 
                "uid": forward.uid.as_str(),
                "forwards": {
                    "uid": login_user_id  // 转发者
                },
                "forwards_count": forward.forwards_count + 1
            });
            let mut mu = Mutation::new();
            mu.set_set_json(&p)?;
            txn.mutate(mu).await?;
        }        
    }

    // 遗留：创建feed

    // 遗留：通知提及用户，以及粉丝
    if inputs.entities.is_some() {
        let mentions = &inputs.clone().entities.unwrap().mentions;
    }

    let res = txn.commit().await;
    if res.is_err() {
        let msg = get_text!(lang_ctx, lang, "common.action_failure").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    return Res::ok(json!({ "data": { "status_id": status_id } }));
}

/*
  1，功能：在小组发帖
  2，path: /group/{id}/publish
  3，body：贴文信息
  4，logged_user：当前登录用户信息
*/ 
pub async fn group_status_publish(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser,
    body: web::Json<status_model::input::CreateStatusInput>,
    group_id: web::Path<String>
) -> ServiceResult {
    let client = &data.into_inner();
    let inputs = &body.into_inner();
    let login_user = logged_user.to_slim();
    let login_user_id = login_user.uid.clone();
    let lang = parse::get_accept_language(&req);

    // 查询小组
    let mut group_vars = HashMap::new(); 
    group_vars.insert("$uid", group_id.as_str());
    group_vars.insert("$logged_user_id", login_user_id.as_str());
    let _group: Option<IsGroupMember> = client
        .find_one(dgraph::USER_IS_GROUP_MEMBER, group_vars)
        .await?;
    
    // 小组不存在
    if _group.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    let group = _group.unwrap();
 
    // 不是该小组的成员
    if group.members.is_empty() {
        let msg = get_text!(lang_ctx, lang, "group.err.publish.unauthorized").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    let member = group.members.first().unwrap();

    // 检查最近一次发帖时间是否少于60秒，防止恶意刷帖
    if member.last_publish_at.is_some() {
        let last_publish_at = member.last_publish_at.clone().unwrap();
        if date::compare_from_now_secs(&last_publish_at) < 60 {
            let msg = get_text!(lang_ctx, lang, "common.err.too_frequently").unwrap();
            return Res::err(json!({ "message": msg }));
        }
    }

    // 是否被禁言 
    if member.forbidden_date.is_some() {
        let forbidden_date = member.forbidden_date.clone().unwrap();
        if date::before_from_now(&forbidden_date) {
            let msg = get_text!(lang_ctx, lang, "group.err.forbidden").unwrap();
            return Res::err(json!({ "message": msg }));
        }
    }
    
    let mut txn = client.mutate_txn();

    // 1，发布小组帖子
    let created_at = date::get_utc_now();
    let mut status_json = service::parse_status(inputs,login_user_id.clone());
    status_json["group"] = json!({ "uid": group_id.as_str() });   
    let p = json!({
        "uid": group_id.as_str(),
        "statuses": status_json,
        "members": {
            "uid": member.uid.as_str(),
            "last_publish_at": created_at.clone()
        },        
        "statuses_count": group.statuses_count + 1
    });
    let mut mu = Mutation::new();
    mu.set_set_json(&p)?;
    let user_resp = txn.mutate(mu).await?;
    let status_id = &user_resp.uids["status"];

    // 2，如果是转发，更新原帖子的转发信息
    if inputs.forward_to_status.is_some() {
        let forward_id = inputs.clone().forward_to_status.unwrap();
        let mut vars = HashMap::new(); 
        vars.insert("$uid", forward_id.as_str());
        let _forward: Option<status_model::item::StatusInfo> = client
            .find_one(dgraph::STATUS_INFO, vars)
            .await?;

        // 如果原贴还在，则更新转发信息
        if _forward.is_some() {
            let forward = _forward.unwrap();
            let p = json!({ 
                "uid": forward.uid.as_str(),
                "forwards": {
                    "uid": login_user_id  // 转发者
                },
                "forwards_count": forward.forwards_count + 1
            });
            let mut mu = Mutation::new();
            mu.set_set_json(&p)?;
            txn.mutate(mu).await?;
        }        
    }

    // 遗留：创建feed

    // 遗留：通知提及用户，以及粉丝
    if inputs.entities.is_some() {
        let mentions = &inputs.clone().entities.unwrap().mentions;
    }

    let res = txn.commit().await;
    if res.is_err() {
        let msg = get_text!(lang_ctx, lang, "common.action_failure").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    return Res::ok(json!({ "data": { "status_id": status_id } }));
}



/*
  1，功能：喜欢/点赞该贴子
  2，path: /status/{id}/favorite
  3，logged_user：当前登录用户信息
*/ 
pub async fn status_favorite(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser,
    status_id: web::Path<String>
) -> ServiceResult {
    let client = &data.into_inner();
    let login_user = logged_user.to_slim();
    let logged_user_id = login_user.uid.clone();
    let lang = parse::get_accept_language(&req);

    let mut vars = HashMap::new(); 
    vars.insert("$uid", status_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());
    let _status: Option<status_model::item::Favorite> = client
        .find_one(dgraph::STATUS_OR_REPLY_FAVORITE_INFO, vars)
        .await?;

    // 帖子不存在
    if _status.is_none() {
        let msg = get_text!(lang_ctx, lang, "status.err.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    let status = _status.unwrap();
    
    let mut txn = client.mutate_txn();

    // 已点赞，则取消点赞
    if status.is_favorite {
        // 1，删除边
        let p1 = json!({
            "uid": status_id.as_str(),
            "favorites":{
                "uid": logged_user_id
            }
        });
        let mut mu_del = Mutation::new();
        mu_del.set_delete_json(&p1)?;
        txn.mutate(mu_del).await?;

        // 2，更新点赞数量
        let p2 = json!({
            "uid": status_id.as_str(),
            "favorites_count": status.favorites_count - 1
        });
        let mut mu_mut = Mutation::new();        
        mu_mut.set_set_json(&p2)?;        
        txn.mutate(mu_mut).await?;
    } 
    // 没点赞，则添加点赞
    else {        
        let p = json!({
            "uid": status_id.as_str(),
            "favorites":{
                "uid": logged_user_id
            },
            "favorites_count": status.favorites_count + 1
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
  1，功能：删除贴子
  2，path: /status/{id}
  3，logged_user：当前登录用户信息
*/ 
pub async fn status_delete(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser,
    status_id: web::Path<String>
) -> ServiceResult {
    let client = &data.into_inner();
    let login_user = logged_user.to_slim();
    let login_user_id = login_user.uid.clone();
    let lang = parse::get_accept_language(&req);

    // 查找帖子
    let mut vars = HashMap::new(); 
    vars.insert("$uid", status_id.as_str());
    let _status: Option<status_model::item::QueryStatus> = client
        .find_one(dgraph::STATUS_INFO, vars)
        .await?;
    
    // 若帖子不存在
    if _status.is_none() {
        let msg = get_text!(lang_ctx, lang, "status.err.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    // 非帖子的作者，不能删除该帖子
    let status = _status.unwrap();
    if status.user.uid != login_user_id {
        let msg = get_text!(lang_ctx, lang, "status.err.delete.not_owner").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    let mut txn = client.mutate_txn();
    
    // 1，删除帖子，以及实体信息
    let mut _vars = HashMap::new(); 
    _vars.insert("$uid", status_id.as_str());
    let mut mu = Mutation::new();
    mu.set_delete_nquads(dgraph::STATUS_OR_REPLY_DELETE);
    txn.upsert_with_vars(dgraph::STATUS_OR_REPLY_DELETE_QUERY, _vars, mu).await?;

    // 2，删除关联该帖子的边
    let mut mu_del = Mutation::new();
    // 小组帖子
    if status.group.is_some() {
        let group = status.group.unwrap();
        let group_status_delete = json!({
            "uid": group.uid.as_str(),
            "statuses":{
                "uid": status_id.as_str()
            }
        });
        mu_del.set_delete_json(&group_status_delete)?;
        txn.mutate(mu_del).await?;

        // 更新小组帖子数量
        let group_status_count = json!({
            "uid": group.uid.as_str(),
            "statuses_count": group.statuses_count - 1
        });
        let mut mu_mut = Mutation::new();
        mu_mut.set_set_json(&group_status_count)?;        
        txn.mutate(mu_mut).await?;
    }
    // 非小组帖子
    else{
        let user_status_delete = json!({
            "uid": login_user_id.as_str(),
            "statuses":{
                "uid": status_id.as_str()
            }
        });
        mu_del.set_delete_json(&user_status_delete)?;
        txn.mutate(mu_del).await?;

        // 更新用户帖子数量
        let mut vars = HashMap::new(); 
        vars.insert("$uid", login_user_id.as_str());
        let _user: Option<status_model::item::StatusCreater> = client
            .find_one(dgraph::STATUS_CREATE, vars)
            .await?;
        let user = _user.unwrap();
        let user_status_count = json!({
            "uid": login_user_id.as_str(),
            "statuses_count": user.statuses_count - 1
        });
        let mut mu_user_count = Mutation::new();
        mu_user_count.set_set_json(&user_status_count)?;        
        txn.mutate(mu_user_count).await?;
    }
    
    // 3，如果是转发的帖子，则更新原帖子的转发信息
    if status.forward_to_status.is_some() {
        let forward = status.forward_to_status.unwrap();
        // 删除边
        let p1 = json!({
            "uid": forward.uid.as_str(),
            "forwards":{
                "uid": login_user_id.as_str()
            }
        });
        let mut mu_del = Mutation::new();
        mu_del.set_delete_json(&p1)?;
        txn.mutate(mu_del).await?;

        // 转发数量减1
        let p2 = json!({
            "uid": forward.uid.as_str(),            
            "forwards_count": forward.forwards_count - 1
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
    Res::ok(json!({ "message": msg }))
}
