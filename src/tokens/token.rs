#[derive(Debug, PartialEq)]
pub enum Token {
    Identifier(i32, String),
    Int(i32, String),
    Float(i32, String),
    LPar(i32),
    RPar(i32),
    LBrace(i32),
    RBrace(i32),
    LBrack(i32),
    RBrack(i32),
    Comma(i32),
}
