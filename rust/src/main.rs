use actix_web::{App, HttpServer};
use handlers::similarity::similarity;

mod handlers;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .service(similarity)
    })
    .bind(("127.0.0.1", 8081))?
    .run()
    .await
}