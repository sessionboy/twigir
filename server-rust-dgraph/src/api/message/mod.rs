use actix_web::{ web };

pub mod dgraph;
pub mod handlers;

use handlers::{
    mutation,
    query
};

pub fn route(cfg: &mut web::ServiceConfig) {
    cfg.service( 
        web::scope("/message")
            .service(web::resource("").route(web::post().to(mutation::message_create)))                       
                     
    );
}
