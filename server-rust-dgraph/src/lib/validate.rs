use actix_web::web::Json;
use validator::{Validate, ValidationErrors};
use serde_json::{json, Value as JsonValue};

// 参数校验
pub fn validate<T>(params: &Json<T>) -> Result<(), Vec<JsonValue>>
where
    T: Validate,
{
    match params.validate() {
        Ok(()) => Ok(()),
        Err(error) => {
            Err(collect_errors(error))
        },
    }
}

fn collect_errors(error: ValidationErrors) -> Vec<JsonValue> {
    let mut errors: Vec<JsonValue> = Vec::new();
    error
        .field_errors()
        .into_iter()
        .for_each(|e| {
            for item in e.1 {
                errors.push(json!({
                    "value": item.params.get("value"),
                    "path": e.0,
                    "message": item.message
                }));
            }
        });
    errors
}

