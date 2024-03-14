use actix_web::{ web };

pub mod dgraph;
pub mod handlers;
pub mod services;

use handlers::{
    mutate_group,
    mutate_admin,
    query
};
use crate::api::status::handlers::mutate_status;

pub fn route(cfg: &mut web::ServiceConfig) {
    cfg.service( 
        web::scope("/group")
            .service(web::resource("").route(web::post().to(mutate_group::group_create)))                       
            .service(web::resource("/{id}").route(web::get().to(query::group))) 
            // 小组发帖
            .service(web::resource("/{id}/publish").route(web::post().to(mutate_status::group_status_publish)))
            .service(web::resource("/{id}/status").route(web::get().to(query::group_status))) 
            .service(web::resource("/{id}/update").route(web::put().to(mutate_group::group_update)))    
            .service(web::resource("/{id}/destroy").route(web::delete().to(mutate_group::group_delete)))   
            // 小组管理 
            .service(web::resource("/{id}/join")
                .route(web::post().to(mutate_admin::group_join))
            ) 
            .service(web::resource("/{id}/leave")
                .route(web::post().to(mutate_admin::group_leave))
            )      
            .service(web::resource("/{id}/add_admin/{member_id}")
                .route(web::post().to(mutate_admin::group_add_admin))
            )  
            .service(web::resource("/{id}/fire_admin/{member_id}")
                .route(web::post().to(mutate_admin::group_fire_admin))
            )     
            .service(web::resource("/{id}/forbidden/{member_id}")
                .route(web::post().to(mutate_admin::group_forbidden))
            )            
    );
}
