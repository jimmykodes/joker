#[derive(Debug, PartialEq)]
pub enum OpCode {
    OpReturn,
}

#[derive(Debug)]
pub struct Chunk {
    pub code: Vec<OpCode>,
}

impl Chunk {
    pub fn new() -> Chunk {
        Chunk { code: Vec::new() }
    }

    pub fn push(&mut self, op: OpCode) {
        self.code.push(op);
    }
}
