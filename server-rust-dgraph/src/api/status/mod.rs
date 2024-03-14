use actix_web::{ web };

pub mod dgraph;
pub mod handlers;
pub mod services;

use handlers::{
    mutate_status,
    mutate_reply,
    query_status,
    query_reply
};

pub fn route(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/status")
            .service(web::resource("").route(web::post().to(mutate_status::status_publish))) 
            .service(web::resource("/user_home").route(web::get().to(query_status::status_user_homepage)))  
            .service(web::resource("/{id}").route(web::get().to(query_status::status)))  
            .service(web::resource("/{id}/replies").route(web::get().to(query_reply::status_replies)))  
            .service(web::resource("/{id}/favorite").route(web::post().to(mutate_status::status_favorite)))  
            .service(web::resource("/{id}/destroy").route(web::delete().to(mutate_status::status_delete)))           
    ).service(
        web::scope("/reply")   
            .service(web::resource("").route(web::post().to(mutate_reply::reply_publish))) 
            .service(web::resource("/{id}").route(web::get().to(query_reply::reply)))  
            .service(web::resource("/{id}/replies").route(web::get().to(query_reply::reply_replies)))  
            .service(web::resource("/{id}/favorite").route(web::post().to(mutate_reply::reply_favorite))) 
            .service(web::resource("/{id}/destroy").route(web::delete().to(mutate_reply::reply_delete)))    
    );
}
