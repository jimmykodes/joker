use crate::tokens::token::Token;

pub fn lexer(input: &str) -> Vec<Token> {
    let mut tokens = Vec::new();
    let mut i = 0;
    let mut line = 1;
    loop {
        if i >= input.len() {
            break;
        }
        let c = input.chars().nth(i).unwrap();
        match c {
            '\n' => line += 1,
            '(' => tokens.push(Token::LPar(line)),
            ')' => tokens.push(Token::RPar(line)),
            '{' => tokens.push(Token::LBrace(line)),
            '}' => tokens.push(Token::RBrace(line)),
            '[' => tokens.push(Token::LBrack(line)),
            ']' => tokens.push(Token::RBrack(line)),
            ',' => tokens.push(Token::Comma(line)),
            '0'..='9' => {
                let (tok, consumed) = number(&line, &input[i..]);
                i += consumed;
                tokens.push(tok);
            }
            'a'..='z' | 'A'..='Z' | '_' => {
                let (tok, consumed) = ident(&line, &input[i..]);
                i += consumed;
                tokens.push(tok);
            }
            _ => (),
        }
        i += 1
    }
    return tokens;
}

fn ident(ln: &i32, input: &str) -> (Token, usize) {
    let mut ident = String::new();
    let mut l = 0;

    for (i, c) in input.chars().enumerate() {
        match c {
            'a'..='z' | 'A'..='Z' | '0'..='9' | '_' => ident.push(c),
            _ => {
                l = i - 1;
                break;
            }
        }
    }
    return (Token::Identifier(*ln, ident), l);
}

fn number(ln: &i32, input: &str) -> (Token, usize) {
    let mut num = String::new();
    let mut l = 0;
    let mut int = true;

    for (i, c) in input.chars().enumerate() {
        match c {
            '_' => (),
            '0'..='9' => num.push(c),
            '.' => {
                num.push(c);
                int = false;
            }
            _ => {
                l = i - 1;
                break;
            }
        }
    }

    let tok = if int {
        Token::Int(*ln, num)
    } else {
        Token::Float(*ln, num)
    };

    return (tok, l);
}
