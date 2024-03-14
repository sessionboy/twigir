use serde::{ Deserialize, Serialize };

// 详细的验证信息
#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Verified {
    pub verify_type: String,
    pub name: String,
    pub level: u8,
    pub description: String,
    pub created_at: Option<String>,
    pub updated_at: Option<String>
}

// 缩减版验证信息，用于展示
#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SlimVerified {
    pub uid: String,
    pub name: String,
    pub description: String
}
