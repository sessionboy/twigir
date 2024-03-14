use actix_web::{ web };

pub mod dgraph;
pub mod handlers;

use handlers::{
    mutation,
    query
};

pub fn route(cfg: &mut web::ServiceConfig) {
    cfg.service( 
        web::scope("/notification")
            .service(web::resource("").route(web::post().to(mutation::notification_create)))                       
                     
    );
}
