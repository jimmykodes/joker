mod opcodes;
use opcodes::{Chunk, OpCode};

fn main() {
    let mut c = Chunk::new();
    c.push(OpCode::OpReturn);
    println!("{:?}", c);
}
