use serde::{ Serialize, Deserialize };

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Field {
  pub field_name: String,
  pub alias_name: Option<String>,
  pub model_type: String,
  pub schema_str: String,

  pub unique: bool,
  pub default: Option<String>,
  pub default_type: Option<String>
}

// DgraphOption,
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DgraphModel {
  pub model_type: String,
  pub schema: Vec<Field>
}

// DgraphModel
pub fn model() {
  // 
}