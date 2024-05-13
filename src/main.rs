mod lexer {
    pub mod lexer;
}
mod tokens {
    pub mod token;
}

use lexer::lexer::lexer;

fn main() {
    let l = lexer(
        r#"
fn add(a, b){
    return a + b;
}

add(12, 4)
"#,
    );
    println!("{:?}", l);
}
